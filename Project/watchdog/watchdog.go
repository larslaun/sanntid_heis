package watchdog

import (
	"Elev-project/Network-go-master/network/peers"
	"Elev-project/collector"
	"Elev-project/distributor"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/requests"
	"Elev-project/settings"

	"fmt"
	"strconv"
	"time"
)

func LocalWatchdog(floors chan int, elev *elevator.Elevator, elevOrderTx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator) {
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration)
	for {
		select {
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev) {
				elev.Available = false
				distributor.RedistributeFaultyElevOrders(elevOrderTx, elevStateRx, elevators, elev)
			} else {
				watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			}
		case <-floors:
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			elev.Available = true
		}
	}
}

func NetworkWatchdog(peerUpdateCh chan peers.PeerUpdate, localElev *elevator.Elevator, elevators *[settings.NumElevs]elevator.Elevator, recoveryElevators *[settings.NumElevs]elevator.Elevator, elevOrderTx chan collector.ElevatorOrder, elevStateRx chan elevator.Elevator) {
	for {
		select {
		case peers := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peers.Peers)
			fmt.Printf("  New:      %q\n", peers.New)
			fmt.Printf("  Lost:     %q\n", peers.Lost)

			localID, _ := strconv.Atoi(localElev.ID)


			if peers.New != "" {
				newElev, _ := strconv.Atoi(peers.New)
				elevators[newElev].Available = true
				recoveryElevators[newElev].Available = true
				//fmt.Print("Will recover this now:\n")
				//elevator.Elevator_print(recoveryElevators[newElev])
				distributor.RecoverCabOrders(elevOrderTx, elevStateRx, elevators, &recoveryElevators[newElev])
			}

			//fmt.Print("\nNew elevator:\n")
			//elevator.Elevator_print(elevators[newElev])

			lostElevs := peers.Lost
			for _, s := range lostElevs {
				s, _ := strconv.Atoi(s)
				fmt.Printf("s val: %d\n", s)
				elevators[s].Available = false
				recoveryElevators[s] = elevators[s]
				//fmt.Printf("\nLost elevator ID %d:\n", s)
				//fmt.Print("Recovery state saved:\n")
				//elevator.Elevator_print(recoveryElevators[s])

				
				if len(peers.Peers) ==  0{
					//Legg inn if not localElev her???
					elevators[localID].Available = true
					distributor.RedistributeFaultyElevOrders(elevOrderTx, elevStateRx, elevators, &elevators[s])
				
				}else if localElev.ID == peers.Peers[0] {
					//elevator.Elevator_print(elevators[s])
					distributor.RedistributeFaultyElevOrders(elevOrderTx, elevStateRx, elevators, &elevators[s])
				}
			}
		}
	}
}