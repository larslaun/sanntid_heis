package cost

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/elevatorDriver/requests"
	"Elev-project/settings"
)


// Calculates an estimate for the time an elevator takes to go from a list of requests until they are executed.
func TimeToIdle(elevSim elevator.Elevator) int {

	var Duration int = 0

	ClearForSafety(&elevSim, &Duration)
	//fmt.Printf("\ncost print 1: %d\n", Duration)

	switch elevSim.Behaviour {
	case elevator.EB_Idle:
		elevSim.Dirn = requests.RequestsChooseDirection(elevSim).Dirn
		if elevSim.Dirn == elevio.MD_Stop {
			//fmt.Printf("\nElevator ID "+elevSim.ID+" calculated duration1: %d\n", Duration)
			//elevator.Elevator_print(elevSim)
			return Duration
		}
	case elevator.EB_Moving:
		Duration += settings.TRAVELTIME
		elevSim.Floor += int(elevSim.Dirn)
	case elevator.EB_DoorOpen:
		Duration += settings.DOOROPENTIME

	}
	for {
		if requests.Requests_shouldStop(elevSim) {

			elevSim = CostClearAtCurrentFloor(elevSim)
			Duration += settings.DOOROPENTIME
			elevSim.Dirn = requests.RequestsChooseDirection(elevSim).Dirn
			if elevSim.Dirn == elevio.MD_Stop {
				//fmt.Printf("\nElevator ID "+elevSim.ID+" calculated duration2: %d\n", Duration)
				//elevator.Elevator_print(elevSim)
				return Duration
			}
		}
		elevSim.Floor += int(elevSim.Dirn)
		Duration += settings.TRAVELTIME
	}
}

// Augmented clear at ClearAtCurrentFloor function, such that it does not affect the real elevator. Takes an elevator
// of type Elevator as an argument, and returns
func CostClearAtCurrentFloor(elevOld elevator.Elevator) elevator.Elevator {
	var elev elevator.Elevator = elevOld

	for btn := 0; btn < settings.N_BUTTONS; btn++ {
		if elev.Requests[elev.Floor][btn] {
			elev.Requests[elev.Floor][btn] = false
		}
	}

	return elev
}

// Clears all requests on the elevators current floor, and adds one DOOROPENTIME to the estimated runtime. Takes two pointers
// Elevator e and the cost as arguments.
func ClearForSafety(e *elevator.Elevator, cost *int) {
	for floor := 0; floor < settings.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++ {
			if requests.RequestsShouldClearImmediately(*e, floor, btn) && (e.Requests[floor][btn]) {
				e.Requests[floor][btn] = false
				*cost += settings.DOOROPENTIME
			}
		}
	}
}
