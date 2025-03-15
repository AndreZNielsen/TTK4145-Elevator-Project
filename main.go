package main

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/network"
	"root/reciver"
	"root/sharedData"
	"root/transmitter"
)

var elevator_1_ip = "localhost:15657"

func main() {
    fmt.Println("Started!")

    localEventRecived 	:= make(chan elevator.LocalEvent)
    aliveTimer 			:= make(chan bool)
    remoteEventRecived 	:= make(chan [3]int)
    disconnected 		:= make(chan string)

	externalData := sharedData.InitExternalData()

    var elev elevator.Elevator
    elevator.FSM_MakeElevator(&elev, elevator_1_ip, config.Num_floors)
    elevator.Start_if_idle(&elev)
    go elevator.FSM_DetectLocalEvents(localEventRecived)

    network.Start_network(remoteEventRecived, disconnected, externalData)       // We could separate connections from the other shared data? Maybe an idea. OR only pass externalData.RemoteElevatorConnections!
    transmitter.Send_Elevator_data(elevator.GetElevatorData(&elev), externalData)
    go reciver.AliveTimer(aliveTimer)

    for {
        select {
        case localEvent := <-localEventRecived:
            elevator.FSM_HandleLocalEvent(&elev, localEvent, externalData)
			elevator.SetAllLights(&elev, externalData)
			// Transmitt ? No, because we only transmitt changes. It would not be possible to put it here. 
			// assign here? Cause it requires the external data

        case remoteEvent := <-remoteEventRecived:
			elevator.FSM_HandleRemoteEvent(&elev, externalData, remoteEvent)
            

        case id := <-disconnected:
            go network.Network_reconnector(remoteEventRecived, disconnected, id, externalData)

        case <-aliveTimer:

        }
    }
}