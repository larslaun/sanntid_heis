package fsm

import (
	"Driver-go/elevio"
	"Driver-go/elevator"
	"Driver-go/requests"
	"fmt"
)

func fsm_onInitBetweenFloors(e Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	e.dirn = elevio.MD_Down
	e.behaviour = elevator.EB_Moving
}

func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType){
	fmt.Printf("\n\n%T(%T, %T)", "__FUNCTION__", btn_floor, )

	
}