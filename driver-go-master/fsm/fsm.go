package fsm

import (
	"Driver-go/elevio"
	"driver-go-master/elevator"
	"fmt"
)

func fsm_onInitBetweenFloors(){
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = elevator.EB_Moving
}

func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType){
	fmt.Printf("\n\n%T(%T, %T)", "__FUNCTION__", btn_floor, )
}