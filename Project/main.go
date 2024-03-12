package main

import (
	"Elev-project/Network-go-master/network/bcast"
	"Elev-project/Network-go-master/network/peers"
	"Elev-project/communicationHandler/collector"
	"Elev-project/communicationHandler/distributor"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/settings"
	"Elev-project/watchdog"

	//"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/fsm"

	//"Elev-project/driver-go-master/cost_function"

	//"fmt"
	"os"
	//"os/exec"
	//"time"
)

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	var elevPort string
	var id string

	args := os.Args

	id = args[1]
	elevPort = args[2]

	if elevPort == "" {
		elevPort = "15657"
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	elevStateTx := make(chan elevator.Elevator)
	elevStateRx := make(chan elevator.Elevator)
	go bcast.Transmitter(20010, elevStateTx)
	go bcast.Receiver(20010, elevStateRx)

	elevStateRx2 := make(chan elevator.Elevator)
	go bcast.Receiver(20010, elevStateRx2)

	//Må finne ut at av hvilke porter som kan brukes
	elevOrderTx := make(chan elevator.ElevatorOrder)
	elevOrderRx := make(chan elevator.ElevatorOrder)
	go bcast.Transmitter(21010, elevOrderTx)
	go bcast.Receiver(21010, elevOrderRx)

	var elev elevator.Elevator
	//This is where process pairs were
	elevio.Init("localhost:"+elevPort, settings.N_FLOORS)
	fsm.Elev_init(&elev, id)
	elevators := collector.ElevatorsInit()
	recoveryElevators := collector.ElevatorsInit()
	for i := 0; i < settings.NumElevs; i++ {
		elevator.Elevator_print(elevators[i])
	}

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstruction := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstruction)
	go elevio.PollStopButton(drv_stop)

	watchdog_floors := make(chan int)
	watchdog_elevOrderTx := make(chan elevator.ElevatorOrder)
	go bcast.Transmitter(21010, watchdog_elevOrderTx)

	watchdog_elevStateRx := make(chan elevator.Elevator)
	go bcast.Receiver(20010, watchdog_elevStateRx)

	go elevio.PollFloorSensor(watchdog_floors)
	go watchdog.LocalWatchdog(watchdog_floors, &elev, watchdog_elevOrderTx, elevOrderRx, watchdog_elevStateRx, &elevators)
	go watchdog.NetworkWatchdog(peerUpdateCh, &elev, &elevators, &recoveryElevators, watchdog_elevOrderTx, elevOrderRx, watchdog_elevStateRx)

	go collector.CollectStates(elevStateRx, &elevators, &elev)
	go distributor.DistributeState(elevStateTx, &elev)

	go fsm.Fsm_server(elevStateRx2, elevOrderRx, elevOrderTx, drv_buttons, drv_floors, drv_obstruction, drv_stop, &elev, &elevators)

	//i := cost_function.TimeToIdle(elev)
	//fmt.Printf("\nTime to idle: %d\n", i)

	select {}
}
