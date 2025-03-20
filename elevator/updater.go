package elevator

import (
	"fmt"
	"root/assigner"
	"root/config"
	"root/elevio"
	"root/sharedData"
	"root/transmitter"
)



func UpdatesharedHallRequests(elevator *Elevator, externalData *sharedData.ExternalData, update config.Update) {
    if update.Value && update.ButtonType != 2 { // ignores updates to cab requests (update[1] != 2)
        externalData.HallRequests[update.Floor][update.ButtonType] = true
    } else if update.ButtonType != 2 {
        externalData.HallRequests[update.Floor][update.ButtonType] = false
    }
  
    
}

func AssignLocalHallRequests(elevator *Elevator, externalData *sharedData.ExternalData) {
    localData := GetElevatorData(elevator)
    remoteData := externalData.RemoteElevatorData
    sharedHallRequests := externalData.HallRequests

    fmt.Println(localData)
    fmt.Println(remoteData)

    // Prevents invalid data from crashing the assigner
    // We might be able to find a better solution here.
    // The same problem can happen with data from remoteElevators. That will require a separate fix. Probably similar, but not implemented yet
    fmt.Println("Direction: ",localData.Direction)
    fmt.Println("Behaviour: ", localData.Behavior)

    correctedLocalData := localData
    
    if localData.Floor == 0 && localData.Direction == "down" || localData.Floor == 3 && localData.Direction == "up" {
        fmt.Println("Invalid data sent to assigner executealbe, hard coded fix triggered!")
        correctedLocalData.Direction = "stop"

    } 
    
    updatedRequests := assigner.Assigner(correctedLocalData, remoteData, sharedHallRequests)
    elevator.requests = MakeRequests(updatedRequests, GetCabRequests(elevator.requests))
}




func Send_Elevator_data(elevatorData config.Elevator_data, externalData *sharedData.ExternalData) {
    transmitter.Send_Elevator_data(elevatorData, externalData)
}

func Start_if_idle(elevator *Elevator) {

    if elevator.behaviour == Behaviour_idle {
        pair := elevator.SelectNextDirection()
        elevator.direction = pair.dir
        elevator.behaviour = pair.behaviour
        if elevator.behaviour == Behaviour_door_open {
            DoorOpen(elevator)
        }
        elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
    }
}