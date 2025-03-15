package elevator

import (
    "root/assigner"
    "root/sharedData"
    "root/transmitter"
    "root/elevio"
    Config "root/config"
    "fmt"
)

func UpdatesharedHallRequests(elevator *Elevator, externalData *sharedData.ExternalData, update [3]int) {
    sharedHallRequests := externalData.HallRequests
    if update[2] == 1 && update[1] != 2 { // ignores updates to cab requests (update[1] != 2)
        sharedHallRequests[update[0]][update[1]] = true
    } else if update[1] != 2 {
        sharedHallRequests[update[0]][update[1]] = false
    }
  
    // ChangeLocalHallRequests(elevator)      This one is called anyway!?
}

func Transmitt_update_and_update_localHallRequests(elevator *Elevator, update_val [3]int, externalData *sharedData.ExternalData) { // sends the hall requests update to the other elevator and updates the local hall requests
    UpdatesharedHallRequests(elevator, externalData, update_val)     // call this in main instead, as it requires externalData
    transmitter.Send_update(update_val, externalData)
}

func ChangeLocalHallRequests(elevator *Elevator, externalData *sharedData.ExternalData) {
    localData := GetElevatorData(elevator)
    remoteData := externalData.RemoteElevatorData
    sharedHallRequests := externalData.HallRequests

    fmt.Println(localData)
    fmt.Println(remoteData)

    // Prevents invalid data from crashing the assigner
    if localData.Floor != -1 && !(localData.Floor == 0 && localData.Direction == "down") && !(localData.Floor == 3 && localData.Direction == "up") {
        updatedRequests := assigner.Assigner(localData, remoteData, sharedHallRequests)
        elevator.requests = MakeRequests(updatedRequests, GetCabRequests(elevator.requests))

        Start_if_idle(elevator) // I have a problem with the way this is called. WIll look into it later
        // SetAllLights(elevator, &externalData)
        elevator.print()
    }
}

func Send_Elevator_data(elevatorData Config.Elevator_data, externalData *sharedData.ExternalData) {
    transmitter.Send_Elevator_data(elevatorData, externalData)
}

func Start_if_idle(elevator *Elevator) {
    switch elevator.behaviour {
    case Behaviour_idle:
        pair := elevator.SelectNextDirection()
        elevator.direction = pair.dir
        elevator.behaviour = pair.behaviour
        elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
    case Behaviour_door_open:
        StartTimer()
        elevio.SetDoorOpenLamp(true)
    }
}