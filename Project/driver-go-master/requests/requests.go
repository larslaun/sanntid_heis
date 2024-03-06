package requests

import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
)


type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection //endret type fra Dirn til elevio.MotorDirection
	Behaviour elevator.ElevatorBehaviour
}

// checks if there are any requests for the elevator above it's current floor
// by incrementing through each element in the "boolean" requests matrix.
func RequestsAbove(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < elevator.N_FLOORS; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {  //må btn loop endres til å iterere gjennom enum, ref
			if e.Requests[f][btn] { //la til ==1 på disse NOT. fjernet igjen pga. endret request til bool
				return true
			}
		}
	}
	return false
}

// checks if there are any requests for the elevator below it's current floor.
func RequestsBelow(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

// checks if there are any requests for the elevator at it's current floor
func RequestsHere(e elevator.Elevator) bool {
	for btn := 0; btn < elevator.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

//checks if the request matrix for an elevator is empty
func HasRequests(e elevator.Elevator) bool {
	return (RequestsBelow(e) || RequestsAbove(e) || RequestsHere(e))
}

//decides wether the elevator should move up, stop og move down based on if there are any requests for the elevator.
//if the elevator is already moving up, it will check for requests above it's current floor first and handle them.

//->"Continue in the current direction of travel if there are any further requests in that direction"

func RequestsChooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevio.MD_Up:
		if RequestsAbove(e) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_DoorOpen}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
		}
	case elevio.MD_Down:
		if RequestsBelow(e) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_Moving}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
		}
	case elevio.MD_Stop:
		if RequestsHere(e) {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_DoorOpen}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_Moving}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_Moving} 
		}else {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
		}
	default:
		return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
	}
}

// checks if the elevator should stop at it's current floor or not. It will only stop if the cab has ordered it to or there is a
// a request in the direction it is already moving.
func Requests_shouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case elevio.MD_Down:
		return e.Requests[e.Floor][elevio.BT_HallDown] || e.Requests[e.Floor][elevio.BT_Cab] || !RequestsBelow(e) //fjernet int conversion
	case elevio.MD_Up:
		return e.Requests[e.Floor][elevio.BT_HallUp] || e.Requests[e.Floor][elevio.BT_Cab] || !RequestsAbove(e)
	case elevio.MD_Stop:
		fallthrough
	default:
		return true
	}
}

// function where you can spesify a specific request type and it returns wether the request should be cleared or not.
func RequestsShouldClearImmediately(e elevator.Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return e.Floor == btnFloor &&
		((e.Dirn == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.Dirn == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			e.Dirn == elevio.MD_Stop ||
			btnType == elevio.BT_Cab)
}

// function clears request from the cab at the current floor.
// if the elevator is going up and there are no more requests above or requests UP at the current floor, it will clear the DOWN-request.
// It also clears requests for UP as default, there are either no requests there, or it continues to go UP.
// If the elevetor state is Stop, it clears both UP and DOWN hall calls, probably only one of them at that floor, since the elevator door will open for one of (the first) the requests.
func RequestsClearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	e.Requests[e.Floor][elevio.BT_Cab] = false //endret fra 0 til false

	switch e.Dirn {
	case elevio.MD_Up:
		if !RequestsAbove(e) && !e.Requests[e.Floor][elevio.BT_HallUp] {
			e.Requests[e.Floor][elevio.BT_HallDown] = false
		}
		e.Requests[e.Floor][elevio.BT_HallUp] = false
	case elevio.MD_Down:
		if !RequestsBelow(e) && !e.Requests[e.Floor][elevio.BT_HallDown] {
			e.Requests[e.Floor][elevio.BT_HallUp] = false
		}
		e.Requests[e.Floor][elevio.BT_HallDown] = false
	default:
		e.Requests[e.Floor][elevio.BT_HallUp] = false
		e.Requests[e.Floor][elevio.BT_HallDown] = false
	}

	return e
}