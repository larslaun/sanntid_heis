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

func DistributeState(elevStateTx chan elevator.Elevator, localElev *elevator.Elevator) {
	for {
		elevStateTx <- *localElev
		time.Sleep(20 * time.Millisecond)
	}
}

func DistributeOrder(buttonPress elevio.ButtonEvent, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator, localID int) {

	elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, *elevatorArray, localID)

	if elevatorArray[localID].NetworkAvailable == false {
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
			case receivedState := <-elevStateRx:
				if receivedState.ID == elevOrder.RecipientID {
					if receivedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] || receivedState.Floor == elevOrder.Order.Floor {
						return
					} 
				}

			case <-time.After(time.Millisecond * 40):
				transmissionFailures++

				if transmissionFailures >= settings.MaxTransmissionFailures {

					ReceiverID, _ := strconv.Atoi(elevOrder.RecipientID)
					elevatorArray[ReceiverID].NetworkAvailable = false
					elevOrder = hallAssigner.ChooseOptimalElev(buttonPress, *elevatorArray, localID)

					if buttonPress.Button == elevio.BT_Cab {
						elevOrder.RecipientID = localElev.ID
						elevatorArray[ReceiverID].NetworkAvailable = true
					}
					transmissionFailures = 0
				}
				elevOrderTx <- elevOrder
			}
		} 
	}

}

func RedistributeFaultyElevOrders(elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator, localID int) {
	fmt.Print("\nRedistribute initiated\n")
	faultyElevID, _ := strconv.Atoi(faultyElev.ID)

	for floor := 0; floor < settings.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
			if faultyElev.Requests[floor][btn] {
				faultyElev.Requests[floor][btn] = false
				elevatorArray[faultyElevID].Requests[floor][btn] = false
				hallCall := elevio.ButtonEvent{Floor: floor, Button: btn}
				go DistributeOrder(hallCall, elevOrderTx, elevOrderRx, elevStateRx, elevatorArray, faultyElev, localID)
			}
		}
	}
}

func RecoverCabOrders(elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator, localID int) {
	fmt.Print("\nCab recovery initiated\n")

	for floor := 0; floor < settings.N_FLOORS; floor++ {
		if faultyElev.Requests[floor][elevio.BT_Cab] {
			hallCall := elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
			go DistributeOrder(hallCall, elevOrderTx, elevOrderRx, elevStateRx, elevatorArray, faultyElev, localID)
			faultyElev.Requests[floor][elevio.BT_Cab] = false
		}
	}
}
