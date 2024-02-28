package cost_function

import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/requests"
	"fmt"
	"time"
)

const doorOpenTime = 3
const travelTime = 5

const N_FLOORS int = 4
const N_BUTTONS int = 3

func TimeToIdle(elevSim elevator.Elevator) int {

	var Duration int = 0

	switch elevSim.Behaviour {
	case elevator.EB_Idle:
		elevSim.Dirn = requests.RequestsChooseDirection(elevSim).Dirn
		if elevSim.Dirn == elevio.MD_Stop {
			return Duration
		}
	case elevator.EB_Moving:
		Duration += travelTime
		elevSim.Floor += int(elevSim.Dirn)
	case elevator.EB_DoorOpen:
		Duration += doorOpenTime

	}
	i:=0
	for {
		fmt.Printf("%d\n" , i)
		elevator.Elevator_print(elevSim)
		if requests.Requests_shouldStop(elevSim) {
			
			elevSim = requests.RequestsClearAtCurrentFloor(elevSim)
			Duration += doorOpenTime
			elevSim.Dirn = requests.RequestsChooseDirection(elevSim).Dirn
			if elevSim.Dirn == elevio.MD_Stop {
				return Duration
			}
		}
		i++
		time.Sleep(500 * time.Millisecond)
	}
}

func CostClearAtCurrentFloor(elevOld elevator.Elevator) elevator.Elevator {
	var elev elevator.Elevator = elevOld

	for btn := 0; btn < N_BUTTONS; btn++{
		if elev.Requests[elev.Floor][btn]{
			elev.Requests[elev.Floor][btn] = false
		}
	}

	return elev
}


