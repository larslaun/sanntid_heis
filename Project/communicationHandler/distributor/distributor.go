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
			receivedElevID, _ := strconv.Atoi(newState.ID)
			elevatorArray[receivedElevID] = newState

		case newOrder := <- orderEvent:
			//elevator.PrintElevator(elevatorArray[localID])
			//fmt.Print("Sending new order: ")
			//fmt.Printf("%+v\n", newOrder.Order)
			//pre 

			elevOrder := hallAssigner.ChooseOptimalElev(newOrder, elevatorArray, localID)

			fmt.Printf("Sending to elev ID: %s\n" ,elevOrder.RecipientID)

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
						receivedElevID, _ := strconv.Atoi(receivedState.ID)
						elevatorArray[receivedElevID] = receivedState
						if receivedState.ID == elevOrder.RecipientID {
							if receivedState.Requests[elevOrder.Order.Floor][elevOrder.Order.Button] || receivedState.Floor == elevOrder.Order.Floor {
								break out
							} 
						}

					case <-time.After(settings.TRANSMISSION_RATE):  
						transmissionFailures++
						fmt.Printf("transmission fails: %d\n", transmissionFailures)

						if transmissionFailures >= settings.MaxTransmissionFailures {
							ReceiverID, _ := strconv.Atoi(elevOrder.RecipientID)
							elevatorArray[ReceiverID].NetworkAvailable = false
							elevOrder = hallAssigner.ChooseOptimalElev(newOrder, elevatorArray, localID)

							if elevOrder.Order.Button == elevio.BT_Cab {
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
}

func RedistributeFaultyElevOrders(orderEvent chan elevator.ElevatorOrder, elevatorArray *[settings.N_ELEVS]elevator.Elevator, faultyElev *elevator.Elevator, localID int, distributeElevState chan elevator.Elevator) {
	fmt.Print("\nRedistribute initiated\n")
	faultyElevID, _ := strconv.Atoi(faultyElev.ID)
	shouldRedistribute := false

	distributeElevState <- *faultyElev

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
					
					faultyElev.Requests[floor][btn] = false
					elevatorArray[faultyElevID].Requests[floor][btn] = false
					distributeElevState <- *faultyElev

					hallCall := elevio.ButtonEvent{Floor: floor, Button: btn}
					order := elevator.ElevatorOrder{RecipientID: faultyElev.ID, Order: hallCall}
					
					orderEvent <- order	
				}
			}
		}
	}
}

func RecoverCabOrders(orderEvent chan elevator.ElevatorOrder, distributeElevState chan elevator.Elevator,faultyElev *elevator.Elevator) {
	fmt.Print("\nCab recovery initiated\n")

	for floor := 0; floor < settings.N_FLOORS; floor++ {
		if faultyElev.Requests[floor][elevio.BT_Cab] {
			faultyElev.Requests[floor][elevio.BT_Cab] = false
			distributeElevState <- *faultyElev

			cabCall := elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
			order := elevator.ElevatorOrder{RecipientID: faultyElev.ID, Order: cabCall}
			orderEvent <- order
		}
	}
}
