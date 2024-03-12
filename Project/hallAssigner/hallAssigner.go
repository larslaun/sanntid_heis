package hallAssigner

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/hallAssigner/cost"
	"Elev-project/settings"
	"strconv"
)

func ChooseOptimalElev(buttonPress elevio.ButtonEvent, elevators [settings.N_ELEVS]elevator.Elevator, localID int) elevator.ElevatorOrder {

	var optimalElevID string
	var lowestCost = 1000000
	var currCost int

	var order elevator.ElevatorOrder

	for i := 0; i < settings.N_ELEVS; i++ {
		if elevators[i].Available {
			//fmt.Printf("Calculating cost for ID %d:", i)
			//elevator.PrintElevator(elevators[i])
			elevators[i].Requests[buttonPress.Floor][buttonPress.Button] = true

			currCost = cost.TimeToIdle(elevators[i])
			//fmt.Printf("Cost for elevator ID %d is the following: %d\n", i, currCost)
			if currCost < lowestCost {
				optimalElevID = strconv.Itoa(i)
				lowestCost = currCost
				order = elevator.ElevatorOrder{RecipientID: optimalElevID, Order: buttonPress}
			}
		}
	}

	if elevators[localID].Available == false {
		order = elevator.ElevatorOrder{RecipientID: elevators[localID].ID, Order: buttonPress}
	}

	//fmt.Printf("Optimal ID calculated: " + optimalElevID + "\n")

	return order
}
