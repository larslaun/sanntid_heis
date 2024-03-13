package hallAssigner

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/hallAssigner/cost"
	"Elev-project/settings"
	"strconv"
)

//Estimates which elevator sholud serve an incomming request and returns it as an ElevatorOrder
func ChooseOptimalElev(buttonPress elevio.ButtonEvent, elevatorArray [settings.N_ELEVS]elevator.Elevator, localID int) elevator.ElevatorOrder {

	var optimalElevID string
	var lowestCost = 1000000
	var currCost int
	var order elevator.ElevatorOrder

	for elev := 0; elev < settings.N_ELEVS; elev++ {
		if elevatorArray[elev].Available &&  elevatorArray[elev].NetworkAvailable{
			elevatorArray[elev].Requests[buttonPress.Floor][buttonPress.Button] = true
			currCost = cost.TimeToIdle(elevatorArray[elev])
			if currCost < lowestCost {
				optimalElevID = strconv.Itoa(elev)
				lowestCost = currCost
				order = elevator.ElevatorOrder{RecipientID: optimalElevID, Order: buttonPress}
			}
		}
	}
	if !elevatorArray[localID].NetworkAvailable {
		order = elevator.ElevatorOrder{RecipientID: elevatorArray[localID].ID, Order: buttonPress}
	}
	return order
}
