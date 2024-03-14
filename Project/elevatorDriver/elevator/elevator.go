package elevator

import (
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/settings"
	"fmt"
	"strconv"
)

//bytte til elevatorArray

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
	NetworkAvailable bool
	Obstruction bool
}

type ElevatorOrder struct{
	RecipientID string
	Order elevio.ButtonEvent
}



func dirnToString(motorDir elevio.MotorDirection) string{
	if motorDir == elevio.MD_Up{
		return "MD_Up"
	} else if motorDir==elevio.MD_Down {
		return "MD_Down"
	}else{
		return "MD_Stop"
	}
}

func behaviourToString(eb ElevatorBehaviour) string{
	if eb == EB_Idle{
		return "EB_Idle"  
	} else if eb==EB_DoorOpen {
		return "EB_DoorOpen"
	}else{
		return "EB_Moving"
	} 
}


func PrintElevator(es Elevator){
	fmt.Print("\nElevator ID: ")
	fmt.Print(es.ID)
	fmt.Printf("\nAvailable: ")
	fmt.Print(es.Available)
	fmt.Printf("\nNetwork Available: ")
	fmt.Print(es.NetworkAvailable)
	fmt.Printf("\nObstruction: ")
	fmt.Print(es.Obstruction)
	fmt.Print("\n")
	fmt.Print("  +--------------------+\n")
	fmt.Printf("  |floor = %-2d|\n", es.Floor)
    fmt.Printf("  |dirn  = %-12.12s|\n", dirnToString(es.Dirn))
    fmt.Printf("  |behav = %-12.12s|\n", behaviourToString(es.Behaviour)) 
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
			}
		}
		fmt.Print("|\n");
	}
	fmt.Print("  +--------------------+\n");
	
}


func InitializeElevStates(elev *Elevator, elevID string){  
	elev.Floor = 1  
	elev.Dirn = elevio.MD_Stop
	elev.Behaviour = EB_Idle
	elev.ID = elevID
	elev.Available = false
	elev.NetworkAvailable = false
	elev.Obstruction = false 
}


func ElevatorArrayInit() [settings.N_ELEVS]Elevator{
	var elevatorArray = [settings.N_ELEVS]Elevator{}

	for i := 0; i < settings.N_ELEVS; i++ {
		InitializeElevStates(&elevatorArray[i], strconv.Itoa(i))
		PrintElevator(elevatorArray[i])
	}
	return elevatorArray
}