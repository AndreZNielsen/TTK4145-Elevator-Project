package main

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/network"
	//"root/reciver"
	"flag"
	"root/reviver"
	SharedData "root/sharedData"
	"root/transmitter"
)

var elevator_1_ip = "localhost:12345"

/*
hvordan kjøre:
start to simulatorer med port 12345 og 12346 (./SimElevatorServer --port ______ i simulator mappen)
kjør go run -ldflags="-X root/config.Elevator_id=A" main.go
og så go run -ldflags="-X root/config.Elevator_id=B" main2.go
på samme maskin
*/

func main() {
    fmt.Println("Started!")
    // Define flags
	var isRestart bool
	var cabBackup string

	// Parse the flags
	flag.BoolVar(&isRestart, "isRestart", false, "Indicates if restart is required")
	flag.StringVar(&cabBackup, "cabBackup", "", "Space-separated list for CabBackup")
	flag.Parse()

    var elev elevator.Elevator

    localEventRecived 	:= make(chan elevator.LocalEvent)
    remoteEventRecived 	:= make(chan config.RemoteEvent)
    disconnected 		:= make(chan string)

	
	sharedData := SharedData.InitSharedData()
    externalConn := SharedData.InitExternalConn()
	network.StartPeerNetwork(remoteEventRecived, disconnected, sharedData, externalConn)
    
    
    elevator.FSM_MakeElevator(&elev, elevator_1_ip, config.Num_floors)
    go elevator.FSM_DetectLocalEvents(localEventRecived)
    fmt.Println(cabBackup)
    if isRestart{
        elevator.RestorCabRequests(&elev,cabBackup)

        transmitter.RequestHallRequests(externalConn, sharedData.HallRequests, config.RemoteIDs[0])

    }
    
    go reviver.StartReviver(&elev)

    transmitter.Send_Elevator_data(elevator.GetElevatorData(&elev), externalConn) 
    for {
        select {
        case localEvent := <-localEventRecived:
            elevator.FSM_HandleLocalEvent(&elev, localEvent, sharedData, externalConn)
			elevator.SetAllLights(&elev, sharedData)
            elevator.Send_Elevator_data(elevator.GetElevatorData(&elev), externalConn)

        case remoteEvent := <-remoteEventRecived:
			elevator.FSM_HandleRemoteEvent(&elev, sharedData, remoteEvent, *externalConn)
            

        case id := <-disconnected:     
            if externalConn.ConnectedConn[id]{
                externalConn.ConnectedConn[id]=false
                network.StopAliveTimer(id)
                fmt.Println("disconnect triggered")

                go network.ReconnectPeer(remoteEventRecived, disconnected, id, sharedData, externalConn,&elev)
            }
        }
    }
}

