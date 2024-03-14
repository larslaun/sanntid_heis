package watchdog

import (
	"Elev-project/communicationHandler/distributor"
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/requests"
	"Elev-project/networkDriver/network/peers"
	"Elev-project/settings"

	"fmt"
	"strconv"
	"time"
)

func LocalWatchdog(floor chan int, elev *elevator.Elevator, distributeElevState chan elevator.Elevator, orderEvent chan elevator.ElevatorOrder,elevatorArray *[settings.N_ELEVS]elevator.Elevator) {
	watchdogTimer := time.NewTimer(settings.WATCHDOG_TIMEOUT_DURATION)
	idleFlag := true
	localID, _ := strconv.Atoi(elev.ID)
	timeoutCounter := 0

	for {
		select {
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev){
				if idleFlag{
					idleFlag = false
					watchdogTimer.Reset(settings.WATCHDOG_TIMEOUT_DURATION)
				}else{
					elev.Available = false
					elevatorArray[localID].Available = false
					distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, elev, localID, distributeElevState)
					watchdogTimer.Reset(settings.WATCHDOG_TIMEOUT_DURATION)
				}
			} else {
				timeoutCounter++
				if timeoutCounter == settings.MAX_WATCHDOG_TIMEOUT{
					timeoutCounter = 0
					elev.Available = true
				}
				watchdogTimer.Reset(settings.WATCHDOG_TIMEOUT_DURATION)
				idleFlag = true
			}
		case <-floor:
			elev.Available = true
			idleFlag = true
			if !watchdogTimer.Stop(){
				<-watchdogTimer.C
			}
			watchdogTimer.Reset(settings.WATCHDOG_TIMEOUT_DURATION)
			
		}
	}
}


func NetworkWatchdog(peerUpdateCh chan peers.PeerUpdate, localElev *elevator.Elevator, elevatorArray *[settings.N_ELEVS]elevator.Elevator, distributeElevState chan elevator.Elevator, orderEvent chan elevator.ElevatorOrder) {
	localID, _ := strconv.Atoi(localElev.ID)
	recoveryElevators := elevator.ElevatorArrayInit()

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
				distributor.RecoverCabOrders(orderEvent, distributeElevState,&recoveryElevators[newElev])
			}

			lostElevs := peers.Lost
			for _, elevID := range lostElevs {
				if elevID != localElev.ID{
					elevID, _ := strconv.Atoi(elevID)
					elevatorArray[elevID].NetworkAvailable = false
				}
			}

			if len(peers.Peers)==0{
				elevatorArray[localID].NetworkAvailable = false
				localElev.NetworkAvailable = false
				distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, &elevatorArray[localID], localID, distributeElevState)
			}

			for _, elevID := range lostElevs {
				if elevID != localElev.ID{
					elevID, _ := strconv.Atoi(elevID)
					elevatorArray[elevID].NetworkAvailable = false
					recoveryElevators[elevID] = elevatorArray[elevID]

					if len(peers.Peers) !=  0{
						if localElev.ID == peers.Peers[0] {
							distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, &elevatorArray[elevID], localID, distributeElevState) 
						}
					}else{
						distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, &elevatorArray[elevID], localID, distributeElevState)
					}
				}
			}
		}
	}
}