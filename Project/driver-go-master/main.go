package main

import (
	
	"Driver-go/elevio"
	"Driver-go/fsm"
	"Driver-go/networkfunc"
	"Driver-go/elevator"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	conn :=	networkfunc.InitConn("localhost", "20008")
	defer conn.Close()

	var testelev elevator.Elevator
	elevator.Elevator_uninitialized(&testelev)
	
	go networkfunc.WriteStateToUDP(conn, testelev)
	go networkfunc.ReadFromUDP(conn)


	select{
	}
	


	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	


	
	fsm.Fsm_server(drv_buttons, drv_floors, drv_obstr, drv_stop)
	
}
