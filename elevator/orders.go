package elevator

import (
	"root/sharedData"
	"root/customStructs"
)

func (e *Elevator) RequestsAbove() bool {
    for f := e.floor + 1; f < Num_floors; f++ {
        for btn := 0; btn < Num_buttons; btn++ {
            if e.Requests[f][btn] {
                return true
            }
        }
    }
    return false
}

func (e *Elevator) RequestsBelow() bool {
    for f := 0; f < e.floor; f++ {
        for btn := 0; btn < Num_buttons; btn++ {
            if e.Requests[f][btn] {
                return true
            }
        }
    }
    return false
}

func (e *Elevator) RequestsHere() bool {
    for btn := 0; btn < Num_buttons; btn++ {
        if e.Requests[e.floor][btn] {
            return true
        }
    }
    return false
}

func (e *Elevator) SelectNextDirection() DirBehaviourPair {
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
    case Dir_stop:
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

func (e *Elevator) ShouldStop() bool {
    switch e.direction {
    case Dir_down:
        return e.Requests[e.floor][Btn_halldown] || e.Requests[e.floor][Btn_hallcab] || !e.RequestsBelow()
    case Dir_up:
        return e.Requests[e.floor][Btn_hallup] || e.Requests[e.floor][Btn_hallcab] || !e.RequestsAbove()
    default:
        return true
    }
}

func (e *Elevator) RequestsClearImmediately(buttonFloor int, buttonType Button) bool {
    if e.direction != Dir_stop {
        return false
    }
    
    return e.floor == buttonFloor && (
        (e.direction == Dir_up && buttonType == Btn_hallup) ||
        (e.direction == Dir_down && buttonType == Btn_halldown) ||
        e.direction == Dir_stop ||
        buttonType == Btn_hallcab)
}

func (e *Elevator) ClearRequestsAtFloor(SharedData *sharedData.SharedData) []customStructs.Update {

    updates := []customStructs.Update{}

    e.Requests[e.floor][Btn_hallcab] = false
    switch e.direction {
    case Dir_up:
        if !e.RequestsAbove() && !e.Requests[e.floor][Btn_hallup] {
            e.Requests[e.floor][Btn_halldown] = false
            updates = append(updates, customStructs.Update{Floor: e.floor, ButtonType: int(Btn_halldown), Value: false})
        }
        e.Requests[e.floor][Btn_hallup] = false
        updates = append(updates, customStructs.Update{Floor: e.floor, ButtonType: int(Btn_hallup), Value: false})
    
    case Dir_down:
        if !e.RequestsBelow() && !e.Requests[e.floor][Btn_halldown] {
            e.Requests[e.floor][Btn_hallup] = false
            updates = append(updates, customStructs.Update{Floor: e.floor, ButtonType: int(Btn_hallup), Value: false})
        }
        e.Requests[e.floor][Btn_halldown] = false
        updates = append(updates, customStructs.Update{Floor: e.floor, ButtonType: int(Btn_halldown), Value: false})
    default:
            e.Requests[e.floor][Btn_hallup] = false
            updates = append(updates, customStructs.Update{Floor: e.floor, ButtonType: int(Btn_hallup), Value: false})
        
            e.Requests[e.floor][Btn_halldown] = false
            updates = append(updates, customStructs.Update{Floor: e.floor, ButtonType: int(Btn_halldown), Value: false})
    }
    return updates
}
