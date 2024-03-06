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

func HallMatrix(elevators *[settings.NumElevs]elevator.Elevator) [elevator.N_FLOORS][elevator.N_BUTTONS - 1] bool {
	//lager en matrise med nuller
	hallMatrix := make([][]bool, elevator.N_FLOORS)
    for i := range hallMatrix {
        hallMatrix[i] = make([]bool, elevator.N_BUTTONS - 1) //tar bare med hall_requests
    } 

	//går gjennom hvert Hall-element i hver heis sin matrise og OR'er med hvert element i resultMatrix
	for id := 0; id < len(elevators); id++ {
		for floor := 0; floor < elevator.N_FLOORS; floor++ {
			for btn := elevio.BT_HallUp; btn < elevio.BT_HallDown; btn ++{
				hallMatrix[floor][btn] = hallMatrix[floor][btn] || elevators[id].Requests[floor][btn]
			}
		}
	}

	return hallMatrix
}

//psuedo distributor
//Receives buttonpress, then calculates optimal elevator wiht cost func,then sends elevOrder which includes order and ID of elev.


func DistributeOrder(buttonPress chan elevio.ButtonEvent,  elevOrderTx chan collector.ElevatorOrder,  elevators *[settings.NumElevs]elevator.Elevator){
	for{
		select{
		case buttonPress:=<-buttonPress:

			//Problem her sannsynligvis. Får ikke tak i heis states
			elevator.Elevator_print(elevators[0])
		
			elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators) //choose optimalelev must calculate cost func for all elevs and create order to optimal elevator
			
			fmt.Printf("\nOptimal elev calculated:\n")
			fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
			fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
			fmt.Printf("Button: %d \n", elevOrder.Order.Button)
				
			elevOrderTx<-elevOrder
		}
	}
}


func RedistributeOrders(elevOrderTx chan collector.ElevatorOrder, elevators *[settings.NumElevs]elevator.Elevator, faultyElev elevator.Elevator){
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevator.N_BUTTONS - 1; btn++ {   //-1 to skip cab buttons
			if faultyElev.requests[floor][btn] == true{
				hallCall elevio.ButtonEvent := elevio.ButtonEvent{Floor: floor, Button: btn}
				DistributeOrder(hallCall, elevOrderTx, elevators)	
			}
		}
	}
}