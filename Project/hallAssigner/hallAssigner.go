package hallAssigner

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/hallAssigner/cost"
	"Elev-project/settings"
	"strconv"
)

//Estimates which elevator sholud serve an incomming request and returns it as an ElevatorOrder
func ChooseOptimalElev(buttonPress elevio.ButtonEvent, elevators [settings.NumElevs]elevator.Elevator, localID int) elevator.ElevatorOrder {

	var optimalElevID string
	var lowestCost = 1000000
	var currCost int
	var order elevator.ElevatorOrder

	for i := 0; i < settings.NumElevs; i++ {
		if elevators[i].Available {
			elevators[i].Requests[buttonPress.Floor][buttonPress.Button] = true
			currCost = cost.TimeToIdle(elevators[i])
			if currCost < lowestCost {
				optimalElevID = strconv.Itoa(i)
				lowestCost = currCost
				order = elevator.ElevatorOrder{RecipientID: optimalElevID, Order: buttonPress}
			}
		}
	}
	if !elevators[localID].Available {
		order = elevator.ElevatorOrder{RecipientID: elevators[localID].ID, Order: buttonPress}
	}
	return order
}
