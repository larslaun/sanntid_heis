package collector

import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	//"fmt"
	//"Elev-project/driver-go-master/fsm"
	"strconv"
	"Elev-project/settings"
)



type ElevatorOrder struct{
	RecipientID string
	Order elevio.ButtonEvent
}




func ElevatorsInit() [settings.NumElevs]elevator.Elevator{
	var elevators = [settings.NumElevs]elevator.Elevator{}

	for i := 0; i < settings.NumElevs; i++ {
		elevator.Elevator_uninitialized(&elevators[i], strconv.Itoa(i))
		elevator.Elevator_print(elevators[i])
	}
	return elevators
}


//Function for collecting states of different elevators. 
//Should change so length of array is not hardcoded. Global var??
func CollectStates(elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator){
	for{
		select {
		case newState := <-elevStateRx:
			elevID, _ := strconv.Atoi(newState.ID)
			elevators[elevID] = newState

			
			
		}
	}
}



//Function for collecting orders broadcasted
//Orders are stored to reciever elevators state.
//Should only recipient store order??

//Maybe scrap this function and handle in fsm instead, then we can use already made Fsm_onRequestButtonPress function
/*
func CollectOrders(elevOrderRx chan ElevatorOrder, elevators *[settings.NumElevs]elevator.Elevator){
	for{
		select{
		case newOrder := <-elevOrderRx:
			fmt.Print("Recieved new order")
			RecipientID, _ := strconv.Atoi(newOrder.RecipientID)
			elevators[RecipientID].Requests[newOrder.Floor][newOrder.Button] = true
		}
	}

}
*/