package main


type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
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


//function where you can spesify a specific request type and it returns wether the request should be cleared or not. 
func requestsShouldClearImmediately(e Elevator, btnFloor int, btnType Button) int {
	return  e.floor == btnFloor && 
			(
			(e.dirn == D_Up && btnType == B_HallUp) || 
			(e.dirn == D_Down && btnType == B_HallDown) || 
			e.dirn == D_Stop || 
			btnType == B_Cab
			)	
}



//function clears request from the cab at the current floor. 
//if the elevator is going up and there are no more requests above or requests UP at the current floor, it will clear the DOWN-request. 
//It also clears requests for UP as default, there are either no requests there, or it continues to go UP. 
//If the elevetor state is Stop, it clears both UP and DOWN hall calls, probably only one of them at that floor, since the elevator door will open for one of (the first) the requests. 
func requestsClearAtCurrentFloor(e Elevator) Elevator {
	e.requests[e.floor][B_Cab] = 0

	switch e.dirn {
	case D_Up:
		if !requestsAbove(e) && !e.requests[e.floor][B_HallUp] {
			e.requests[e.floor][B_HallDown] = 0
		}
		e.requests[e.floor][B_HallUp] = 0
	case D_Down:
		if !requestsBelow(e) && !e.requests[e.floor][B_HallDown] {
			e.requests[e.floor][B_HallUp] = 0
		}
		e.requests[e.floor][B_HallDown] = 0
	case D_Stop, default:
		e.requests[e.floor][B_HallUp] = 0
		e.requests[e.floor][B_HallDown] = 0
	}

	return e
}