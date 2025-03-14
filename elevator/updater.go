package elevator

import (
    "root/assigner"
    "root/sharedData"
    "root/transmitter"
    "root/elevio"
    Config "root/config"
    "fmt"
)

func UpdatesharedHallRequests(elevator *Elevator, update [3]int) {
    sharedHallRequests := sharedData.GetsharedHallRequests()
    if update[2] == 1 && update[1] != 2 { // ignores updates to cab requests (update[1] != 2)
        sharedHallRequests[update[0]][update[1]] = true
    } else if update[1] != 2 {
        sharedHallRequests[update[0]][update[1]] = false
    }
    sharedData.ChangeSharedHallRequests(sharedHallRequests)
    ChangeLocalHallRequests(elevator)
}

func Transmitt_update_and_update_localHallRequests(elevator *Elevator, update_val [3]int) { // sends the hall requests update to the other elevator and updates the local hall requests
    UpdatesharedHallRequests(elevator, update_val)
    transmitter.Send_update(update_val)
}

func ChangeLocalHallRequests(elevator *Elevator) {
    fmt.Println(GetElevatorData(elevator))
    fmt.Println(sharedData.GetRemoteElevatorData())
    if GetElevatorData(elevator).Floor != -1 && !(GetElevatorData(elevator).Floor == 0 && GetElevatorData(elevator).Direction == "down") && !(GetElevatorData(elevator).Floor == 3 && GetElevatorData(elevator).Direction == "up") { // stops the elevator data from crashing the assigner
        elevator.requests = MakeRequests(assigner.Assigner(GetElevatorData(elevator), sharedData.GetRemoteElevatorData(), sharedData.GetsharedHallRequests()), GetCabRequests(elevator.requests))
        Start_if_idle(elevator)
        SetAllLights(elevator)
        elevator.print()
    }
}

func Send_Elevator_data(elevatorData Config.Elevator_data) {
    transmitter.Send_Elevator_data(elevatorData)
}

func Start_if_idle(elevator *Elevator) {
    switch elevator.behaviour {
    case Behaviour_idle:
        pair := elevator.RequestsChooseDirection()
        elevator.direction = pair.dir
        elevator.behaviour = pair.behaviour
        elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
    case Behaviour_door_open:
        StartTimer()
        elevio.SetDoorOpenLamp(true)
    }
}