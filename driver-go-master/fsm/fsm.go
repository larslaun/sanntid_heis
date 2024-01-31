package fsm

import (
	"Driver-go/elevio"
	"Driver-go/elevator"
	"Driver-go/requests"
	"fmt"
)


func fsm_onRequestButtonPress(buttons chan elevio.ButtonEvent, elevator chan elevator.Elevator){
	select{
		case 
	}
}



func fsm_onInitBetweenFloors(e Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	e.dirn = elevio.MD_Down
	e.behaviour = elevator.EB_Moving
}

