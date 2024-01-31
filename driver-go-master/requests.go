package main


type DirnBehaviourPair struct {
	dirn      Dirn
	Behaviour ElevatorBehaviour
}


//checks if there are any requests for the elevator above it's current floor 
//by incrementing through each element in the "boolean" requests matrix.
func requestsAbove(e Elevator) int {
	for f := e.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

//checks if there are any requests for the elevator below it's current floor.
func requestsBelow(e Elevator) int {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

//checks if there are any requests for the elevator at it's current floor 
func requestsHere(e Elevator) int {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return 1
		}
	}
	return 0
}


//decides wether the elevator should move up, stop og move down based on if there are any requests for the elevator.
//if the elevator is already moving up, it will check for requests above it's current floor first and handle them.

//->"Continue in the current direction of travel if there are any further requests in that direction"

func requestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Stop:
		if requestsHere(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}


//checks if the elevator should stop at it's current floor or not. It will only stop if the cab has ordered it to or there is a 
//a request in the direction it is already moving. 
func requests_shouldStop(e Elevator) int {
	switch e.dirn {
	case D_Down:
		return int(e.Requests[e.floor][B_HallDown] || e.Requests[e.floor][B_Cab] || !requestsBelow(e))
	case D_Up:
		return int(e.Requests[e.floor][B_HallUp] || e.Requests[e.floor][B_Cab] || !requestsAbove(e))
	case D_Stop:
		fallthrough
	default:
		return 1
	}
}


