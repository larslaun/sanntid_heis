package fsm

import (
	"Elev-project/communicationHandler/distributor"
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/elevatorDriver/requests"
	"Elev-project/settings"
	"fmt"
	"time"
	"strconv"
)

func initBetweenFloors(elev *elevator.Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Dirn = elevio.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func FsmServer(elevStateRx chan elevator.Elevator, elevOrderRx chan elevator.ElevatorOrder, elevOrderTx chan elevator.ElevatorOrder, buttons chan elevio.ButtonEvent, floors chan int, obstruction chan bool, stop chan bool, elev *elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator) {
	go updateLights(elevators, elev)
	localID, _ := strconv.Atoi(elev.ID)

	for {
		select {
		case receivedOrder := <-elevOrderRx:
			if receivedOrder.RecipientID == elev.ID {
				fmt.Print("Received new order: ")
				fmt.Printf("%+v\n", receivedOrder.Order)
				onRequestButtonPress(receivedOrder.Order, elev)
			}

		case buttonPress := <-buttons:
			go distributor.DistributeOrder(buttonPress, elevOrderTx, elevOrderRx, elevStateRx, elevators, elev, localID)
			

		case currentFloor := <-floors:
			onFloorArrival(currentFloor, elev)
			

		case obstrState := <-obstruction:
			fmt.Printf("%+v\n", obstrState)
			/*
			elev.Obstruction = obstrState

			//While the obstruction  is true, onFloorArrival should continue to run, holding the door open.
			onFloorArrival(elev.Floor, elev)
			*/

			//fix later
		case stopState := <-stop:
			fmt.Printf("%+v\n", stopState)
			//elevio.SetStopLamp(stopState)
		}
	}

}

func onRequestButtonPress(buttons elevio.ButtonEvent, elev *elevator.Elevator) {


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
		var pair requests.DirnBehaviourPair = requests.ChooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)

			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })

			*elev = requests.ClearRequestAtCurrentFloor(*elev)

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)
		case elevator.EB_Idle:
		}
	}
	//SetCabLights(*elev)
}

func onFloorArrival(newFloor int, elev *elevator.Elevator) {


	elev.Floor = newFloor //dobbeltsjekk at det faktisk er den nye etasjen som blir tatt inn her
	elevio.SetFloorIndicator(newFloor)

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if requests.ShouldStop(*elev) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			*elev = requests.ClearRequestAtCurrentFloor(*elev)

			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })

			//SetCabLights(*elev)
			elev.Behaviour = elevator.EB_DoorOpen

		}
	case elevator.EB_DoorOpen:
		elevio.SetDoorOpenLamp(true)
		time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })
	}
}

func SetCabLights(elev elevator.Elevator) {
	for floor := 0; floor < settings.N_FLOORS; floor++ {
		if elev.Requests[floor][elevio.BT_Cab] {
			elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
		} else {
			elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
		}
	}
}

func SetHallLights(elevators *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator) {
	
	//making a matrix with zeros
	hallMatrix := make([][]bool, settings.N_FLOORS)
	for i := range hallMatrix {
		hallMatrix[i] = make([]bool, settings.N_BUTTONS-1) //only including hall-requests
	}

	//Iterating through each Hall-request in every elevator's matrix and OR'ing with every element in the hallMatrix.
	//This creates a "common" boolean matrix for hallCalls used to light every hall call button of the same type.
	for id := 0; id < len(elevators); id++ {
		if elevators[id].NetworkAvailable{
			for floor := 0; floor < settings.N_FLOORS; floor++ {
				for btn := elevio.BT_HallUp; btn <= elevio.BT_HallDown; btn++ {
					hallMatrix[floor][btn] = hallMatrix[floor][btn] || elevators[id].Requests[floor][btn]
				}
			}
		}
	}
	
	if localElev.NetworkAvailable == false{ //turn of hall calls from other elevators in case of network loss
		for floor := 0; floor < settings.N_FLOORS; floor++ {
			for btn := elevio.BT_HallUp; btn <= elevio.BT_HallDown; btn++ {
				hallMatrix[floor][btn] = localElev.Requests[floor][btn]
			}
		}
	}


	for floor := 0; floor < settings.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn <= elevio.BT_HallDown; btn++ {
			if hallMatrix[floor][btn] {
				elevio.SetButtonLamp(btn, floor, true)
			} else {
				elevio.SetButtonLamp(btn, floor, false)
			}
		}
	}
}


func updateLights(elevators *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator){
	for{
		SetHallLights(elevators, localElev)
		SetCabLights(*localElev)
		time.Sleep(20 * time.Millisecond)
	}
}


func onDoorTimeout(elev *elevator.Elevator) {

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		var pair requests.DirnBehaviourPair = requests.ChooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour

		switch elev.Behaviour {
		case elevator.EB_DoorOpen:
			time.AfterFunc(settings.DoorOpenDuration, func() { onDoorTimeout(elev) })
			*elev = requests.ClearRequestAtCurrentFloor(*elev)

		case elevator.EB_Moving:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elev.Dirn) //Var ikke i elev_algo men tror den må være her
		case elevator.EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elev.Dirn)
		}

	}

}

func ElevatorInit(elev *elevator.Elevator, elevID string) {
	elevator.InitializeElevStates(elev, elevID)
	initBetweenFloors(elev)
	
	SetCabLights(*elev)
	for floor := 0; floor < settings.N_FLOORS; floor++ {
		for btn := elevio.BT_HallUp; btn < elevio.BT_Cab+1; btn++ {
			elevio.SetButtonLamp(btn, floor, false)
		}
	}

	print("\nElevator initialized at following state: \n")
	elevator.PrintElevator(*elev)
}
