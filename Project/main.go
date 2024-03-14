package main

import (
	"Elev-project/communicationHandler/collector"
	"Elev-project/communicationHandler/distributor"
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/elevatorDriver/fsm"
	"Elev-project/networkDriver/network/bcast"
	"Elev-project/networkDriver/network/peers"
	"Elev-project/settings"
	"Elev-project/watchdog"
	"os"
	"strconv"
)

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	
	id := "0"
	commPort := 20008
	elevPort := "15657"

	args := os.Args

	if len(args)>1{
	id = args[1]
	}
	if len(args)>2{
	commPort, _ = strconv.Atoi(args[2])
	}
	if len(args)>3{
		elevPort = args[3]
	}
	id_int , _ := strconv.Atoi(id)



	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	elevStateTx := make(chan elevator.Elevator)
	elevStateRx := make(chan elevator.Elevator)
	go bcast.Transmitter(commPort, elevStateTx)
	go bcast.Receiver(commPort, elevStateRx)

	elevOrderTx := make(chan elevator.ElevatorOrder)
	elevOrderRx := make(chan elevator.ElevatorOrder)
	go bcast.Transmitter(commPort+1000, elevOrderTx)
	go bcast.Receiver(commPort+1000, elevOrderRx)



	var elev elevator.Elevator
	elevio.Init("localhost:"+elevPort, settings.N_FLOORS)
	fsm.ElevatorInit(&elev, id)
	elevatorArray := elevator.ElevatorArrayInit()



	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstruction := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstruction)
	go elevio.PollStopButton(drv_stop)


	orderEvent := make(chan elevator.ElevatorOrder, 20)
	distributeElevState := make(chan elevator.Elevator,1)
	go distributor.DistributeOrder(orderEvent, elevOrderTx, elevOrderRx, distributeElevState, id_int)
	go collector.CollectStates(elevStateRx, &elevatorArray, &elev, distributeElevState)
	go distributor.DistributeState(elevStateTx, &elev)


	watchdog_floors := make(chan int)
	go elevio.PollFloorSensor(watchdog_floors)

	go watchdog.LocalWatchdog(watchdog_floors, &elev, distributeElevState, orderEvent, &elevatorArray)
	go watchdog.NetworkWatchdog(peerUpdateCh, &elev, &elevatorArray, distributeElevState, orderEvent)


	go fsm.FsmServer(elevOrderRx, orderEvent, drv_buttons, drv_floors, drv_obstruction, drv_stop, &elev, &elevatorArray)


	select {}
}
