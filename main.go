package main

import (
    "fmt"
    "root/elevator"
    "root/network"
    "root/reciver"
    "root/transmitter"
    "root/config"
)

var elevator_1_ip = "localhost:12345"

func main() {
    fmt.Println("Started!")

    localEventRecived 	:= make(chan elevator.LocalEvent)
    aliveTimer 			:= make(chan bool)
    remoteEventRecived 	:= make(chan [3]int)
    disconnected 		:= make(chan string)

    var elev elevator.Elevator
    elevator.FSM_MakeElevator(&elev, elevator_1_ip, config.Num_floors)
    elevator.Start_if_idle(&elev)
    go elevator.FSM_DetectLocalEvents(localEventRecived)

    network.Start_network(remoteEventRecived, disconnected)
    transmitter.Send_Elevator_data(elevator.GetElevatorData(&elev))
    go reciver.AliveTimer(aliveTimer)

    for {
        select {
        case localEvent := <-localEventRecived:
            elevator.FSM_HandleLocalEvent(&elev, localEvent)

        case remoteEvent := <-remoteEventRecived:
            elevator.UpdatesharedHallRequests(&elev, remoteEvent)
            elevator.ChangeLocalHallRequests(&elev)
            elevator.SetAllLights(&elev)

        case id := <-disconnected:
            go network.Network_reconnector(remoteEventRecived, disconnected, id)

        case <-aliveTimer:

        }
    }
}