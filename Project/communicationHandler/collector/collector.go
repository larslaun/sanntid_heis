package collector

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/settings"
	"strconv"
)

func CollectStates(elevStateRx chan elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator, distributeElevState chan elevator.Elevator){
	for{
		select {
		case newState := <-elevStateRx:
			elevID, _ := strconv.Atoi(newState.ID)
			elevatorArray[elevID] = newState
			distributeElevState<-newState
		}
	}
}



