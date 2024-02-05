package fsm

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/requests"
	"fmt"
)


func fsm_server(buttons chan elevio.ButtonEvent, floors chan int, obstr chan bool, stop chan bool, elev elevator.Elevator){

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
			print("TIMER INSERT HERE")
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
		
	

	}




}



func fsm_onInitBetweenFloors(e Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	e.dirn = elevio.MD_Down
	e.behaviour = elevator.EB_Moving
}

