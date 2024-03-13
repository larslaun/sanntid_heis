package requests

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/settings"
)


type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection 
	Behaviour elevator.ElevatorBehaviour
}

// checks if there are any requests for the elevator above it's current floor
// by incrementing through each element in the "boolean" requests matrix.
func RequestsAbove(elev elevator.Elevator) bool {
	for f := elev.Floor + 1; f < settings.N_FLOORS; f++ {
		for btn := 0; btn < settings.N_BUTTONS; btn++ {  
			if elev.Requests[f][btn] { 
				return true
			}
		}
	}
	return false
}

// checks if there are any requests for the elevator below it's current floor.
func RequestsBelow(elev elevator.Elevator) bool {

	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < settings.N_BUTTONS; btn++ {
			if elev.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}

// checks if there are any requests for the elevator at it's current floor
func RequestsHere(elev elevator.Elevator) bool {
	for btn := 0; btn < settings.N_BUTTONS; btn++ {
		if elev.Requests[elev.Floor][btn] {
			return true
		}
	}
	return false
}

//checks if the request matrix for an elevator is empty
func HasRequests(elev elevator.Elevator) bool {
	for floor:=0; floor<settings.N_FLOORS; floor++{
		for btn := 0; btn < settings.N_BUTTONS; btn++ {
			if elev.Requests[floor][btn] {
				return true
			}
		}
	}
	return false
}

//decides wether the elevator should move up, stop or move down based on if there are any requests for the elevator.
//if the elevator is already moving up, it will check for requests above it's current floor first and handle them.

//->"Continue in the current direction of travel if there are any further requests in that direction"

func ChooseDirection(elev elevator.Elevator) DirnBehaviourPair {
	switch elev.Dirn {
	case elevio.MD_Up:
		if RequestsAbove(elev) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_Moving}
		} else if RequestsHere(elev) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_DoorOpen}
		} else if RequestsBelow(elev) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
		}
		
	case elevio.MD_Down:
		if RequestsBelow(elev) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_Moving}
		} else if RequestsHere(elev) {
			return DirnBehaviourPair{elevio.MD_Down, elevator.EB_DoorOpen}
		} else if RequestsAbove(elev) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_Idle}
		}

	case elevio.MD_Stop:
		if RequestsHere(elev) {
			return DirnBehaviourPair{elevio.MD_Stop, elevator.EB_DoorOpen}
		} else if RequestsAbove(elev) {
			return DirnBehaviourPair{elevio.MD_Up, elevator.EB_Moving}
		} else if RequestsBelow(elev) {
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
func ShouldStop(elev elevator.Elevator) bool {
	switch elev.Dirn {
	case elevio.MD_Down:
		return elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] || !RequestsBelow(elev) 

	case elevio.MD_Up:
		return elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] || !RequestsAbove(elev)

	case elevio.MD_Stop:
		fallthrough

	default:
		return true
	}
}

// function where you can spesify a specific request type and it returns wether the request should be cleared or not.
func RequestsShouldClearImmediately(elev elevator.Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return elev.Floor == btnFloor &&
		((elev.Dirn == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(elev.Dirn == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			elev.Dirn == elevio.MD_Stop ||
			btnType == elevio.BT_Cab)
}

// function clears request from the cab at the current floor.
// if the elevator is going up and there are no more requests above or requests UP at the current floor, it will clear the DOWN-request.
// It also clears requests for UP as default, there are either no requests there, or it continues to go UP.
func ClearRequestAtCurrentFloor(elev elevator.Elevator) elevator.Elevator {
	elev.Requests[elev.Floor][elevio.BT_Cab] = false 

	switch elev.Dirn {
	case elevio.MD_Up:
		if !RequestsAbove(elev) && !elev.Requests[elev.Floor][elevio.BT_HallUp] {
			elev.Requests[elev.Floor][elevio.BT_HallDown] = false
		}
		elev.Requests[elev.Floor][elevio.BT_HallUp] = false

	case elevio.MD_Down:
		if !RequestsBelow(elev) && !elev.Requests[elev.Floor][elevio.BT_HallDown] {
			elev.Requests[elev.Floor][elevio.BT_HallUp] = false
		}
		elev.Requests[elev.Floor][elevio.BT_HallDown] = false

	default:
		elev.Requests[elev.Floor][elevio.BT_HallUp] = false
		elev.Requests[elev.Floor][elevio.BT_HallDown] = false
	}

	return elev
}