package elevator

import (
	"Elev-project/driver-go-master/elevio"
	"Elev-project/settings"
	"fmt"
	"strconv"
)

//Deklarerer her foreløpig


type ElevatorBehaviour int

const (
	EB_Idle   ElevatorBehaviour = 0
	EB_DoorOpen                = 1
	EB_Moving                = 2
)


type Elevator struct{
	Floor int
	Dirn elevio.MotorDirection
	Requests[settings.N_FLOORS][settings.N_BUTTONS] bool 
	Behaviour ElevatorBehaviour
	ID string
	Available bool 	
	Obstruction bool
}

type ElevatorOrder struct{
	RecipientID string
	Order elevio.ButtonEvent
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
	fmt.Print("\nElevator ID: ")
	fmt.Print(es.ID)
	fmt.Print("\n")
	fmt.Printf("Availible: ")
	fmt.Print(es.Available)
	fmt.Print("\n")
	fmt.Printf("Obstruction: ")
	fmt.Print(es.Obstruction)
	fmt.Print("\n")
	fmt.Print("  +--------------------+\n")
	fmt.Printf("  |floor = %-2d|\n", es.Floor)
    fmt.Printf("  |dirn  = %-12.12s|\n", Elevio_dirn_toString(es.Dirn))
    fmt.Printf("  |behav = %-12.12s|\n", Eb_toString(es.Behaviour)) 
	fmt.Print("  +--------------------+\n")
	fmt.Print("  |  | up  | dn  | cab |\n")
	for f := settings.N_FLOORS-1; f >= 0; f--{
		fmt.Printf("  | %d", f);
		for btn := 0; btn < settings.N_BUTTONS; btn++{
			if (f == settings.N_FLOORS-1 && btn == int(elevio.BT_HallUp))  || (f == 0 && btn == elevio.BT_HallDown){
				fmt.Print("|     ");
			} else {
				if es.Requests[f][btn]{
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


func Elevator_uninitialized(es *Elevator, elevID string){  //initialize elevator, passing pointer
	es.Floor = 1  //OBS! endret denne til 1 istedenfor -1 pga. index feil. Kan virke som dette funker siden man poller floor uansett og det vil bli endret til riktig
	es.Dirn = elevio.MD_Stop
	es.Behaviour = EB_Idle
	es.ID = elevID
	es.Available = false
	es.Obstruction = false 
}


func ElevatorsInit() [settings.NumElevs]Elevator{
	var elevators = [settings.NumElevs]Elevator{}

	for i := 0; i < settings.NumElevs; i++ {
		Elevator_uninitialized(&elevators[i], strconv.Itoa(i))
		Elevator_print(elevators[i])
	}
	return elevators
}