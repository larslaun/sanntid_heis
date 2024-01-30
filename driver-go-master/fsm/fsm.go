package fsm

import (
	"Driver-go/elevio"
	"driver-go_master/elevator"
	"driver-go_master/elevio"
)

func fsm_onInitBetweenFloors(){
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = elevator.EB_Moving
}