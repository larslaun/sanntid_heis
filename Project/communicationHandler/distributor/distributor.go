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

func DistributeOrder(orderEvent chan elevator.ElevatorOrder, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, distributeElevState chan elevator.Elevator,localID int) {
	elevatorArray := elevator.ElevatorArrayInit()

	for{
		select{
		case newState := <-distributeElevState:
			recievedElevID, _ = strconv.Atoi(newState.ID)
			elevatorArray[recievedElevID] = newState

		case newOrder := <- orderEvent:
			elevOrder := hallAssigner.ChooseOptimalElev(newOrder, *elevatorArray, localID)


			if elevatorArray[localID].NetworkAvailable == false {
				fmt.Print("\nNo network, store order directly\n")
				elevOrderRx <- elevOrder
			} else {
			
				elevOrderTx <- elevOrder

				transmissionFailures := 0


				out:
				for {
					select {
					case receivedState := <-distributeElevState:
						recievedElevID, _ = strconv.Atoi(newState.ID)
						elevatorArray[recievedElevID] = newState
						if receivedState.ID == elevOrder.RecipientID {
							if receivedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] || receivedState.Floor == elevOrder.Order.Floor {
								break out
							} 
						}

					case <-time.After(time.Millisecond * 40):  //Add to settings
						transmissionFailures++

						if transmissionFailures >= settings.MaxTransmissionFailures {
							ReceiverID, _ := strconv.Atoi(elevOrder.RecipientID)
							elevatorArray[ReceiverID].NetworkAvailable = false
							elevOrder = hallAssigner.ChooseOptimalElev(newOrder, *elevatorArray, localID)

							if buttonPress.Button == elevio.BT_Cab {
								elevatorArray[ReceiverID].NetworkAvailable = true
							}
							transmissionFailures = 0
						}
						elevOrderTx <- elevOrder
					}
				} 
			}
	}

}

func RedistributeFaultyElevOrders(elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder, elevStateRx chan elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator, localID int) {
	fmt.Print("\nRedistribute initiated\n")
	faultyElevID, _ := strconv.Atoi(faultyElev.ID)
	shouldRedistribute := false

	for id := 0; id < settings.N_ELEVS; id++{
		if id != localID && elevatorArray[id].NetworkAvailable{
			shouldRedistribute = true
			fmt.Print("\nShould redistribute\n")
		}
	}
	if faultyElevID != localID{
		shouldRedistribute = true
	}

	if shouldRedistribute{
		for floor := 0; floor < settings.N_FLOORS; floor++ {
			for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
				if faultyElev.Requests[floor][btn] {
					hallCall := elevio.ButtonEvent{Floor: floor, Button: btn}
					go DistributeOrder(hallCall, elevOrderTx, elevOrderRx, elevStateRx, elevatorArray, faultyElev, localID)
					
					faultyElev.Requests[floor][btn] = false
					elevatorArray[faultyElevID].Requests[floor][btn] = false
				}
			}
		}
	}
}

func RecoverCabOrders(elevatorArray *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator, localID int) {
	fmt.Print("\nCab recovery initiated\n")

	for floor := 0; floor < settings.N_FLOORS; floor++ {
		if faultyElev.Requests[floor][elevio.BT_Cab] {
			cabCall := elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
			order := elevator.ElevatorOrder{RecipientID: }
			
			distributeOrder

			faultyElev.Requests[floor][elevio.BT_Cab] = false
		}
	}
}
