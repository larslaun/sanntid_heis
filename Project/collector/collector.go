package collector

import (
	"Elev-project/driver-go-master/elevator"
	"fmt"
	//"Elev-project/driver-go-master/fsm"
	"strconv"
)

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
func CollectStates(elevStateRx chan elevator.Elevator, elevators *[3]elevator.Elevator) [3]elevator.Elevator{
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
