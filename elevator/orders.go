package elevator

// e.requests is a 2D matrix that stores what type of button is pushed at a given floor
//buttons are: BTN_HALLUP, BTN_HALLDOWN, BTN_HALLCAB

func (e *Elevator) RequestsAbove() bool {
	//self explainatory, the different buttons are BTN_HALLUP, BTN_HALLDOWN, BTN_HALLCAB
	for f := e.floor + 1; f < Num_floors; f++ {
		for btn := 0; btn < Num_buttons; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func (e *Elevator) RequestsBelow() bool {
	//also self explainatory
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < Num_buttons; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func (e *Elevator) RequestsHere() bool {
	//shouldnt need to explain this either
	for btn := 0; btn < Num_buttons; btn++ {
		if e.requests[e.floor][btn] {
			return true
		}
	}
	return false
}

func (e *Elevator) RequestsChooseDirection() DirBehaviourPair {
	switch e.direction {
	case Dir_up:
		if e.RequestsAbove() {
			return DirBehaviourPair{Dir_up, Behaviour_moving}
		} else if e.RequestsHere() {
			return DirBehaviourPair{Dir_stop, Behaviour_door_open}
		} else if e.RequestsBelow() {
			return DirBehaviourPair{Dir_down, Behaviour_moving}
		} else {
			return DirBehaviourPair{Dir_stop, Behaviour_idle}
		}
	case Dir_down:
		if e.RequestsBelow() {
			return DirBehaviourPair{Dir_down, Behaviour_moving}
		} else if e.RequestsHere() {
			return DirBehaviourPair{Dir_stop, Behaviour_door_open}
		} else if e.RequestsAbove() {
			return DirBehaviourPair{Dir_up, Behaviour_moving}
		} else {
			return DirBehaviourPair{Dir_stop, Behaviour_idle}
		}
	case Dir_stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
		if e.RequestsHere() {
			return DirBehaviourPair{Dir_stop, Behaviour_door_open}
		} else if e.RequestsAbove() {
			return DirBehaviourPair{Dir_up, Behaviour_moving}
		} else if e.RequestsBelow() {
			return DirBehaviourPair{Dir_down, Behaviour_moving}
		} else {
			return DirBehaviourPair{Dir_stop, Behaviour_idle}
		}
	default:
		return DirBehaviourPair{Dir_stop, Behaviour_idle}
	}
}

func (e *Elevator) RequestsShouldStop() bool {
	switch e.direction {
	case Dir_down:
		return e.requests[e.floor][Btn_halldown] || e.requests[e.floor][Btn_hallcab] || !e.RequestsBelow()
	case Dir_up:
		return e.requests[e.floor][Btn_hallup] || e.requests[e.floor][Btn_hallcab] || !e.RequestsAbove()
	default:
		return true
	}
}

func (e *Elevator) RequestsShouldClearImmediately(buttonFloor int, buttonType Button) bool {
	switch e.config.clearRequestVariation {
	case CV_All:
		return e.floor == buttonFloor
	case CV_InDirn:
		return e.floor == buttonFloor && ((e.direction == Dir_up && buttonType == Btn_hallup) ||
			(e.direction == Dir_down && buttonType == Btn_halldown) ||
			e.direction == Dir_stop ||
			buttonType == Btn_hallcab ) 
	default:
		return false
	}
}

func RequestsClearAtCurrentFloor(e Elevator) Elevator {
	var update [3]int
	switch e.config.clearRequestVariation {
	case CV_All:
		for btn := 0; btn < Num_buttons; btn++ {
			e.requests[e.floor][btn] = false
			update = [3]int{e.floor, btn, 0}
			go Transmitt_update_and_update_localHallRequests(update)


		}

	case CV_InDirn:
		e.requests[e.floor][Btn_hallcab] = false
		switch e.direction {
		case Dir_up:
			if !e.RequestsAbove() && !e.requests[e.floor][Btn_hallup] {
				e.requests[e.floor][Btn_halldown] = false
				update = [3]int{e.floor, int(Btn_halldown), 0}
				go Transmitt_update_and_update_localHallRequests(update)

			}
			e.requests[e.floor][Btn_hallup] = false
			update = [3]int{e.floor, int(Btn_hallup), 0}
			go Transmitt_update_and_update_localHallRequests(update)


		case Dir_down:
			if !e.RequestsBelow() && !e.requests[e.floor][Btn_halldown] {
				e.requests[e.floor][Btn_hallup] = false
				update = [3]int{e.floor, int(Btn_hallup), 0}
				go Transmitt_update_and_update_localHallRequests(update)

			}
			e.requests[e.floor][Btn_halldown] = false
			update = [3]int{e.floor, int(Btn_halldown), 0}
			go Transmitt_update_and_update_localHallRequests(update)


		default:
			e.requests[e.floor][Btn_hallup] = false
			update = [3]int{e.floor, int(Btn_hallup), 0}
			go Transmitt_update_and_update_localHallRequests(update)


			e.requests[e.floor][Btn_halldown] = false
			update = [3]int{e.floor, int(Btn_halldown), 0}
			go Transmitt_update_and_update_localHallRequests(update)


		}

	}
	SetAllLights()

	return e
}