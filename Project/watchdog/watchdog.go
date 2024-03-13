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

func LocalWatchdog(floor chan int, elev *elevator.Elevator, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder,elevStateRx chan elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator) {
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration)
	idleFlag := true
	localID, _ := strconv.Atoi(elev.ID)

	for {
		select {
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev){
				if idleFlag{
					idleFlag = false
					watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
				}else{
				elev.Available = false
				distributor.RedistributeFaultyElevOrders(elevOrderTx, elevOrderRx, elevStateRx, elevatorArray, elev, localID)
				}
			} else {
				watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
				idleFlag = true
			}
		case <-floor:
			elev.Available = true
			if !watchdogTimer.Stop(){
				<-watchdogTimer.C
			}
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			
		}
	}
}


func NetworkWatchdog(peerUpdateCh chan peers.PeerUpdate, localElev *elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, recoveryElevators *[settings.N_ELEVS]elevator.Elevator, elevOrderTx chan elevator.ElevatorOrder, elevOrderRx chan elevator.ElevatorOrder,elevStateRx chan elevator.Elevator) {
	localID, _ := strconv.Atoi(localElev.ID)
	for {
		select {
		case peers := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peers.Peers)
			fmt.Printf("  New:      %q\n", peers.New)
			fmt.Printf("  Lost:     %q\n", peers.Lost)


			if peers.New != "" {
				newElev, _ := strconv.Atoi(peers.New)
				localElev.NetworkAvailable = true
				elevatorArray[newElev].NetworkAvailable = true
				recoveryElevators[newElev].NetworkAvailable = true
				distributor.RecoverCabOrders(elevOrderTx, elevOrderRx,elevStateRx, elevatorArray, &recoveryElevators[newElev], localID)
			}


			lostElevs := peers.Lost
			for _, s := range lostElevs {
				s, _ := strconv.Atoi(s)
				fmt.Printf("s val: %d\n", s)
				elevatorArray[s].NetworkAvailable = false
				recoveryElevators[s] = elevatorArray[s]

				
				if len(peers.Peers) !=  0{
					if localElev.ID == peers.Peers[0] {
					distributor.RedistributeFaultyElevOrders(elevOrderTx, elevOrderRx,elevStateRx, elevatorArray, &elevatorArray[s], localID)
					}
				}else{
					localElev.NetworkAvailable = false
					distributor.RedistributeFaultyElevOrders(elevOrderTx, elevOrderRx,elevStateRx, elevatorArray, &elevatorArray[s], localID)
				}
			}
		}
	}
}