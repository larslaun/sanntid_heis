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
		localElev.Available = true
		elevStateTx <- *localElev
		time.Sleep(1000 * time.Millisecond)
	}
}



//psuedo distributor
//Receives buttonpress, then calculates optimal elevator wiht cost func,then sends elevOrder which includes order and ID of elev.


func DistributeOrder(buttonPress elevio.ButtonEvent,  elevOrderTx chan collector.ElevatorOrder,  elevators *[settings.NumElevs]elevator.Elevator){
	/*for{
		select{
		case buttonPress:=<-buttonPress:

<<<<<<< HEAD
			//Problem her sannsynligvis. Får ikke tak i heis states
			elevator.Elevator_print(elevators[1])
		
			elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators) //choose optimalelev must calculat cost func for all elevs and create order to optimal elevator
=======
			if buttonPress.Button != elevio.BT_Cab {*/
				//Problem her sannsynligvis. Får ikke tak i heis states
				//elevator.Elevator_print(elevators[0])

>>>>>>> 2ea464e6d2b2f4118d1e6c63a1735f2def87e42c
			
				elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators) //choose optimalelev must calculat cost func for all elevs and create order to optimal elevator
				
				fmt.Printf("\nOptimal elev calculated:\n")
				fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
				fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
				fmt.Printf("Button: %d \n", elevOrder.Order.Button)
					
				elevOrderTx<-elevOrder
			//}
		//}
	//}
}


func RedistributeFaultyElevOrders(elevOrderTx chan collector.ElevatorOrder, elevators *[settings.NumElevs]elevator.Elevator, faultyElev elevator.Elevator){
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := 0; btn < elevator.N_BUTTONS - 1; btn++ {   //-1 to skip cab buttons
			if faultyElev.requests[floor][btn] == true{
				hallCall elevio.ButtonEvent := elevio.ButtonEvent{Floor: floor, Button: btn}
				DistributeOrder(hallCall, elevOrderTx, elevators)	
			}
		}
	}
}