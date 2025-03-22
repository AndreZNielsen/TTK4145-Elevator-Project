package main

import (
	"fmt"
    "root/config"
	"root/elevator"
	"root/network"
	//"root/reciver"
	SharedData "root/sharedData"
	"root/transmitter"
	"root/backup"
    "flag"
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
    remoteEventRecived 	:= make(chan config.Update)
    disconnected 		:= make(chan string)

	
	sharedData := SharedData.InitSharedData()
    externalConn := SharedData.InitExternalConn()
	//network.StartPeerNetwork(remoteEventRecived, disconnected, sharedData, externalConn)
    
    
    elevator.FSM_MakeElevator(&elev, elevator_1_ip, config.Num_floors)
    go elevator.FSM_DetectLocalEvents(localEventRecived)
    fmt.Println(cabBackup)
    if isRestart{
        elevator.RestorCabRequests(&elev,cabBackup)
    }
    
    go backup.Start_backup(&elev)

    transmitter.Send_Elevator_data(elevator.GetElevatorData(&elev), externalConn) 
    // RequestsShouldClearImmediately is bugged. Doesnt allow you to call the elevator from the floor it just left

    for {
        select {
        case localEvent := <-localEventRecived:
            elevator.FSM_HandleLocalEvent(&elev, localEvent, sharedData, externalConn)
			elevator.SetAllLights(&elev, sharedData)
            elevator.Send_Elevator_data(elevator.GetElevatorData(&elev), externalConn)

			// Transmitt ? No, because we only transmitt changes. It would not be possible to put it here. 
            // This is because not all events should be transmitted. If a request is handled immediately, because
            // we are at the same floor, we should not transmit that.
			
            // I think assign should be called here? Cause it requires the external data

            // Either expand HandleLocalEvent to include assign and setlights and stuff, or call these separately here.
            // Some control logic too, maybe. We need a more defined function for that.
            // This should also improve Remoteevent handling, as we can use that function there as well.

        case remoteEvent := <-remoteEventRecived:
			elevator.FSM_HandleRemoteEvent(&elev, sharedData, remoteEvent, *externalConn)
            
        case id := <-disconnected:
            externalConn.ConnectedConn[id]=false
			go network.ReconnectPeer(remoteEventRecived, disconnected, id, sharedData, externalConn,&elev)

        }
    }
}

