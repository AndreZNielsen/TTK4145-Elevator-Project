package elevator

import (
	"fmt"
	"root/assigner"
	Config "root/config"
	"root/elevio"
	"root/sharedData"
	"root/transmitter"
)

func UpdatesharedHallRequests(elevator *Elevator, externalData *sharedData.ExternalData, update [3]int) {
    sharedHallRequests := externalData.HallRequests
    if update[2] == 1 && update[1] != 2 { // ignores updates to cab requests (update[1] != 2)
        sharedHallRequests[update[0]][update[1]] = true
    } else if update[1] != 2 {
        sharedHallRequests[update[0]][update[1]] = false
    }
  
    AssignLocalHallRequests(elevator, externalData)      //This one is called anyway should be called elsewhere
}

func Transmitt_update_and_update_localHallRequests(elevator *Elevator, update_val [3]int, externalData *sharedData.ExternalData) { // sends the hall requests update to the other elevator and updates the local hall requests
    UpdatesharedHallRequests(elevator, externalData, update_val)     // call this in main instead, as it requires externalData
    // transmitter.Send_update(update_val, externalData)
}

func AssignLocalHallRequests(elevator *Elevator, externalData *sharedData.ExternalData) {
    localData := GetElevatorData(elevator)
    remoteData := externalData.RemoteElevatorData
    sharedHallRequests := externalData.HallRequests

    fmt.Println(localData)
    fmt.Println(remoteData)

//////////////////////////////////////////////////////////////////////////////////////////////  
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
    
/////////////////////////////////////////////////////////////////////////////////////////////////

    updatedRequests := assigner.Assigner(correctedLocalData, remoteData, sharedHallRequests)
    elevator.requests = MakeRequests(updatedRequests, GetCabRequests(elevator.requests))

    Start_if_idle(elevator) // I have a problem with the way this is called. WIll look into it later


    // if localData.Floor != -1 && !(localData.Floor == 0 && localData.Direction == "down") && !(localData.Floor == 3 && localData.Direction == "up") {
    //     updatedRequests := assigner.Assigner(localData, remoteData, sharedHallRequests)
    //     elevator.requests = MakeRequests(updatedRequests, GetCabRequests(elevator.requests))

    //     Start_if_idle(elevator) // I have a problem with the way this is called. WIll look into it later
    //     // SetAllLights(elevator, &externalData)
    //     elevator.print()
    // } else {
    //     fmt.Println("Invalid data sent to assigner executealbe!")
    // }
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
        if elevator.behaviour == Behaviour_door_open {
            DoorOpen(elevator)
        }
        elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
    }
}