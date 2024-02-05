package main

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/fsm"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)


	//initalisere heis
	var elev elevator.Elevator
	elevator.Elevator_uninitialized(&elev)
	//

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		fsm.Fsm_server(drv_buttons, drv_floors, drv_obstr, drv_stop, elev)
	}
}
