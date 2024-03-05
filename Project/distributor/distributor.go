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
		elevStateTx <- *localElev
		time.Sleep(1000 * time.Millisecond)
	}
}





//psuedo distributor
//Recieves buttonpress, then calculates optimal elevator wiht cost func,then sends elevOrder which includes order and ID of elev.


func DistributeOrder(buttonPress chan elevio.ButtonEvent,  elevOrderTx chan collector.ElevatorOrder,  elevators [settings.NumElevs]elevator.Elevator){
	for{
		select{
		case buttonPress:=<-buttonPress:
			elevOrder := hallAssigner.ChooseOptimalElev(buttonPress, elevators) //choose optimalelev must calculat cost func for all elevs and create order to optimal elevator
			fmt.Printf("\nOptimal elev calculated:\n")
			fmt.Printf("optimalElevID: " + elevOrder.RecipientID + "\n")
			fmt.Printf("Floor: %d \n", elevOrder.Order.Floor)
			fmt.Printf("Button: %d \n", elevOrder.Order.Button)
			
			elevOrderTx<-elevOrder
		}
	}
}
