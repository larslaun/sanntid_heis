package cost_function

import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/requests"
	"fmt"
)

const DOOROPENTIME = 3
const TRAVELTIME = 5

const N_FLOORS int = 4
const N_BUTTONS int = 3

//Calculates an estimate for the time an elevator takes to go from a list of requests until they are executed.
func TimeToIdle(elevSim elevator.Elevator) int {

	var Duration int = 0

	ClearForSafety(&elevSim, &Duration)
	
	switch elevSim.Behaviour {
	case elevator.EB_Idle:
		elevSim.Dirn = requests.RequestsChooseDirection(elevSim).Dirn
		if elevSim.Dirn == elevio.MD_Stop {
			fmt.Printf("\nElevator ID " + elevSim.ID + " calculated duration: %d\n", Duration)
			return Duration
		}
	case elevator.EB_Moving:
		Duration += TRAVELTIME
		elevSim.Floor += int(elevSim.Dirn)
	case elevator.EB_DoorOpen:
		Duration += DOOROPENTIME

	}
	for {
		if requests.Requests_shouldStop(elevSim) {
			
			elevSim = CostClearAtCurrentFloor(elevSim)
			Duration += DOOROPENTIME
			elevSim.Dirn = requests.RequestsChooseDirection(elevSim).Dirn
			if elevSim.Dirn == elevio.MD_Stop {
				fmt.Printf("\nElevator ID " + elevSim.ID + " calculated duration: %d\n", Duration)
				return Duration
			}
		}
		elevSim.Floor += int(elevSim.Dirn)
		Duration += TRAVELTIME
	}
}

//Augmented clear at ClearAtCurrentFloor function, such that it does not affect the real elevator. Takes an elevator 
// of type Elevator as an argument, and returns 
func CostClearAtCurrentFloor(elevOld elevator.Elevator) elevator.Elevator {
	var elev elevator.Elevator = elevOld

	for btn := 0; btn < N_BUTTONS; btn++{
		if elev.Requests[elev.Floor][btn]{
			elev.Requests[elev.Floor][btn] = false
		}
	}

	return elev
}

//Clears all requests on the elevators current floor, and adds one DOOROPENTIME to the estimated runtime. Takes two pointers
// Elevator e and the cost as arguments.
func ClearForSafety(e *elevator.Elevator, cost *int){
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++{
			if requests.RequestsShouldClearImmediately(*e, floor, btn){
				e.Requests[floor][btn] = false
				*cost += DOOROPENTIME
			}
		}
	}
}
