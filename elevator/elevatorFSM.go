package elevator

import (
    Config "root/config"
    "root/elevio"

)

type LocalEvent struct {
    EventType  string
    Button     elevio.ButtonEvent
    Floor      int
    Obstructed bool
}

func FSM_MakeElevator(elevator *Elevator, elevator_ip string, Num_floors int) {
    elevio.Init(elevator_ip, Num_floors)
    *elevator = MakeUninitializedelevator()
    FSM_InitBetweenFloors(elevator) 
}

func GetElevatorData(elevator *Elevator) Config.Elevator_data {
    return Config.Elevator_data{
        Behavior:    EbToString(elevator.behaviour), 
        Floor:       elevator.floor, 
        Direction:   ElevioDirToString(elevator.direction), 
        CabRequests: GetCabRequests(elevator.requests),
    }
}

func FSM_InitBetweenFloors(elevator *Elevator) {
    elevio.SetMotorDirection(elevio.MD_Down)
    elevator.direction = Dir_down
    elevator.behaviour = Behaviour_moving
}

func FSM_RequestButtonPress(elevator *Elevator, btn_floor int, btn_type Button) {

    switch elevator.behaviour {
    case Behaviour_door_open:
        if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
            StartTimer()
        } else {
            if btn_type == Btn_hallcab {
                elevator.requests[btn_floor][btn_type] = true
            }
            UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1)
        }
    case Behaviour_moving:   
        if btn_type == Btn_hallcab {
            elevator.requests[btn_floor][btn_type] = true
        }
        UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1)

    case Behaviour_idle:
        if btn_type == Btn_hallup {
            elevator.requests[btn_floor][btn_type] = true
        }
        UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1)

        if elevator.floor == btn_floor {
            elevio.SetMotorDirection(elevio.MD_Stop)
            elevio.SetDoorOpenLamp(true)
            elevator.RequestsClearAtCurrentFloor()
            StartTimer()
            SetAllLights(elevator)
            elevator.behaviour = Behaviour_door_open
        } else {
            UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1)
        }
    }
}

func UpdateAndTransmittLocalRequests(elevator *Elevator, btn_floor int, btn_type Button, update int) {
    localUpdate := [3]int{btn_floor, int(btn_type), update}
    go Transmitt_update_and_update_localHallRequests(elevator, localUpdate)
}
func FSM_FloorArrival(elevator *Elevator, newFloor int) {
    elevator.floor = newFloor
    elevio.SetFloorIndicator(elevator.floor)

    switch elevator.behaviour {
    case Behaviour_moving:
        if elevator.RequestsShouldStop() {
            elevio.SetMotorDirection(elevio.MD_Stop)
            elevio.SetDoorOpenLamp(true)
            elevator.RequestsClearAtCurrentFloor()
            StartTimer()
            SetAllLights(elevator)
            elevator.behaviour = Behaviour_door_open
        }
    }
    go Send_Elevator_data(GetElevatorData(elevator))
}

func FSM_DoorTimeout(elevator *Elevator) {
    switch elevator.behaviour {
    case Behaviour_door_open:
        pair := elevator.SelectNextDirection()
        elevator.direction = pair.dir
        elevator.behaviour = pair.behaviour

        switch elevator.behaviour {
        case Behaviour_door_open:
            StartTimer()
            elevator.RequestsClearAtCurrentFloor()
            SetAllLights(elevator)
        case Behaviour_moving, Behaviour_idle:
            elevio.SetDoorOpenLamp(false)
            elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
        }
    }
    go Send_Elevator_data(GetElevatorData(elevator))
}

func FSM_HandleLocalEvent(elevator *Elevator, event LocalEvent) {
    switch event.EventType {
    case "button":
        FSM_RequestButtonPress(elevator, event.Button.Floor, Button(event.Button.Button))
        SetAllLights(elevator)
    case "floor":
        if !IsDoorObstructed() {
            FSM_FloorArrival(elevator, event.Floor)
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
            FSM_DoorTimeout(elevator)
        } else {
            StartTimer()
        }
    }
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
