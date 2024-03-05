package hallAssigner

import(
	"Elev-project/driver-go-main/cost_function"
	"Elev-project/driver-go-main/elevator"
	"Elev-project/driver-go-main/elevio"
	"Elev-project/settings"
	"Elev-project/collector"
	"strconv"
)



func ChooseOptimalElev(buttonPress elevio.ButtonEvent, [setting.NumElevs]elevators elevator.Elevator) ElevatorOrder{
	var optimalElevID string
	var lowestCost = 1000000
	var currCost int	

	var order collector.ElevatorOrder

	for i := 0; i < settings.NumElevs; i++ {
		if elevators[i].Available {
			currCost = cost_function.TimeToIdle(elevators[i])
			if currCost < lowestCost{
				optimalElevID = stconv.Itoa(i)
				lowestCost = currCost
				order := collector.ElevatorOrder{RecipientID: optimalElevID, Order: order}
			}
		}
	}
	return order
}