package settings


import (
	"time"
)


const N_ELEVS int = 3
const N_FLOORS int = 4
const N_BUTTONS int = 3



const MaxTransmissionFailures = 100

const WatchdogTimeoutDuration = time.Duration(10) * time.Second
const DoorOpenDuration = time.Duration(3) * time.Second



//Cost calculation settings
const DOOROPENTIME = 3
const TRAVELTIME = 3


