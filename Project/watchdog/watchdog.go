package watchdog

import (
	"Elev-project/networkDriver/network/peers"
	"Elev-project/communicationHandler/distributor"
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/requests"
	"Elev-project/settings"

	"fmt"
	"strconv"
	"time"
)

func LocalWatchdog(floors chan int, elev *elevator.Elevator, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder,elevStateRx chan elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator) {
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration)
	for {
		select {
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev) {
				fmt.Print("\nWatchdog fired\n")
				elev.Available = false
				distributor.RedistributeFaultyElevOrders(elevOrderTx, elevOrderRx, elevStateRx, elevators, elev)
			} else {
				watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			}
		case <-floors:
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			fmt.Print("New floor reached")
			//elev.Available = true
		}
	}
}


func NetworkWatchdog(peerUpdateCh chan peers.PeerUpdate, localElev *elevator.Elevator, elevators *[settings.N_ELEVS]elevator.Elevator, recoveryElevators *[settings.N_ELEVS]elevator.Elevator, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder,elevStateRx chan elevator.Elevator) {
	for {
		select {
		case peers := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peers.Peers)
			fmt.Printf("  New:      %q\n", peers.New)
			fmt.Printf("  Lost:     %q\n", peers.Lost)


			if peers.New != "" {
				newElev, _ := strconv.Atoi(peers.New)
				localElev.Available = true
				elevators[newElev].Available = true
				recoveryElevators[newElev].Available = true
				distributor.RecoverCabOrders(elevOrderTx, elevOrderRx,elevStateRx, elevators, &recoveryElevators[newElev])
			}


			lostElevs := peers.Lost
			for _, s := range lostElevs {
				s, _ := strconv.Atoi(s)
				fmt.Printf("s val: %d\n", s)
				elevators[s].Available = false
				recoveryElevators[s] = elevators[s]

				
				if len(peers.Peers) !=  0{
					if localElev.ID == peers.Peers[0] {
					distributor.RedistributeFaultyElevOrders(elevOrderTx, elevOrderRx,elevStateRx, elevators, &elevators[s])
					}
				}else{
					localElev.Available = false
				}
			}
		}
	}
}