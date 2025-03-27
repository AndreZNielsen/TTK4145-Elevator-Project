package main

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/network"
	"flag"
    "root/customStructs"
	"root/reviver"
	SharedData "root/sharedData"
	"root/transmitter"
)



func main() {
    fmt.Println("Started!")
    // Define flags
	var isRestart bool
	var cabBackup string

	// Parse the flags
	flag.BoolVar(&isRestart, "isRestart", false, "Indicates if this is a restart")
	flag.StringVar(&cabBackup, "cabBackup", "", "Space-separated list for CabBackup")
	flag.Parse()
    
    var elev elevator.Elevator
    
    localEventRecived 	:= make(chan elevator.LocalEvent)
    remoteEventRecived 	:= make(chan customStructs.RemoteEvent)
    disconnected 		:= make(chan string)

	
	sharedData := SharedData.InitSharedData()
    externalConn := SharedData.InitExternalConn()

    elevator.FSM_MakeElevator(&elev, config.LocalElevatorServerPort, config.Num_floors)

	go network.StartPeerNetwork(remoteEventRecived, disconnected, sharedData, externalConn,&elev)

    go reviver.StartReviver(&elev)

    go elevator.FSM_DetectLocalEvents(localEventRecived)
    
    if isRestart{
        elevator.RestorCabRequests(&elev,cabBackup)
        transmitter.MergeHallRequests(externalConn, sharedData.HallRequests, config.RemoteIDs[0])
    }

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
                go network.ReconnectPeer(remoteEventRecived, disconnected, id, sharedData, externalConn,&elev)
            }
        }
    }
}

