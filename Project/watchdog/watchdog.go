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
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration)
	idleFlag := true
	localID, _ := strconv.Atoi(elev.ID)
	timeoutCounter := 0

	for {
		select {
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev){
				fmt.Print("\n1\n")
				if idleFlag{
					fmt.Print("\n2\n")
					idleFlag = false
					watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
				}else{
					fmt.Print("\n3\n")
					elev.Available = false
					elevatorArray[localID].Available = false
					distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, elev, localID, distributeElevState)
					watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
				}
			} else {
				fmt.Print("\n4\n")
				timeoutCounter++
				if timeoutCounter == 2{
					timeoutCounter = 0
					elev.Available = true
				}
				watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
				idleFlag = true
			}
		case <-floor:
			fmt.Print("\n5\n")
			elev.Available = true
			if !watchdogTimer.Stop(){
				<-watchdogTimer.C
			}
			idleFlag = true
			fmt.Print("\n6\n")
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration)
			
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


			if len(peers.Peers)==0{
				elevatorArray[localID].NetworkAvailable = false
				localElev.NetworkAvailable = false
				distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, &elevatorArray[localID], localID, distributeElevState)
			}
			

			lostElevs := peers.Lost
			for _, elevID := range lostElevs {
				if elevID != localElev.ID{
					elevID, _ := strconv.Atoi(elevID)
					fmt.Printf("Lost elev ID: %d\n", elevID)
					elevatorArray[elevID].NetworkAvailable = false
					recoveryElevators[elevID] = elevatorArray[elevID]

					if len(peers.Peers) !=  0{
						if localElev.ID == peers.Peers[0] {
							distributor.RedistributeFaultyElevOrders(orderEvent, elevatorArray, &elevatorArray[elevID], localID, distributeElevState)
						}
					}
				}
			}
		}
	}
}