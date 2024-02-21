package main

import (
	"Elev-project/Network-go-master/network/bcast"
	"Elev-project/Network-go-master/network/localip"
	"Elev-project/Network-go-master/network/peers"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/fsm"
	"flag"
	"fmt"
	"os"
	"time"
)


func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	

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

	

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	


	// We make channels for sending and receiving our custom data types
	helloTx := make(chan elevator.Elevator)
	helloRx := make(chan elevator.Elevator)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	





	numFloors := 4

	
	elevio.Init("localhost:15657", numFloors)

	
	
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	
	
	var elev elevator.Elevator = fsm.Elev_init()



	// The example message. We just send one of these every second.
	
	go func() {
		var testelev elevator.Elevator
		elevator.Elevator_uninitialized(&testelev)
		
		helloMsg := testelev

		//helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			//helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()
	





	for{
		fsm.Fsm_server(drv_buttons, drv_floors, drv_obstr, drv_stop, elev)
		
		
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			//fmt.Printf("Received: %#v\n", a)
			elevator.Elevator_print(a)
		}
		
	}



}
