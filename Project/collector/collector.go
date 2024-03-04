package collector

import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"fmt"
	//"Elev-project/driver-go-master/fsm"
	"strconv"
)



type ElevatorOrder struct{
	RecipientID string
	order elevio.ButtonEvent
}




func ElevatorsInit(numElevs int) [3]elevator.Elevator{
	var elevators = [3]elevator.Elevator{}

	for i := 0; i < numElevs; i++ {
		elevator.Elevator_uninitialized(&elevators[i], strconv.Itoa(i))
		elevator.Elevator_print(elevators[i])
	}
	return elevators
}


//Function for collecting states of different elevators. 
//Should change so length of array is not hardcoded. Global var??
func CollectStates(elevStateRx chan elevator.Elevator, elevators *[3]elevator.Elevator){
	for{
		select {
		case newState := <-elevStateRx:
			fmt.Print("Recieved new state:")
			elevID, _ := strconv.Atoi(newState.ID)
			elevators[elevID] = newState
			elevator.Elevator_print(elevators[elevID])
		}
	}
}



//Function for collecting orders broadcasted
//Orders are stored to reciever elevators state.
//Should only recipient store order??
func CollectOrders(elevOrderRx chan ElevatorOrder, elevators *[3]elevator.Elevator){
	for{
		select{
		case newOrder := <-elevOrderRx:
			fmt.Print("Recieved new order")
			RecipientID, _ := strconv.Atoi(newOrder.RecipientID)
			elevators[RecipientID].Requests[newOrder.Floor][newOrder.Button] = true
		}
	}

}
