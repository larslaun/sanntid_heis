package elevator

import "Driver-go/elevio/elevator_io.go"
import "fmt"

type ElevatorBehaviour int

const (
	EB_Idle   ElevatorBehaviour = 0
	EB_DoorOpen                = 1
	EB_Moving                = 2
)


type elevator struct{
	floor int
	dirn MotorDirection
	requests[N_FLOORS][N_BUTTONS] int
	behaviour int
	
	//Config
	doorOpenDuration double
	
}



func eb_toString(eb ElevtorBehaviour) sting{
	if eb == EB_Idle{
		fmt.Print("EB_Idle")
	} else if eb==EB_DoorOpen {
		fmt.Print("EB_DoorOpen")
	}else if eb==EB_Moving {
		fmt.Print("EB_Moving")
	}
}


func elevator_print(es Elevator){
	fmt.Print("  +--------------------+\n")
	fmt.Print(
        "  |floor = %-2d          |\n"
        "  |dirn  = %-12.12s|\n"
        "  |behav = %-12.12s|\n",
        es.floor,
        elevio_dirn_toString(es.dirn),
        eb_toString(es.behaviour))
	fmt.Print("  +--------------------+\n");
	fmt.Print("  |  | up  | dn  | cab |\n");
	for(int f = N_FLOORS-1; f >= 0; f--){
		fmt.Print("  | %d", f);
		for btn := 0; btn < N_BUTTONS; btn++{
			if (f == N_FLOORS-1 && btn == B_HallUp)  || 
				(f == 0 && btn == B_HallDown) 
			{
				fmt.Print("|     ");
			} else {
				fmt.Print(es.requests[f][btn] ? "|  #  " : "|  -  ");
			}
		}
		fmt.Print("|\n");
	}
	fmt.Print("  +--------------------+\n");
	
}


func elevator_uninitialized(es *elevator){  //initialize elevator, passing pointer
	es.floor = -1
	es.dirn = MD_Stop
	es.behaviour = EB_Idle
	es.doorOpenDuration = 3.0
}