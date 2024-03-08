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

func DistributeOrder(buttonPress chan elevio.ButtonEvent, elevOrderTx chan collector.ElevatorOrder, elevOrderRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator) {
	for {
		select {
		case buttonPress := <-buttonPress:
			elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators)

			if buttonPress.Button == elevio.BT_Cab {
				elevOrder.RecipientID = localID
			}
			elevOrderTx <- elevOrder

			transmissionFailures := 0

			for {
				select {
				case recievedState := <-elevOrderRx:
					if recievedState.ID == elevOrder.RecipientID {
						if recievedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] {
							return
						}
					}
				case <-time.After(time.Millisecond * 50):
					transmissionFailures++

					if transmissionFailures >= settings.MaxTransmissionFailures {

						RecieverID, _ := strconv.Atoi(elevOrder.RecipientID)
						elevators[RecieverID].Available = false

						elevOrder = hallAssigner.ChooseOptimalElev(buttonPress, elevators)

						if buttonPress.Button == elevio.BT_Cab {
							elevOrder.RecipientID = localID
							elevators[RecieverID].Available = true
						}
						elevOrderTx <- elevOrder
						transmissionFailures = 0
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
	}
}

func RedistributeFaultyElevOrders(elevOrderTx chan collector.ElevatorOrder, elevators *[settings.NumElevs]elevator.Elevator, faultyElev *elevator.Elevator, redistributeSignal chan bool) {
	for {
		select {
		case <-redistributeSignal:
			fmt.Print("\nRedistribute initiated\n")
			for floor := 0; floor < elevator.N_FLOORS; floor++ {
				for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
					if faultyElev.Requests[floor][btn] {
						var hallCall elevio.ButtonEvent = elevio.ButtonEvent{Floor: floor, Button: btn}
						DistributeOrder(hallCall, elevOrderTx, elevators)
						faultyElev.Requests[floor][btn] = false
					}
				}
			}
		}
	}
}
