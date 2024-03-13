package distributor

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/hallAssigner"
	"Elev-project/settings"
	"fmt"
	"strconv"
	"time"
)

// endre recievedState til receivedState og recieverID


func DistributeState(elevStateTx chan elevator.Elevator, localElev *elevator.Elevator) {
	for {
		elevStateTx <- *localElev
		time.Sleep(20 * time.Millisecond)
	}
}


func DistributeOrder(buttonPress elevio.ButtonEvent, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator) {

	localID, _ := strconv.Atoi(localElev.ID)
	elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, *elevators, localID)


	if localElev.NetworkAvailable == false {
		fmt.Print("\n10\n") //kan fjernes 
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
					if recievedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] || recievedState.Floor == elevOrder.Order.Floor {
						return
					}

			case <-time.After(time.Millisecond * 40):
				transmissionFailures++
				//fmt.Printf("Transmission failures: %d\n", transmissionFailures)

				if transmissionFailures >= settings.MaxTransmissionFailures {
					RecieverID, _ := strconv.Atoi(elevOrder.RecipientID)
					elevators[RecieverID].NetworkAvailable = false

					elevOrder = hallAssigner.ChooseOptimalElev(buttonPress, *elevators, localID)
					if buttonPress.Button == elevio.BT_Cab {
						elevOrder.RecipientID = localElev.ID
						elevators[RecieverID].NetworkAvailable = true
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

func RedistributeFaultyElevOrders(elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator) {
	fmt.Print("\nRedistribute initiated\n")
	for floor := 0; floor < settings.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
			if faultyElev.Requests[floor][btn] {
				hallCall := elevio.ButtonEvent{Floor: floor, Button: btn}
				go DistributeOrder(hallCall, elevOrderTx, elevOrderRx, elevStateRx, elevators, faultyElev)

				faultyElev.Requests[floor][btn] = false
			}
		}
	}
}

func RecoverCabOrders(elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator) {
	fmt.Print("\nCab recovery initiated\n")
	for floor := 0; floor < settings.N_FLOORS; floor++ {
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
