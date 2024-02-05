package fsm

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/requests"
	"fmt"
)

func Fsm_onInitBetweenFloors(elev *elevator.Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Dirn = elevio.MD_Down	
	elev.Behaviour = elevator.EB_Moving
}



func Fsm_server(buttons chan elevio.ButtonEvent, floors chan int, obstr chan bool, stop chan bool, elev elevator.Elevator){

	select {
	case a := <-buttons:
		fmt.Printf("%+v\n", a)
		fsm_onRequestButtonPress(a, elev)
		
		

	case a := <-floors:
		fmt.Printf("%+v\n", a)
		fsm_onFloorArrival(a, elev)
		

	case a := <-obstr:
		fmt.Printf("%+v\n", a)
		//lag ny funksjon her eller finnes det allerede?
		

	case a := <-stop:
		fmt.Printf("%+v\n", a)
		//lag ny funksjon her eller finnes det allerede? tror det sto noe om at det
		//ikke var definert noen oppførsel. kan velge selv?
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
			elev.Requests[buttons.Floor][buttons.Button] = true
		}
		
	case elevator.EB_Moving:
		elev.Requests[buttons.Floor][buttons.Button] = true
		
	case elevator.EB_Idle:
		elev.Requests[buttons.Floor][buttons.Button] = true
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



func fsm_onFloorArrival(newFloor int, elev elevator.Elevator){

	elevator.Elevator_print(elev)
	
	elev.Floor = newFloor   //dobbeltsjekk at det faktisk er den nye etasjen som blir tatt inn her
	elevio.SetFloorIndicator(newFloor)

	switch elev.Behaviour{
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(elev){
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elev = requests.RequestsClearAtCurrentFloor(elev)
			print("TIMER INSERT HERE")  //Sett inn timer
			//timer_start(elev.doorOpenDuration)
			setAllLights(elev)
			elev.Behaviour = elevator.EB_DoorOpen
		}
	}
	print("\nNew state:\n")
	elevator.Elevator_print(elev)
	
}







func setAllLights(elev elevator.Elevator){
	for floor := 0; floor < elevator.N_FLOORS; floor++{
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++{  //Tror dette er fikset
			if elev.Requests[floor][btn]{
				elevio.SetButtonLamp(btn, floor, true) 
			}
		}
	}
}

