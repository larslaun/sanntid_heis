package hallAssigner

import(
	"Elev-project/driver-go-master/cost_function"
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"Elev-project/settings"
	"Elev-project/collector"
	"strconv"
	"fmt"
)



func ChooseOptimalElev(buttonPress elevio.ButtonEvent, elevators [settings.NumElevs]elevator.Elevator) collector.ElevatorOrder{
	
	var optimalElevID string
	var lowestCost = 1000000
	var currCost int	

	var order collector.ElevatorOrder

	for i := 0; i < settings.NumElevs; i++ {
		if elevators[i].Available {
			currCost = cost_function.TimeToIdle(elevators[i])
			if currCost < lowestCost{
				optimalElevID = strconv.Itoa(i)
				lowestCost = currCost
				order = collector.ElevatorOrder{RecipientID: optimalElevID, Order: buttonPress}
			}
		} 
	}
	fmt.Printf("\n COST CALCULATED: %d\n", currCost)
	return order
}