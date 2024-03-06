package main

import (
	"Elev-project/collector"
	"Elev-project/distributor"
	"Elev-project/driver-go-master/elevator"

	"Elev-project/Network-go-master/network/bcast"
	"Elev-project/Network-go-master/network/localip"
	"Elev-project/Network-go-master/network/peers"
	//"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/fsm"
	//"Elev-project/driver-go-master/cost_function"
	"flag"
	"fmt"
	"os"
	//"os/exec"
	//"time"
)

func main() {

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var elevPort string
	flag.StringVar(&elevPort, "port", "", "port of elev")
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	elevStateTx := make(chan elevator.Elevator)
	elevStateRx := make(chan elevator.Elevator)
	go bcast.Transmitter(20008, elevStateTx)
	go bcast.Receiver(20008, elevStateRx)
	
	//Må finne ut at av hvilke porter som kan brukes
	elevOrderTx := make(chan collector.ElevatorOrder)
	elevOrderRx := make(chan collector.ElevatorOrder)
	go bcast.Transmitter(21008, elevOrderTx)
	go bcast.Receiver(21008, elevOrderRx)

	var elev elevator.Elevator
	//This is where process pairs were
	numFloors := 4
	elevio.Init("localhost:"+elevPort, numFloors)
	fsm.Elev_init(&elev, id)
	elevators := collector.ElevatorsInit()
	for i := 0; i <3; i++ {
		elevator.Elevator_print(elevators[i])
	}

	
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	
	go collector.CollectStates(elevStateRx, &elevators)
	go distributor.DistributeState(elevStateTx, &elev)
	go distributor.DistributeOrder(drv_buttons, elevOrderTx, &elevators)


	go fsm.Fsm_server(elevOrderRx, drv_floors, drv_obstr, drv_stop, &elev)



	
	
	//i := cost_function.TimeToIdle(elev)
	//fmt.Printf("\nTime to idle: %d\n", i)

	select {
	/*
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
	*/
	}
}
