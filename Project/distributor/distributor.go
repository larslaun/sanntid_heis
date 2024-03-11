package distributor

import (
	"Elev-project/collector"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/hallAssigner"
	"Elev-project/settings"
	"fmt"
	"strconv"
	"time"
)

func DistributeState(elevStateTx chan elevator.Elevator, elevStateRx chan elevator.Elevator, localElev *elevator.Elevator) {
	for {
		//localElev.Available = true
		if localElev.Available == false {
			fmt.Print("\nthomAS\n")

			elevStateRx <- *localElev

		} else {
			elevStateTx <- *localElev
			
		}
		time.Sleep(20 * time.Millisecond)
	}
}

//psuedo distributor
//Receives buttonpress, then calculates optimal elevator wiht cost func,then sends elevOrder which includes order and ID of elev.

func DistributeOrder(buttonPress elevio.ButtonEvent, elevOrderTx chan collector.ElevatorOrder, elevOrderRx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, localElev *elevator.Elevator) {

	localID, _ := strconv.Atoi(localElev.ID)
	elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, *elevators, localID)

	//elevator.Elevator_print(*localElev)

	/*
		fmt.Printf("\nOptimal elev calculated:\n")
		fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
		fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
		fmt.Printf("Button: %d \n", elevOrder.Order.Button)
	*/

	if localElev.Available == false {
		fmt.Print("\n10\n")
		elevOrderRx <- elevOrder
	} else {
		if buttonPress.Button == elevio.BT_Cab {
			elevOrder.RecipientID = localElev.ID
		}
		elevOrderTx <- elevOrder

		transmissionFailures := 0

		for {
			select {
			case recievedState := <-elevStateRx:
				if recievedState.ID == elevOrder.RecipientID {
					//fmt.Print("1")
					//fmt.Print(recievedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button])
					if recievedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] || recievedState.Floor == elevOrder.Order.Floor {
						//fmt.Print("2")
						fmt.Print("CORRECT state recieved\n")
						return
					} else {
						fmt.Print("Wrong state recieved\n")
					}
				}

			case <-time.After(time.Millisecond * 40):
				transmissionFailures++
				//fmt.Printf("Transmission failures: %d\n", transmissionFailures)

				if transmissionFailures >= settings.MaxTransmissionFailures {

					RecieverID, _ := strconv.Atoi(elevOrder.RecipientID)
					elevators[RecieverID].Available = false

					elevOrder = hallAssigner.ChooseOptimalElev(buttonPress, *elevators, localID)

					if buttonPress.Button == elevio.BT_Cab {
						elevOrder.RecipientID = localElev.ID
						elevators[RecieverID].Available = true
					}
					transmissionFailures = 0
				}
				elevOrderTx <- elevOrder
			}
		}
	}
	/*
		fmt.Printf("\nOptimal elev calculated:\n")
		fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
		fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
		fmt.Printf("Button: %d \n", elevOrder.Order.Button)
	*/

}

func RedistributeFaultyElevOrders(elevOrderTx chan collector.ElevatorOrder, elevOrderRx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, faultyElev *elevator.Elevator) {
	fmt.Print("\nRedistribute initiated\n")
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
			if faultyElev.Requests[floor][btn] {
				hallCall := elevio.ButtonEvent{Floor: floor, Button: btn}
				go DistributeOrder(hallCall, elevOrderTx, elevOrderRx, elevStateRx, elevators, faultyElev)

				//buttonPress <- elevio.ButtonEvent{Floor: floor, Button: btn}
				//DistributeOrder(hallCall, elevOrderTx, elevStateRx, elevators, faultyElev)
				faultyElev.Requests[floor][btn] = false
			}
		}
	}
}

func RecoverCabOrders(elevOrderTx chan collector.ElevatorOrder, elevOrderRx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, faultyElev *elevator.Elevator) {
	fmt.Print("\nCab recovery initiated\n")
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		if faultyElev.Requests[floor][elevio.BT_Cab] {
			hallCall := elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
			go DistributeOrder(hallCall, elevOrderTx, elevOrderRx, elevStateRx, elevators, faultyElev)

			//hallCall := make(chan elevio.ButtonEvent)
			//buttonPress <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
			//DistributeOrder(hallCall, elevOrderTx, elevStateRx, elevators, faultyElev)
			faultyElev.Requests[floor][elevio.BT_Cab] = false
		}
	}
}
