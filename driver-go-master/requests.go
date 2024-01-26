package main

func requestsAbove(e Elevator) int {
	for f := e.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func requestsBelow(e Elevator) int {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func requestsHere(e Elevator) int {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return 1
		}
	}
	return 1
}


