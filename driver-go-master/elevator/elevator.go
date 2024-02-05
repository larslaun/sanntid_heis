package elevator

import "Driver-go/elevio"
import "fmt"


//Deklarerer her foreløpig
const N_FLOORS int = 4
const N_BUTTONS int = 3


type ElevatorBehaviour int

const (
	EB_Idle   ElevatorBehaviour = 0
	EB_DoorOpen                = 1
	EB_Moving                = 2
)


type Elevator struct{
	Floor int
	Dirn elevio.MotorDirection
	Requests[N_FLOORS][N_BUTTONS] int
	Behaviour ElevatorBehaviour
	
	//Config
	DoorOpenDuration float64
	
}


// return or print directly??
func Elevio_dirn_toString(md elevio.MotorDirection) string{
	if md == elevio.MD_Up{
		return "MD_Up"
	} else if md==elevio.MD_Down {
		return "MD_Down"
	}else{
		return "MD_Stop"
	}
}

func Eb_toString(eb ElevatorBehaviour) string{
	if eb == EB_Idle{
		return "EB_Idle"  
	} else if eb==EB_DoorOpen {
		return "EB_DoorOpen"
	}else{
		return "EB_Moving"
	} 
}


func Elevator_print(es Elevator){
	fmt.Print("  +--------------------+\n")
	fmt.Printf("  |floor = %-2d|\n", es.Floor)
    fmt.Printf("  |dirn  = %-12.12s|\n", Elevio_dirn_toString(es.Dirn))
    fmt.Print("  |behav = %-12.12s|\n", Eb_toString(es.Behaviour)) 
	fmt.Print("  +--------------------+\n")
	fmt.Print("  |  | up  | dn  | cab |\n")
	for f := N_FLOORS-1; f >= 0; f--{
		fmt.Print("  | %d", f);
		for btn := 0; btn < N_BUTTONS; btn++{
			if (f == N_FLOORS-1 && btn == int(elevio.BT_HallUp))  || (f == 0 && btn == elevio.BT_HallDown){
				fmt.Print("|     ");
			} else {
				if es.Requests[f][btn] == 1{
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


func Elevator_uninitialized(es *Elevator){  //initialize elevator, passing pointer
	es.Floor = -1
	es.Dirn = elevio.MD_Stop
	es.Behaviour = EB_Idle
	es.DoorOpenDuration = 3.0
}