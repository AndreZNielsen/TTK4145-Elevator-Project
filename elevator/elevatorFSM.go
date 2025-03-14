package elevator

import (
    Config "root/config"
    "root/elevio"
    "root/sharedData"
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
    return Config.Elevator_data{Behavior: EbToString(elevator.behaviour), Floor: elevator.floor, Direction: ElevioDirToString(elevator.direction), CabRequests: GetCabRequests(elevator.requests)}
}

func FSM_InitBetweenFloors(elevator *Elevator) {
    elevio.SetMotorDirection(elevio.MD_Down)
    elevator.direction = Dir_down
    elevator.behaviour = Behaviour_moving
}

func FSM_RequestButtonPress(elevator *Elevator, btn_floor int, btn_type Button) {
    var localUpdate [3]int

    switch elevator.behaviour {
    case Behaviour_door_open:
        if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
            StartTimer()
        } else {
            if btn_type == Btn_hallcab {
                elevator.requests[btn_floor][btn_type] = true
            }
            localUpdate = [3]int{btn_floor, int(btn_type), 1}
            go Transmitt_update_and_update_localHallRequests(elevator, localUpdate)
        }
    case Behaviour_moving:
        if btn_type == Btn_hallcab {
            elevator.requests[btn_floor][btn_type] = true
        }
        localUpdate = [3]int{btn_floor, int(btn_type), 1}
        go Transmitt_update_and_update_localHallRequests(elevator, localUpdate)

    case Behaviour_idle:
        if btn_type == Btn_hallup {
            elevator.requests[btn_floor][btn_type] = true
        }

        if elevator.floor == btn_floor {
            elevio.SetMotorDirection(elevio.MD_Stop)
            elevio.SetDoorOpenLamp(true)
            elevator.RequestsClearAtCurrentFloor()
            StartTimer()
            SetAllLights(elevator)
            elevator.behaviour = Behaviour_door_open
        } else {
            localUpdate = [3]int{btn_floor, int(btn_type), 1}
            go Transmitt_update_and_update_localHallRequests(elevator, localUpdate)
        }
    }
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
        pair := elevator.RequestsChooseDirection()
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
    buttonEvents := make(chan elevio.ButtonEvent)
    floorEvents := make(chan int)
    obstructionEvents := make(chan bool)
    timerEvents := make(chan bool)

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

var doorObstructed bool

func DoorObstructed(elevator *Elevator) {
    doorObstructed = true
    if elevator.behaviour == Behaviour_door_open {
        StartTimer()
    }
}

func DoorUnobstructed(elevator *Elevator) {
    doorObstructed = false
    if elevator.behaviour == Behaviour_door_open {
        StartTimer()
    }
}

func IsDoorObstructed() bool {
    return doorObstructed
}

func GetCabRequests(matrix [Num_floors][3]bool) []bool {
    var column []bool
    for i := 0; i < len(matrix); i++ {
        column = append(column, matrix[i][2])
    }
    return column
}

func GetHallRequests(matrix [Num_floors][3]bool) [][2]bool {
    var newMatrix [][2]bool

    // Extract columns 1 and 2 (index 0 and 1)
    for i := 0; i < len(matrix); i++ {
        newMatrix = append(newMatrix, [2]bool{matrix[i][0], matrix[i][1]})
    }
    return newMatrix
}

func MakeRequests(HallRequests [][2]bool, GetCabRequests []bool) [Num_floors][3]bool {
    var result [Num_floors][3]bool

    for i := 0; i < Num_floors; i++ {
        result[i][0] = HallRequests[i][0]
        result[i][1] = HallRequests[i][1]
        result[i][2] = GetCabRequests[i]
    }
    return result
}

func GetElevator(elevator *Elevator) Elevator {
    return *elevator
}

func SetAllLights(elevator *Elevator) {
    //Basically just takes the requests from the button presses and lights up the corresponding button lights
    requests := MakeRequests(sharedData.GetsharedHallRequests(), GetCabRequests(elevator.requests))
    for floor := 0; floor < Num_floors; floor++ {
        for btn := 0; btn < Num_buttons; btn++ {
            elevio.SetButtonLamp(elevio.ButtonType(btn), floor, requests[floor][btn])
        }
    }
}