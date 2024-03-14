package fsm

import (
	"Elev-project/communicationHandler/distributor"
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/elevatorDriver/requests"
	"Elev-project/settings"
	"fmt"
	"strconv"
	"time"
)


//bytte elev til localElev?
//skrive button og ikke buttons i alle funksjoner
//floors til floor
//endre pair til newBehaviourPair

func initBetweenFloors(elev *elevator.Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Dirn = elevio.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func FsmServer(elevStateRx chan elevator.Elevator, elevOrderRx chan elevator.ElevatorOrder, elevOrderTx chan elevator.ElevatorOrder, buttonEvent chan elevio.ButtonEvent, floor chan int, obstruction chan bool, stop chan bool, elev *elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator) {
	go updateLights(elevatorArray, elev)
	localID, _ := strconv.Atoi(elev.ID)
	doorTimeout := time.NewTimer(settings.DoorOpenDuration)
	resetTimer := make(chan bool, 4)

	for {
		select {
		case receivedOrder := <-elevOrderRx:
			if receivedOrder.RecipientID == elev.ID {
				fmt.Print("Received new order: ")
				fmt.Printf("%+v\n", receivedOrder.Order)
				onRequestButtonPress(receivedOrder.Order, elev, resetTimer)
			}

		case buttonPress := <-buttonEvent:
			go distributor.DistributeOrder(buttonPress, elevOrderTx, elevOrderRx, elevStateRx, elevatorArray, elev, localID)

		case currentFloor := <-floor:
			onFloorArrival(currentFloor, elev, resetTimer)


		case obstrState := <-obstruction:
			fmt.Printf("%+v\n", obstrState)
			elev.Obstruction = obstrState

		case stopState := <-stop:
			fmt.Printf("%+v\n", stopState)
			panic("Stop button pressed")

		case <-doorTimeout.C:
			onDoorTimeout(elev, resetTimer)

		case <-resetTimer:
			doorTimeout.Reset(settings.DoorOpenDuration)

		}

	}
}


func onRequestButtonPress(buttonEvent elevio.ButtonEvent, elev *elevator.Elevator, resetTimer chan bool) {

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:

		if requests.RequestsShouldClearImmediately(*elev, buttonEvent.Floor, buttonEvent.Button) {

			resetTimer <- true
			elev.Behaviour = elevator.EB_DoorOpen
      
		} else {
			elev.Requests[buttonEvent.Floor][buttonEvent.Button] = true
		}

	case elevator.EB_Moving:
		elev.Requests[buttonEvent.Floor][buttonEvent.Button] = true

	case elevator.EB_Idle:

		elev.Requests[buttonEvent.Floor][buttonEvent.Button] = true
		var newBehaviourPair requests.DirnBehaviourPair = requests.ChooseDirection(*elev)
		elev.Dirn = newBehaviourPair.Dirn
		elev.Behaviour = newBehaviourPair.Behaviour

		switch elev.Behaviour {
		case elevator.EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			resetTimer <- true

			*elev = requests.ClearRequestAtCurrentFloor(*elev)

		case elevator.EB_Moving:
			elevio.SetMotorDirection(elev.Dirn)

		case elevator.EB_Idle:
		}
	}
}



func onFloorArrival(newFloor int, elev *elevator.Elevator, resetTimer chan bool) {
	elev.Floor = newFloor 

	elevio.SetFloorIndicator(newFloor)

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if requests.ShouldStop(*elev) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)

			*elev = requests.ClearRequestAtCurrentFloor(*elev)
			resetTimer <- true

			elev.Behaviour = elevator.EB_DoorOpen //skal den inni if over??

		}
	case elevator.EB_DoorOpen:
		elevio.SetDoorOpenLamp(true)

		resetTimer <- true

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


func SetHallLights(elevatorArray *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator) {

	hallMatrix := make([][]bool, settings.N_FLOORS)
	for i := range hallMatrix {
		hallMatrix[i] = make([]bool, settings.N_BUTTONS-1)
	}


	for id := 0; id < len(elevatorArray); id++ {
		if elevatorArray[id].NetworkAvailable {
			for floor := 0; floor < settings.N_FLOORS; floor++ {
				for btn := elevio.BT_HallUp; btn <= elevio.BT_HallDown; btn++ {
					hallMatrix[floor][btn] = hallMatrix[floor][btn] || elevatorArray[id].Requests[floor][btn]
				}
			}
		}
	}

	if !localElev.NetworkAvailable {
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

func updateLights(elevatorArray *[settings.N_ELEVS]elevator.Elevator, localElev *elevator.Elevator) {
	for {
		SetHallLights(elevatorArray, localElev)
		SetCabLights(*localElev)
		time.Sleep(200 * time.Millisecond)
	}
}

func onDoorTimeout(elev *elevator.Elevator, resetTimer chan bool) {
	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		var newBehaviourPair requests.DirnBehaviourPair = requests.ChooseDirection(*elev)
		elev.Dirn = newBehaviourPair.Dirn
		elev.Behaviour = newBehaviourPair.Behaviour


		if elev.Obstruction {
			elev.Behaviour = elevator.EB_DoorOpen
		}

		switch elev.Behaviour {
		case elevator.EB_DoorOpen:
			*elev = requests.ClearRequestAtCurrentFloor(*elev)
			resetTimer <- true

		case elevator.EB_Moving:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elev.Dirn) 

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
