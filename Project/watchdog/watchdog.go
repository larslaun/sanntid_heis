package watchdog


import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/requests"
	"Elev-project/settings"
	"fmt"
	"time"

)




func LocalWatchdog(floors chan int, elev *elevator.Elevator, redistributeSignal chan bool){
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration * time.Second)
	for{
		select{
		case <-watchdogTimer.C:
			if requests.HasRequests(*elev){
				redistributeSignal <- true
				elev.Available = false
			} else{
				watchdogTimer.Reset(settings.WatchdogTimeoutDuration * time.Second)
			}
		case <-floors:
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration * time.Second)
			fmt.Print("\nWatchdog reset\n")
			elev.Available = true
		}
	}
}



