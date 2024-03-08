package watchdog

import (
	"Elev-project/Network-go-master/network/peers"
	"Elev-project/distributor"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/requests"
	"Elev-project/settings"
	"Elev-project/collector"
	"fmt"

	//"fmt"
	"strconv"
	"time"
)

func LocalWatchdog(floors chan int, elev *elevator.Elevator, elevOrderTx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator) {
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration)
	for {
		select {
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev) {
				distributor.RedistributeFaultyElevOrders(elevOrderTx, elevStateRx, elevators, elev)
				elev.Available = false
			} else {
				watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			}
		case <-floors:
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			elev.Available = true
		}
	}
}

func NetworkWatchdog(peerUpdateCh chan peers.PeerUpdate, elevators *[settings.NumElevs]elevator.Elevator, recoveryElevators *[settings.NumElevs]elevator.Elevator) {
	for {
		select {
		case peers := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peers.Peers)
			fmt.Printf("  New:      %q\n", peers.New)
			fmt.Printf("  Lost:     %q\n", peers.Lost)

			newElev, _ := strconv.Atoi(peers.New)
			elevators[newElev].Available = true

			
			
			//fmt.Print("\nNew elevator:\n")
			//elevator.Elevator_print(elevators[newElev])

			lostElevs := peers.Lost
			for _, s := range lostElevs{
				s, _ := strconv.Atoi(s)
				elevators[s].Available = false
				recoveryElevators[s].Requests = elevators[s].Requests
				fmt.Printf("\nLost elevator ID %d:\n", s)
				elevator.Elevator_print(elevators[s])
			} 
		}
	}
}
