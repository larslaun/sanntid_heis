package elevator

import "Driver-go/elevio/elevator_io.go"
import "fmt"

type ElevatorBehavour int

const (
	EB_Idle   ElevatorBehavour = 0
	EB_DoorOpen                = 1
	EB_Moving                = 2
)


type elevator struct{
	floor int
	dirn MotorDirection
	requests[N_FLOORS][N_BUTTONS] int
	behaviour int

	//Mulig dette kan droppes
	type config struct{
		clearReaquestVariant ClearRequestVariantconst
		doorOpenDuration double
	}
}



func eb_toString(eb ElevtorBehaviour) sting{
	if eb == EB_Idle{
		fmt.Print("EB_Idle")
	}
}

