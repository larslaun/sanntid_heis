package hallAssigner

import (
	"Elev-project/collector"
	"Elev-project/driver-go-master/cost_function"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/settings"
	"strconv"
	//"fmt"
)



func ChooseOptimalElev(buttonPress elevio.ButtonEvent, elevators [settings.NumElevs]elevator.Elevator) collector.ElevatorOrder{
	
	var optimalElevID string
	var lowestCost = 1000000
	var currCost int	

	var order collector.ElevatorOrder

	for i := 0; i < settings.NumElevs; i++ {
		if elevators[i].Available {
			//fmt.Printf("Calculating cost for ID %d:", i)
			//elevator.Elevator_print(elevators[i])
			elevators[i].Requests[buttonPress.Floor][buttonPress.Button] = true

			currCost = cost_function.TimeToIdle(elevators[i])
			//fmt.Printf("Cost for elevator ID %d is the following: %d\n", i, currCost)
			if currCost < lowestCost{
				optimalElevID = strconv.Itoa(i)
				lowestCost = currCost
				order = collector.ElevatorOrder{RecipientID: optimalElevID, Order: buttonPress}
			}
		} 
	}

	return order
}