package main

import (
	"Elev-project/Network-go-master/network/bcast"
	"Elev-project/Network-go-master/network/localip"
	"Elev-project/Network-go-master/network/peers"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/driver-go-master/fsm"
	"Elev-project/driver-go-master/cost_function"
	"flag"
	"fmt"
	"os"
	"os/exec"
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
	elevStateTx := make(chan elevator.Elevator)
	elevStateRx := make(chan elevator.Elevator)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(20008, elevStateTx)
	go bcast.Receiver(20008, elevStateRx)

	var elev elevator.Elevator

	//Processing pairs
	print("This is slave\n")
	timer1 := time.NewTimer(2 * time.Second)
	
	backupLoop:
		for {
			select {
			case elev = <-elevStateRx:
				fmt.Print("\n\nElev msg recieved:\n")
				elevator.Elevator_print(elev)
				fmt.Print("\n\n")
				timer1.Reset(2 * time.Second)
			case <-timer1.C:
				break backupLoop
			}
		}
	fmt.Print("Spawning backup\n")
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	print("This is now master\n")



	numFloors := 4
	elevio.Init("localhost:15657", numFloors)
	fsm.Elev_init(&elev, id)
	

	


	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// The example message. We just send one of these every second.

	go func() {

		for {
			//helloMsg.Iter++
			elevStateTx <- elev
			time.Sleep(500 * time.Millisecond)
		}
	}()

	

	for {
		//fsm.Fsm_server(drv_buttons, drv_floors, drv_obstr, drv_stop, &elev)
		
		fmt.Print("\n\nElev print main:\n")
		elevator.Elevator_print(elev)
		fmt.Print("\n\n")
		

		
		i := cost_function.TimeToIdle(elev)
		fmt.Printf("\nTime to idle: %d\n", i)

		select {

		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			fsm.Fsm_onRequestButtonPress(a, &elev)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			fsm.Fsm_onFloorArrival(a, &elev)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			//lag ny funksjon her eller finnes det allerede?

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			//lag ny funksjon her eller finnes det allerede? tror det sto noe om at det
			//ikke var definert noen oppførsel. kan velge selv?

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-elevStateRx:
		
			fmt.Print("\n\nElev msg recieved:\n")
			elevator.Elevator_print(a)
			fmt.Print("\n\n")
			
		}

	}

}
