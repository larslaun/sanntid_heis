package collector

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/settings"
	"strconv"
)



//Function for collecting states of different elevators. 
//Should change so length of array is not hardcoded. Global var??
func CollectStates(elevStateRx chan elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator){
	for{
		select {
		case newState := <-elevStateRx:
			elevID, _ := strconv.Atoi(newState.ID)
			elevators[elevID] = newState
		}

	}
}



