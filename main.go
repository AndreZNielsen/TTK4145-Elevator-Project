package main

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/network"
	// "root/reciver"
	SharedData "root/sharedData"
	// "root/transmitter"
    "root/backup"
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
    localEventRecived 	:= make(chan elevator.LocalEvent)
    // aliveTimer 			:= make(chan bool)
    remoteEventRecived 	:= make(chan [3]int)
    disconnected 		:= make(chan string)


	sharedData := SharedData.InitSharedData()
    sharedConn := SharedData.InitExternalConn()
    go network.StartPeerNetwork(remoteEventRecived,disconnected,sharedData,sharedConn)
    var elev elevator.Elevator
    elevator.FSM_MakeElevator(&elev, "localhost:12345", config.Num_floors)
    elevator.Start_if_idle(&elev)
    go elevator.FSM_DetectLocalEvents(localEventRecived)
    go backup.Start_backup()
    // network.Start_network(remoteEventRecived, disconnected, externalData)       // I think we should only pass externalData.RemoteElevatorConnections, if only that is needed!
    // transmitter.Send_Elevator_data(elevator.GetElevatorData(&elev), externalData) 
    // go reciver.AliveTimer(aliveTimer)



    // Buttons only work when door is open, why is this?! This is fixed
    // RequestsShouldClearImmediately is bugged. Doesnt allow you to call the elevator from the floor it just left

    for {
        select {
        case localEvent := <-localEventRecived:
            elevator.FSM_HandleLocalEvent(&elev, localEvent, sharedData)
			elevator.SetAllLights(&elev, sharedData)
			// Transmitt ? No, because we only transmitt changes. It would not be possible to put it here. 
            // This is because not all events should be transmitted. If a request is handled immediately, because
            // we are at the same floor, we should not transmit that.
			
            // I think assign should be called here? Cause it requires the external data

            // Either expand HandleLocalEvent to include assign and setlights and stuff, or call these separately here.
            // Some control logic too, maybe. We need a more defined function for that.
            // This should also improve Remoteevent handling, as we can use that function there as well.

        case remoteEvent := <-remoteEventRecived:
			elevator.FSM_HandleRemoteEvent(&elev, sharedData, remoteEvent)
            fmt.Println("It happend :/")

        case id := <-disconnected:
            go network.ReconnectPeer(remoteEventRecived, disconnected, id, sharedData,sharedConn)

        // case <-aliveTimer:

        }
    }
}