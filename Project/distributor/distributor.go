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

func DistributeState(elevStateTx chan elevator.Elevator, localElev *elevator.Elevator) {
	for {
		//localElev.Available = true
		elevStateTx <- *localElev
		time.Sleep(50 * time.Millisecond)
	}
}

//psuedo distributor
//Receives buttonpress, then calculates optimal elevator wiht cost func,then sends elevOrder which includes order and ID of elev.

func DistributeOrder(buttonPress chan elevio.ButtonEvent, elevOrderTx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, localElev *elevator.Elevator) {
		select {
		case buttonPress := <-buttonPress:
			elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators)
			
			elevator.Elevator_print(*localElev)

			fmt.Printf("\nOptimal elev calculated:\n")
			fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
			fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
			fmt.Printf("Button: %d \n", elevOrder.Order.Button)

			if buttonPress.Button == elevio.BT_Cab {
				elevOrder.RecipientID = localElev.ID
			}
			elevOrderTx <- elevOrder

			transmissionFailures := 0

			sendNewMsg := time.NewTimer(time.Millisecond * 200)

			for {
				select {
				case recievedState := <-elevStateRx:
					if recievedState.ID == elevOrder.RecipientID {
						fmt.Print("1")
						fmt.Print(recievedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button])
						if recievedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] {
							fmt.Print("2")
							return
						}
					}
				case <- sendNewMsg.C:
					transmissionFailures++
					fmt.Print("3")

					if transmissionFailures >= settings.MaxTransmissionFailures {

						RecieverID, _ := strconv.Atoi(elevOrder.RecipientID)
						elevators[RecieverID].Available = false

						elevOrder = hallAssigner.ChooseOptimalElev(buttonPress, elevators)

						if buttonPress.Button == elevio.BT_Cab {
							elevOrder.RecipientID = localElev.ID
							elevators[RecieverID].Available = true
						}
						transmissionFailures = 0
					}
					elevOrderTx <- elevOrder
				}
			}

			/*
				fmt.Printf("\nOptimal elev calculated:\n")
				fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
				fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
				fmt.Printf("Button: %d \n", elevOrder.Order.Button)
			*/

		}
	
}

func RedistributeFaultyElevOrders(elevOrderTx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, faultyElev *elevator.Elevator) {
	fmt.Print("\nRedistribute initiated\n")
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
			if faultyElev.Requests[floor][btn] {
				hallCall := make(chan elevio.ButtonEvent)
				hallCall <- elevio.ButtonEvent{Floor: floor, Button: btn}
				DistributeOrder(hallCall, elevOrderTx, elevStateRx, elevators, faultyElev)
				faultyElev.Requests[floor][btn] = false
			}
		}
	}
}
