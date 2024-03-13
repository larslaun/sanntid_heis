package cost

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/elevatorDriver/requests"
	"Elev-project/settings"
)


// Calculates an estimate for the time an elevator takes to go from a list of requests until they are executed.
func TimeToIdle(elevSim elevator.Elevator) int {

	var duration int = 0
	clearForSafety(&elevSim, &duration)

	switch elevSim.Behaviour {
	case elevator.EB_Idle:
		elevSim.Dirn = requests.ChooseDirection(elevSim).Dirn
		if elevSim.Dirn == elevio.MD_Stop {
			return duration
		}
	case elevator.EB_Moving:
		duration += settings.TRAVELTIME
		elevSim.Floor += int(elevSim.Dirn)

	case elevator.EB_DoorOpen:
		duration += settings.DOOROPENTIME
	}

	for {
		if requests.ShouldStop(elevSim) {
			elevSim = costClearAtCurrentFloor(elevSim)
			duration += settings.DOOROPENTIME
			elevSim.Dirn = requests.ChooseDirection(elevSim).Dirn
			if elevSim.Dirn == elevio.MD_Stop {
				return duration
			}
		}
		elevSim.Floor += int(elevSim.Dirn)
		duration += settings.TRAVELTIME
	}
}

// Augmented clear at ClearAtCurrentFloor function, so that it does not affect the real elevator. 
func costClearAtCurrentFloor(elevOld elevator.Elevator) elevator.Elevator {
	var elev elevator.Elevator = elevOld

	for btn := 0; btn < settings.N_BUTTONS; btn++ {
		if elev.Requests[elev.Floor][btn] {
			elev.Requests[elev.Floor][btn] = false
		}
	}
	return elev
}

// Clears all requests on the elevators current floor, and adds one DOOROPENTIME to the estimated runtime.
func clearForSafety(elev *elevator.Elevator, cost *int) {
	for floor := 0; floor < settings.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++ {
			if requests.RequestsShouldClearImmediately(*elev, floor, btn) && (elev.Requests[floor][btn]) {
				elev.Requests[floor][btn] = false
				*cost += settings.DOOROPENTIME
			}
		}
	}
}
