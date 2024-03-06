package watchdog


import (
	"Elev-project/driver-go-master/elevator"
	"Elev-project/driver-go-master/elevio"
	"fmt"

)

const Time int = 3

WatchdogTimer := time.NewTimer(Time * time.Second)

/*
Loop:
	for {
		select {
		case a <- 
			
		case 
			
		}
	}


func localWatchdog(floors chan int){
	watchdogTimer := time.NewTimer(settings.WatchdogTimeoutDuration * time.Second)
	for{
		select{
		case <-watchdogTimer.C:
			if not empty{
				redistribute FLAGG?
			}
		case <-floors:
			watchdogTimer.Reset(settings.WatchdogTimeoutDuration * time.Second)
		}
	}
}



