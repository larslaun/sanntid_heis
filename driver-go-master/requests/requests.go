package requests

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"fmt"
)


type DirnBehaviourPair struct {
	Dirn      Dirn
	Behaviour ElevatorBehaviour
}

// checks if there are any requests for the elevator above it's current floor
// by incrementing through each element in the "boolean" requests matrix.
func RequestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

// checks if there are any requests for the elevator below it's current floor.
func RequestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

// checks if there are any requests for the elevator at it's current floor
func RequestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return 1
		}
	}
	return 0
}

//decides wether the elevator should move up, stop og move down based on if there are any requests for the elevator.
//if the elevator is already moving up, it will check for requests above it's current floor first and handle them.

//->"Continue in the current direction of travel if there are any further requests in that direction"

func RequestsChooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case D_Up:
		if RequestsAbove(e) {
			return DirnBehaviourPair{MD_Up, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{MD_Up, EB_DoorOpen}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{MD_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{MD_Stop, EB_Idle}
		}
	case D_Down:
		if RequestsBelow(e) {
			return DirnBehaviourPair{MD_Down, EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{MD_Down, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{MD_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{MD_Stop, EB_Idle}
		}
	case D_Stop:
		if RequestsHere(e) {
			return DirnBehaviourPair{MD_Stop, EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{MD_Up, EB_Moving}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{MD_Down, EB_Moving} 
		}else {
			return DirnBehaviourPair{MD_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{MD_Stop, EB_Idle}
	}
}

// checks if the elevator should stop at it's current floor or not. It will only stop if the cab has ordered it to or there is a
// a request in the direction it is already moving.
func Requests_shouldStop(e Elevator) bool {
	switch e.Dirn {
	case MD_Down:
		return int(e.Requests[e.Floor][BT_HallDown] || e.Requests[e.Floor][BT_Cab] || !RequestsBelow(e))
	case MD_Up:
		return int(e.Requests[e.Floor][BT_HallUp] || e.Requests[e.Floor][BT_Cab] || !RequestsAbove(e))
	case MD_Stop:
		fallthrough
	default:
		return 1
	}
}

// function where you can spesify a specific request type and it returns wether the request should be cleared or not.
func RequestsShouldClearImmediately(e Elevator, btnFloor int, btnType Button) bool {
	return e.Floor == btnFloor &&
		((e.Dirn == MD_Up && btnType == BT_HallUp) ||
			(e.Dirn == MD_Down && btnType == BT_HallDown) ||
			e.Dirn == MD_Stop ||
			btnType == BT_Cab)
}

// function clears request from the cab at the current floor.
// if the elevator is going up and there are no more requests above or requests UP at the current floor, it will clear the DOWN-request.
// It also clears requests for UP as default, there are either no requests there, or it continues to go UP.
// If the elevetor state is Stop, it clears both UP and DOWN hall calls, probably only one of them at that floor, since the elevator door will open for one of (the first) the requests.
func RequestsClearAtCurrentFloor(e Elevator) Elevator {
	e.Requests[e.Floor][B_Cab] = 0

	switch e.Dirn {
	case MD_Up:
		if !RequestsAbove(e) && !e.Requests[e.Floor][BT_HallUp] {
			e.Requests[e.Floor][BT_HallDown] = 0
		}
		e.Requests[e.Floor][BT_HallUp] = 0
	case MD_Down:
		if !RequestsBelow(e) && !e.Requests[e.Floor][BT_HallDown] {
			e.Requests[e.Floor][BT_HallUp] = 0
		}
		e.Requests[e.Floor][BT_HallDown] = 0
	default:
		e.Requests[e.Floor][BT_HallUp] = 0
		e.Requests[e.Floor][BT_HallDown] = 0
	}

	return e
}