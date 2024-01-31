package elevator

import "Driver-go/elevio"
import "fmt"

type ElevatorBehaviour int

const (
	EB_Idle   ElevatorBehaviour = 0
	EB_DoorOpen                = 1
	EB_Moving                = 2
)


type Elevator struct{
	floor int
	dirn elevio.MotorDirection
	requests[N_FLOORS][N_BUTTONS] int
	behaviour int
	
	//Config
	doorOpenDuration float64
	
}

func elevio_dirn_toString(md elevio.MotorDirection){
	if md == elevio.MD_Up{
		fmt.Print("MD_Up")
	} else if md==MD_Down {
		fmt.Print("MD_Down")
	}else if md==MD_Stop {
		fmt.Print("MD_Stop")
	}
}

func eb_toString(eb ElevatorBehaviour) string{
	if eb == EB_Idle{
		return "EB_Idle"  // return or print directly
	} else if eb==EB_DoorOpen {
		return "EB_DoorOpen"
	}else if eb==EB_Moving {
		return "EB_Moving"
	}
}


func elevator_print(es elevator){
	fmt.Print("  +--------------------+\n")
	fmt.Printf("  |floor = %-2d|\n", es.floor)
    fmt.Printf("  |dirn  = %-12.12s|\n", elevio_dirn_toString(es.dirn))
    fmt.Print("  |behav = %-12.12s|\n", eb_toString(es.behaviour)) 
	fmt.Print("  +--------------------+\n")
	fmt.Print("  |  | up  | dn  | cab |\n")
	for f := N_FLOORS-1; f >= 0; f--{
		fmt.Print("  | %d", f);
		for btn := 0; btn < N_BUTTONS; btn++{
			if (f == N_FLOORS-1 && btn == B_HallUp)  || (f == 0 && btn == B_HallDown){
				fmt.Print("|     ");
			} else {
				if es.requests[f][btn] == 1{
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
				//fmt.Print(es.requests[f][btn] ? "|  #  " : "|  -  "); replaced by if sentence over ^
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