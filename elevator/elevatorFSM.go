package elevator

import (
	"root/elevio"
	"root/sharedData"
)

type LocalEvent struct {
    EventType  string
    Button     elevio.ButtonEvent
    Floor      int
    Obstructed bool
}

// type RemoteEvent struct{ // This might be a good idea to implement
//     Floor int
//     Button int
//     Update int
// }

func FSM_MakeElevator(elevator *Elevator, elevator_ip string, Num_floors int) {
    elevio.Init(elevator_ip, Num_floors)
    *elevator = MakeUninitializedelevator()
    FSM_InitBetweenFloors(elevator) 
}



func FSM_InitBetweenFloors(elevator *Elevator) { // Create Move-down function
    elevio.SetMotorDirection(elevio.MD_Down)
    elevator.direction = Dir_down
    elevator.behaviour = Behaviour_moving
}

func FSM_RequestButtonPress(elevator *Elevator, btn_floor int, btn_type Button, externalData *sharedData.ExternalData) {
// Could look something like:

// If clearImmediately : return // Does the light need to turn on, even if the request is cleared immediately?

// Update local requests
// Transmitt update
// assign
// SetLights
// Start_if_idle // We kind of have two control functions at the moment? This and start_if_idle?
    switch elevator.behaviour {
    case Behaviour_door_open:
        if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
            StartTimer()
        } else {
            if btn_type == Btn_hallcab {
                elevator.requests[btn_floor][btn_type] = true
            }
            UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1, externalData)
        }
    case Behaviour_moving:   
        if btn_type == Btn_hallcab {
            elevator.requests[btn_floor][btn_type] = true
        }
        UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1, externalData)

    case Behaviour_idle:
        if btn_type == Btn_hallup {
            elevator.requests[btn_floor][btn_type] = true
        }
        UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1, externalData)

        if elevator.floor == btn_floor {
            elevio.SetMotorDirection(elevio.MD_Stop)
            elevio.SetDoorOpenLamp(true)
            elevator.RequestsClearAtCurrentFloor(externalData)
            StartTimer()
            //SetAllLights(elevator)
            elevator.behaviour = Behaviour_door_open
        } else {
            UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1, externalData)
        }
    }
}

func UpdateAndTransmittLocalRequests(elevator *Elevator, btn_floor int, btn_type Button, update int, externalData *sharedData.ExternalData) {
    localUpdate := [3]int{btn_floor, int(btn_type), update}
    go Transmitt_update_and_update_localHallRequests(elevator, localUpdate, externalData)
}


// ControlMovement-function, start_if_idle is not good enough. We can make one that is also capable of stopping

func FSM_FloorArrival(elevator *Elevator, newFloor int, externalData *sharedData.ExternalData) {
    elevator.floor = newFloor
    elevio.SetFloorIndicator(elevator.floor)

    switch elevator.behaviour {
    case Behaviour_moving:
        if elevator.RequestsShouldStop() {
            elevio.SetMotorDirection(elevio.MD_Stop)
            elevio.SetDoorOpenLamp(true)
            elevator.RequestsClearAtCurrentFloor(externalData)
            StartTimer()
            //SetAllLights(elevator)
            elevator.behaviour = Behaviour_door_open
        }
    }
    go Send_Elevator_data(GetElevatorData(elevator), externalData)
}

func FSM_DoorTimeout(elevator *Elevator, externalData *sharedData.ExternalData) {
    switch elevator.behaviour {
    case Behaviour_door_open:
        pair := elevator.SelectNextDirection()
        elevator.direction = pair.dir
        elevator.behaviour = pair.behaviour

        switch elevator.behaviour {
        case Behaviour_door_open:
            StartTimer()
            elevator.RequestsClearAtCurrentFloor(externalData)
            //SetAllLights(elevator)
        case Behaviour_moving, Behaviour_idle:
            elevio.SetDoorOpenLamp(false)
            elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
        }
    }
    go Send_Elevator_data(GetElevatorData(elevator), externalData)
}


// Here I think it would be reasonable to create functions for each case
// like FSM_HandleButtonEvent, FSM_HandleFloorEvent, FSM_HandleObstructedEvent, FSM_HandleTimerEvent
// The cases are very different from one another, so I think this makes sense.

// FSM_HandleButtonEvent should do what FSM_RequestButtonPress does now, but in an organized way
func FSM_HandleLocalEvent(elevator *Elevator, event LocalEvent, externalData *sharedData.ExternalData) {
    switch event.EventType {
    case "button":
        FSM_RequestButtonPress(elevator, event.Button.Floor, Button(event.Button.Button), externalData)
        //SetAllLights(elevator)
    case "floor":
        if !IsDoorObstructed() {
            FSM_FloorArrival(elevator, event.Floor, externalData)
        }
    case "obstructed":
        if event.Obstructed {
            DoorObstructed(elevator)
        } else {
            DoorUnobstructed(elevator)
        }
    case "timer":
        if !IsDoorObstructed() {
            StopTimer()
            FSM_DoorTimeout(elevator, externalData)
        } else {
            StartTimer()
        }
    }
}

func FSM_HandleRemoteEvent(elevator *Elevator, externalData *sharedData.ExternalData, event [3]int) { // Ideally this should say RemoteEvent, instead of [3]int, maybe fix this later
    UpdatesharedHallRequests(elevator, externalData, event)
    ChangeLocalHallRequests(elevator, externalData)
    SetAllLights(elevator, externalData)
    //Start_if_idle(elevator) // should be called here instead of in ChangeLocalHallRequests    

    // Once this change is made I am very happy with this function
}


func FSM_DetectLocalEvents(localEvents chan<- LocalEvent) {
    buttonEvents        := make(chan elevio.ButtonEvent)
    floorEvents         := make(chan int)
    obstructionEvents   := make(chan bool)
    timerEvents         := make(chan bool)

    go elevio.PollButtons(buttonEvents)
    go elevio.PollFloorSensor(floorEvents)
    go elevio.PollObstructionSwitch(obstructionEvents)
    go TimerIsDone(timerEvents)

    for {
        select {
        case button := <-buttonEvents:
            localEvents <- LocalEvent{EventType: "button", Button: button}
        case floor := <-floorEvents:
            localEvents <- LocalEvent{EventType: "floor", Floor: floor}
        case obstructed := <-obstructionEvents:
            localEvents <- LocalEvent{EventType: "obstructed", Obstructed: obstructed}
        case <-timerEvents:
            localEvents <- LocalEvent{EventType: "timer"}
        }
    }
}
