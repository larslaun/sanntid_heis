package settings


import (
	"time"
)


const NumElevs = 3
const NumFloors = 4

const MaxTransmissionFailures = 10

const WatchdogTimeoutDuration = time.Duration(10) * time.Second
const DoorOpenDuration = time.Duration(3) * time.Second

