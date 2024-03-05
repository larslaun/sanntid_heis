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

/*
func Redistribute(elevStateTx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, FaultyElevID string){
	FaultyElev elevator.Elevator = elevators[FaultyElevID]
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := 0; btn < elevator.N_BUTTONS - 1; btn++ {  
			if FaultyElev.requests[floor][btn] == true{
				HallCall elevio.ButtonEvent := elevio.ButtonEvent{Floor: floor, Button: btn}

				// DistributeOrder(Hallcall, elevators) ??
			}
			}
		}
	}
}
*/



//psuedo distributor
//Receives buttonpress, then calculates optimal elevator wiht cost func,then sends elevOrder which includes order and ID of elev.


func DistributeOrder(buttonPress chan elevio.ButtonEvent,  elevOrderTx chan collector.ElevatorOrder,  elevators *[settings.NumElevs]elevator.Elevator){
	for{
		select{
		case buttonPress:=<-buttonPress:

			//Problem her sannsynligvis. Får ikke tak i heis states
			elevator.Elevator_print(elevators[0])
		
			elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators) //choose optimalelev must calculat cost func for all elevs and create order to optimal elevator
			
			fmt.Printf("\nOptimal elev calculated:\n")
			fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
			fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
			fmt.Printf("Button: %d \n", elevOrder.Order.Button)
				
			elevOrderTx<-elevOrder
		}
	}
}
