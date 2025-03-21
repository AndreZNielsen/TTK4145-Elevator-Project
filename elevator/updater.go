package elevator

import (
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


    // Prevents invalid data from crashing the assigner
    // We might be able to find a better solution here.
    // The same problem can happen with data from remoteElevators. That will require a separate fix. Probably similar, but not implemented yet
    //fmt.Println("Direction: ",localData.Direction)
    //fmt.Println("Behaviour: ", localData.Behavior)

    correctedLocalData := localData
    
    if localData.Floor == 0 && localData.Direction == "down" || localData.Floor == 3 && localData.Direction == "up" {
        //fmt.Println("Invalid data sent to assigner executealbe, hard coded fix triggered!")
        correctedLocalData.Direction = "stop"

    } 
    
    updatedRequests := assigner.Assigner(correctedLocalData, remoteData, sharedHallRequests, externalConn)
    elevator.requests = MakeRequests(updatedRequests, GetCabRequests(elevator.requests))
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
    }
}