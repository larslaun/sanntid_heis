package distributor

import (
	"Elev-project/collector"
	"Elev-project/hallAssigner"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/settings"


	"time"
	"fmt"
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


func DistributeOrder(buttonPress elevio.ButtonEvent,  elevOrderTx chan collector.ElevatorOrder,  elevators *[settings.NumElevs]elevator.Elevator){
	/*for{
		select{
		case buttonPress:=<-buttonPress:

			if buttonPress.Button != elevio.BT_Cab {*/
				//Problem her sannsynligvis. Får ikke tak i heis states
				//elevator.Elevator_print(elevators[0])

			
				elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators) //choose optimalelev must calculat cost func for all elevs and create order to optimal elevator
				
				/*
				fmt.Printf("\nOptimal elev calculated:\n")
				fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
				fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
				fmt.Printf("Button: %d \n", elevOrder.Order.Button)
				*/
				
				elevOrderTx<-elevOrder

				
			//}
		//}
	//}
}


func RedistributeFaultyElevOrders(elevOrderTx chan collector.ElevatorOrder, elevators *[settings.NumElevs]elevator.Elevator, faultyElev *elevator.Elevator, redistributeSignal chan bool){
	for{
		select{
		case <-redistributeSignal:
			fmt.Print("\nRedistribute initiated\n")
			for floor := 0; floor < elevator.N_FLOORS; floor++ {
				for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {  
					if faultyElev.Requests[floor][btn]{
						var hallCall elevio.ButtonEvent = elevio.ButtonEvent{Floor: floor, Button: btn}
						DistributeOrder(hallCall, elevOrderTx, elevators)	
						faultyElev.Requests[floor][btn] = false
					}
				}
			}
		}
	}
}