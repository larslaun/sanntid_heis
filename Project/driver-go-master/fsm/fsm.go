package fsm

import (
	"Elev-project/collector"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/requests"
	"Elev-project/settings"
	"fmt"
	"time"
)


func Fsm_onInitBetweenFloors(elev *elevator.Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Dirn = elevio.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func Fsm_server(elevStateRx chan elevator.Elevator, elevOrderRx chan collector.ElevatorOrder, floors chan int, obstruction chan bool, stop chan bool, elev *elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator) {

	for {

		//elevator.Elevator_print(*elev)
		select {

		case receivedElev := <-elevStateRx:
			SetHallLights(elevators)
			if elev.ID == receivedElev.ID{
				SetCabLights(receivedElev)
			}

		case a := <-elevOrderRx:
			fmt.Print("Received new order: ")
			fmt.Printf("%+v\n", a.Order)
			Fsm_onRequestButtonPress(a.Order, elev)
		
		case a := <-floors:
			//fmt.Printf("%+v\n", a)
			Fsm_onFloorArrival(a, elev)

		case a := <-obstruction:
			fmt.Printf("%+v\n", a)
			elev.Obstruction = a

			//While the obstruction  is true, onFloorArrival should continue to run, holding the door open. 
			for elev.Obstruction {
				Fsm_onFloorArrival(elev.Floor, elev)
			}
			

		case a := <-stop:
			fmt.Printf("%+v\n", a)
			//lag ny funksjon her eller finnes det allerede? tror det sto noe om at det
			//ikke var definert noen oppførsel. kan velge selv?

		}
	}

}

func Fsm_onRequestButtonPress(buttons elevio.ButtonEvent, elev *elevator.Elevator) {

	//elevator.Elevator_print(*elev)

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		if requests.RequestsShouldClearImmediately(*elev, buttons.Floor, buttons.Button) {
			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })

		} else {
			elev.Requests[buttons.Floor][buttons.Button] = true
		}

	case elevator.EB_Moving:
		elev.Requests[buttons.Floor][buttons.Button] = true

	case elevator.EB_Idle:
		elev.Requests[buttons.Floor][buttons.Button] = true
		var pair requests.DirnBehaviourPair = requests.RequestsChooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)

			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })

			*elev = requests.RequestsClearAtCurrentFloor(*elev)

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)
		case elevator.EB_Idle:
		}
	}
	SetCabLights(*elev)
	//print("\nNew state:\n")
	//elevator.Elevator_print(*elev)
}



func Fsm_onFloorArrival(newFloor int, elev *elevator.Elevator) {

	//elevator.Elevator_print(*elev)

	elev.Floor = newFloor //dobbeltsjekk at det faktisk er den nye etasjen som blir tatt inn her
	elevio.SetFloorIndicator(newFloor)

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(*elev) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			*elev = requests.RequestsClearAtCurrentFloor(*elev)

			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })

			SetCabLights(*elev)
			elev.Behaviour = elevator.EB_DoorOpen

		}
	case elevator.EB_DoorOpen:
		elevio.SetDoorOpenLamp(true)
		time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })
	}
}



func SetCabLights(elev elevator.Elevator) {
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		if elev.Requests[floor][elevio.BT_Cab] {
			elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
		} else {
			elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
		}
	}
}

func SetHallLights(elevators *[settings.NumElevs]elevator.Elevator) {
	//making a matrix with zeros
	hallMatrix := make([][]bool, elevator.N_FLOORS)
	for i := range hallMatrix {
		hallMatrix[i] = make([]bool, elevator.N_BUTTONS-1) //only including hall-requests
	}

	//Iterating through each Hall-request in every elevator's matrix and OR'ing with every element in the hallMatrix.
	//This creates a "common" boolean matrix for hallCalls used to light every hall call button of the same type. 
	for id := 0; id < len(elevators); id++ {
		for floor := 0; floor < elevator.N_FLOORS; floor++ {
			for btn := elevio.BT_HallUp; btn <= elevio.BT_HallDown; btn++ {
				hallMatrix[floor][btn] = hallMatrix[floor][btn] || elevators[id].Requests[floor][btn]
			}
		}
	}

	//Setting the lights using the bools in hallMatrix.
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn <= elevio.BT_HallDown; btn++ {
			if hallMatrix[floor][btn] {
				elevio.SetButtonLamp(btn, floor, true)
			} else {
				elevio.SetButtonLamp(btn, floor, false)
			}
		}
	}
}

func onDoorTimeout(elev *elevator.Elevator) {

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		var pair requests.DirnBehaviourPair = requests.RequestsChooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour

		switch elev.Behaviour {
		case elevator.EB_DoorOpen:
			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })
			*elev = requests.RequestsClearAtCurrentFloor(*elev)

		case elevator.EB_Moving:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elev.Dirn) //Var ikke i elev_algo men tror den må være her
		case elevator.EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elev.Dirn)
		}

	}
	//print("\nNew state:\n")
	//elevator.Elevator_print(*elev)

}

func Elev_init(elev *elevator.Elevator, elevID string) {
	elevator.Elevator_uninitialized(elev, elevID)
	//if elevio.GetFloor() == -1 {
	Fsm_onInitBetweenFloors(elev)
	//}
	SetCabLights(*elev)
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++ {
			elevio.SetButtonLamp(btn, floor, false)
		}
	}

	print("\nElevator initialized at following state: \n")
	elevator.Elevator_print(*elev)
}
