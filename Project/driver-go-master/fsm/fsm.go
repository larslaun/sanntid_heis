package fsm

import (
	"Elev-project/collector"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/requests"
	"Elev-project/distributor"
	"Elev-project/settings"
	"fmt"
	"time"
)

const TimerDuration = time.Duration(3) * time.Second

func Fsm_onInitBetweenFloors(elev *elevator.Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Dirn = elevio.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func Fsm_server(elevOrderRx chan collector.ElevatorOrder, elevOrderTx chan collector.ElevatorOrder, buttons chan elevio.ButtonEvent ,floors chan int, obstr chan bool, stop chan bool, elev *elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator) {
	
	for{

		//elevator.Elevator_print(*elev)
		select {

		case a := <- elevOrderRx: 
			if a.RecipientID==elev.ID{
				fmt.Print("Recieved new order: ")
				fmt.Printf("%+v\n", a.Order)
				Fsm_onRequestButtonPress(a.Order, elev)
			}


		case buttonPress := <-buttons:
			if buttonPress.Button == elevio.BT_Cab{
				fmt.Print("Recieved new cab order: ")
				fmt.Printf("%+v\n", buttonPress)
				Fsm_onRequestButtonPress(buttonPress, elev)
			}else{
				distributor.DistributeOrder(buttonPress, elevOrderTx, elevators)
			}


		case a := <-floors:
			fmt.Printf("%+v\n", a)
			Fsm_onFloorArrival(a, elev)

		case a := <-obstr:
			fmt.Printf("%+v\n", a)
			//lag ny funksjon her eller finnes det allerede?

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

			time.AfterFunc(TimerDuration, func() { onDoorTimeout(elev) })

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

			time.AfterFunc(TimerDuration, func() { onDoorTimeout(elev) })

			*elev = requests.RequestsClearAtCurrentFloor(*elev)

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)
		case elevator.EB_Idle:
		}
	}
	setAllLights(*elev)
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

			time.AfterFunc(TimerDuration, func() { onDoorTimeout(elev) })

			setAllLights(*elev)
			elev.Behaviour = elevator.EB_DoorOpen

		}
	}
	//print("\nNew state:\n")
	//elevator.Elevator_print(*elev)

}

func setAllLights(elev elevator.Elevator) {
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++ { //Tror dette er fikset, ønsker vi +1?
			if elev.Requests[floor][btn] {
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
			time.AfterFunc(TimerDuration, func() { onDoorTimeout(elev) })
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
	setAllLights(*elev)
	for floor := 0; floor < elevator.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++ { 
			elevio.SetButtonLamp(btn, floor, false)
		}
	}

	print("\nElevator initialized at following state: \n")
	elevator.Elevator_print(*elev)	
}