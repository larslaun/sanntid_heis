package fsm

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/requests"
	"fmt"
)


func Fsm_server(buttons chan elevio.ButtonEvent, floors chan int, obstr chan bool, stop chan bool, elev elevator.Elevator){

	select {
	case a := <-buttons:
		fmt.Printf("%+v\n", a)
		fsm_onRequestButtonPress(a, elev)
		
		

	case a := <-floors:
		fmt.Printf("%+v\n", a)
		

	case a := <-obstr:
		fmt.Printf("%+v\n", a)
		
		

	case a := <-stop:
		fmt.Printf("%+v\n", a)
		
	}
}



func fsm_onRequestButtonPress(buttons elevio.ButtonEvent, elev elevator.Elevator){
	
	elevator.Elevator_print(elev)
	
	

	
	switch elev.Behaviour{
	case elevator.EB_DoorOpen:
		if(requests.RequestsShouldClearImmediately(elev, buttons.Floor, buttons.Button)){
			print("TIMER INSERT HERE") //Sett inn timer
			//timer_start(elev.doorOpenDuration)
		} else {
			elev.Requests[buttons.Floor][buttons.Button] = 1
		}
		
	case elevator.EB_Moving:
		elev.Requests[buttons.Floor][buttons.Button] = 1
		
	case elevator.EB_Idle:
		elev.Requests[buttons.Floor][buttons.Button] = 1
		var pair requests.DirnBehaviourPair = requests.RequestsChooseDirection(elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		switch pair.Behaviour{
		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			print("TIMER INSERT HERE")  //Sett inn timer
			//timer_start(elev.doorOpenDuration)
			elev = requests.RequestsClearAtCurrentFloor(elev)
		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)
		case elevator.EB_Idle:	
		}
	}
	setAllLights(elev)
	print("\nNew state:\n")
	elevator.Elevator_print(elev)
}



func fsm_onInitBetweenFloors(e Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	e.dirn = elevio.MD_Down
	e.behaviour = elevator.EB_Moving
}




func setAllLights(elev elevator.Elevator){
	for floor := 0; floor < elevator.N_FLOORS; floor++{
		for elevio.ButtonType btn := 0; btn < elevator.N_BUTTONS; btn++{  //Hvordan iterere gjennom buttontype enum??
			if elev.Requests[floor][btn] == 1{
				elevio.SetButtonLamp(btn, floor, true) 
			}
		}
	}
}

