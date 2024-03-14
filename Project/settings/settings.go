package settings

import (
	"time"
)


//Elevator setting
const N_ELEVS int = 3
const N_FLOORS int = 4
const N_BUTTONS int = 3

const DOOR_OPEN_DURATION = time.Duration(3) * time.Second
const LIGHT_REFRESH__RATE = 200 * time.Millisecond

//Transmission settings
const STATE_TRANSMISSION_RATE = 20 * time.Millisecond
const ORDER_TRANSMISSION_RATE = time.Duration(5) * time.Millisecond
const MAX_TRANSMISSION_FAILURES = 100


//Watchdog settting
const WATCHDOG_TIMEOUT_DURATION = time.Duration(10) * time.Second
const MAX_WATCHDOG_TIMEOUT = 2


//Cost calculation settings
const DOOROPENTIME = 3
const TRAVELTIME = 3


