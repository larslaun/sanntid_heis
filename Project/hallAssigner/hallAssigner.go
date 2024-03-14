package hallAssigner

import (
	"Elev-project/elevatorDriver/elevator"
	"Elev-project/elevatorDriver/elevio"
	"Elev-project/hallAssigner/cost"
	"Elev-project/settings"
	"strconv"
)

//Estimates which elevator sholud serve an incomming request and returns it as an ElevatorOrder
func ChooseOptimalElev(orderRecipientPair elevator.ElevatorOrder, elevatorArray [settings.N_ELEVS]elevator.Elevator, localID int) elevator.ElevatorOrder {

	var optimalElevID string
	var lowestCost = 1000000
	var currCost int
	originalOrderID := orderRecipientPair.RecipientID

	for elev := 0; elev < settings.N_ELEVS; elev++ {
		if elevatorArray[elev].Available &&  elevatorArray[elev].NetworkAvailable{
			elevatorArray[elev].Requests[orderRecipientPair.Order.Floor][orderRecipientPair.Order.Button] = true
			currCost = cost.TimeToIdle(elevatorArray[elev])
			if currCost < lowestCost {
				optimalElevID = strconv.Itoa(elev)
				lowestCost = currCost
				orderRecipientPair = elevator.ElevatorOrder{RecipientID: optimalElevID, Order: orderRecipientPair.Order}
			}
		}
	}
	if !elevatorArray[localID].NetworkAvailable {
		orderRecipientPair = elevator.ElevatorOrder{RecipientID: elevatorArray[localID].ID, Order: orderRecipientPair.Order}
	}
	
	if orderRecipientPair.Order.Button == elevio.BT_Cab{
		orderRecipientPair = elevator.ElevatorOrder{RecipientID: originalOrderID, Order: orderRecipientPair.Order}
	}

	return orderRecipientPair
}
