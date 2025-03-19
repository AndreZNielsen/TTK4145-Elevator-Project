package elevator

import (
	"fmt"
	"root/sharedData"
)

func (e *Elevator) RequestsAbove() bool {
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
    for btn := 0; btn < Num_buttons; btn++ {
        if e.requests[e.floor][btn] {
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
        return e.requests[e.floor][Btn_halldown] || e.requests[e.floor][Btn_hallcab] || !e.RequestsBelow()
    case Dir_up:
        return e.requests[e.floor][Btn_hallup] || e.requests[e.floor][Btn_hallcab] || !e.RequestsAbove()
    default:
        return true
    }
}

func (e *Elevator) RequestsShouldClearImmediately(buttonFloor int, buttonType Button) bool {
    if e.direction != Dir_stop {
        return false
    }
    fmt.Println(e.direction)
    
    switch e.config.clearRequestVariation {
    case CV_All:
        return e.floor == buttonFloor
    case CV_InDirn:
        return e.floor == buttonFloor && (
            (e.direction == Dir_up && buttonType == Btn_hallup) ||
            (e.direction == Dir_down && buttonType == Btn_halldown) ||
            e.direction == Dir_stop ||
            buttonType == Btn_hallcab)
    default:
        return false
    }
}


// If this function can end up needing to send two messages, were kind of screwed.
// Need to make some big changes then.
func (e *Elevator) RequestsClearAtCurrentFloor(externalData *sharedData.ExternalData) [3]int {
    switch e.config.clearRequestVariation {
    // case CV_All: // Is this even possible?
    //     for btn := 0; btn < Num_buttons; btn++ {
    //         e.requests[e.floor][btn] = false
    //         //UpdateAndTransmittLocalRequests(e, e.floor, Button(btn), 0, externalData)
    //         return [3]int{e.floor, btn, 0}
    //     }
    case CV_InDirn:
        e.requests[e.floor][Btn_hallcab] = false
        switch e.direction {
        case Dir_up:
            if !e.RequestsAbove() && !e.requests[e.floor][Btn_hallup] {
                e.requests[e.floor][Btn_halldown] = false
                //UpdateAndTransmittLocalRequests(e, e.floor, Btn_halldown, 0, externalData)
                return [3]int{e.floor, int(Btn_halldown), 0}
            }
            e.requests[e.floor][Btn_hallup] = false
            //UpdateAndTransmittLocalRequests(e, e.floor, Btn_hallup, 0, externalData)
            return [3]int{e.floor, int(Btn_hallup), 0}
        case Dir_down:
            if !e.RequestsBelow() && !e.requests[e.floor][Btn_halldown] {
                e.requests[e.floor][Btn_hallup] = false
                //UpdateAndTransmittLocalRequests(e, e.floor, Btn_hallup, 0, externalData)
                return [3]int{e.floor, int(Btn_hallup), 0}
            }
            e.requests[e.floor][Btn_halldown] = false
            //UpdateAndTransmittLocalRequests(e, e.floor, Btn_halldown, 0, externalData)
            return [3]int{e.floor, int(Btn_halldown), 0}
        // default:
        //     e.requests[e.floor][Btn_hallup] = false
        //     //UpdateAndTransmittLocalRequests(e, e.floor, Btn_hallup, 0, externalData)
        //     return [3]int{e.floor, int(Btn_hallup), 0}
            
        //     e.requests[e.floor][Btn_halldown] = false
        //     //UpdateAndTransmittLocalRequests(e, e.floor, Btn_halldown, 0, externalData)
        //     return [3]int{e.floor, int(Btn_halldown), 0}
        }
    }
    return [3]int{0, 0, 0} // This should never happen!

    
    a // just indicating that this function needs fixing
}