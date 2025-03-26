package elevator

import (

	"fmt"
	"root/assigner"
	"root/config"
	"root/elevio"
	"root/sharedData"
	"root/transmitter"
)



func UpdatesharedHallRequests(elevator *Elevator, sharedData *sharedData.SharedData, update config.Update) {
    if  update.ButtonType != 2 { // ignores updates to cab requests (update[1] != 2)
        sharedData.HallRequests[update.Floor][update.ButtonType] = update.Value
    }
}

func AssignLocalHallRequests(elevator *Elevator, SharedData *sharedData.SharedData, externalConn sharedData.ExternalConn) {
    localData := GetElevatorData(elevator)
    remoteData := SharedData.RemoteElevatorData
    sharedHallRequests := SharedData.HallRequests

    correctedLocalData := localData
    
    // Prevents invalid data from crashing the assigner
    if localData.Floor == 0 && localData.Direction == "down" || localData.Floor == 3 && localData.Direction == "up" {
        correctedLocalData.Direction = "stop"

    } 
    
    updatedRequests := assigner.Assigner(correctedLocalData, remoteData, sharedHallRequests, externalConn)
    elevator.Requests = MakeRequests(updatedRequests, GetCabRequests(elevator.Requests))
}




func Send_Elevator_data(elevatorData config.Elevator_data, externalConn *sharedData.ExternalConn) {
    transmitter.Send_Elevator_data(elevatorData, externalConn)
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
        if elevio.MotorDirection(elevator.direction) != 0 {
        StartStuckTimer()
        fmt.Println("start if idle")
    }

    }
}