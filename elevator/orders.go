package elevator

// e.hallRequests is a 2D matrix that stores hall button requests at a given floor
// e.cabRequests is a 1D matrix that stores cab button requests at a given floor
// buttons are: BTN_HALLUP, BTN_HALLDOWN, BTN_HALLCAB

func (e *Elevator) RequestsAbove() bool {
	for f := e.floor + 1; f < NUM_FLOORS; f++ {
		for btn := 0; btn < 2; btn++ { // Only hall buttons
			if e.hallRequests[f][btn] {
				return true
			}
		}
		if e.cabRequests[f] {
			return true
		}
	}
	return false
}

func (e *Elevator) RequestsBelow() bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < 2; btn++ { // Only hall buttons
			if e.hallRequests[f][btn] {
				return true
			}
		}
		if e.cabRequests[f] {
			return true
		}
	}
	return false
}

func (e *Elevator) RequestsHere() bool {
	for btn := 0; btn < 2; btn++ { // Only hall buttons
		if e.hallRequests[e.floor][btn] {
			return true
		}
	}
	if e.cabRequests[e.floor] {
		return true
	}
	return false
}

func (e *Elevator) RequestsChooseDirection() DirBehaviourPair {
	switch e.direction {
	case DIR_UP:
		if e.RequestsAbove() {
			return DirBehaviourPair{DIR_UP, BEHAVIOUR_MOVING}
		} else if e.RequestsHere() {
			return DirBehaviourPair{DIR_STOP, BEHAVIOUR_DOOR_OPEN}
		} else if e.RequestsBelow() {
			return DirBehaviourPair{DIR_DOWN, BEHAVIOUR_MOVING}
		} else {
			return DirBehaviourPair{DIR_STOP, BEHAVIOUR_IDLE}
		}
	case DIR_DOWN:
		if e.RequestsBelow() {
			return DirBehaviourPair{DIR_DOWN, BEHAVIOUR_MOVING}
		} else if e.RequestsHere() {
			return DirBehaviourPair{DIR_STOP, BEHAVIOUR_DOOR_OPEN}
		} else if e.RequestsAbove() {
			return DirBehaviourPair{DIR_UP, BEHAVIOUR_MOVING}
		} else {
			return DirBehaviourPair{DIR_STOP, BEHAVIOUR_IDLE}
		}
	case DIR_STOP:
		if e.RequestsHere() {
			return DirBehaviourPair{DIR_STOP, BEHAVIOUR_DOOR_OPEN}
		} else if e.RequestsAbove() {
			return DirBehaviourPair{DIR_UP, BEHAVIOUR_MOVING}
		} else if e.RequestsBelow() {
			return DirBehaviourPair{DIR_DOWN, BEHAVIOUR_MOVING}
		} else {
			return DirBehaviourPair{DIR_STOP, BEHAVIOUR_IDLE}
		}
	default:
		return DirBehaviourPair{DIR_STOP, BEHAVIOUR_IDLE}
	}
}

func (e *Elevator) RequestsShouldStop() bool {
	switch e.direction {
	case DIR_DOWN:
		return e.hallRequests[e.floor][BTN_HALLDOWN] || e.cabRequests[e.floor] || !e.RequestsBelow()
	case DIR_UP:
		return e.hallRequests[e.floor][BTN_HALLUP] || e.cabRequests[e.floor] || !e.RequestsAbove()
	default:
		return true
	}
}

func (e *Elevator) RequestsShouldClearImmediately(buttonFloor int, buttonType Button) bool {
	switch e.config.clearRequestVariation {
	case CV_All:
		return e.floor == buttonFloor
	case CV_InDirn:
		return e.floor == buttonFloor && ((e.direction == DIR_UP && buttonType == BTN_HALLUP) ||
			(e.direction == DIR_DOWN && buttonType == BTN_HALLDOWN) ||
			e.direction == DIR_STOP ||
			buttonType == BTN_HALLCAB)
	default:
		return false
	}
}

func RequestsClearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariation {
	case CV_All:
		for btn := 0; btn < 2; btn++ {
			e.hallRequests[e.floor][btn] = false
		}
		e.cabRequests[e.floor] = false

	case CV_InDirn:
		e.cabRequests[e.floor] = false
		switch e.direction {
		case DIR_UP:
			if !e.RequestsAbove() && !e.hallRequests[e.floor][BTN_HALLUP] {
				e.hallRequests[e.floor][BTN_HALLDOWN] = false
			}
			e.hallRequests[e.floor][BTN_HALLUP] = false
		case DIR_DOWN:
			if !e.RequestsBelow() && !e.hallRequests[e.floor][BTN_HALLDOWN] {
				e.hallRequests[e.floor][BTN_HALLUP] = false
			}
			e.hallRequests[e.floor][BTN_HALLDOWN] = false
		default:
			e.hallRequests[e.floor][BTN_HALLUP] = false
			e.hallRequests[e.floor][BTN_HALLDOWN] = false
		}
	}

	return e
}
