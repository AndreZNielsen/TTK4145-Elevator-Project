package reciver

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
	"root/sharedData"
	"root/config"
	"errors"

	
)




func Start_tcp_listen(port string, id string,externalConn *sharedData.ExternalConn) net.Conn {
    // If there is an existing connection for this id, close it.
    if existingConn := externalConn.RemoteElevatorConnections[id]; existingConn != nil {
        existingConn.Close()
    }

    // Start a new listener locally.
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Error starting listener:", err)
        return nil
    }
    defer ln.Close()
    // Accept a new connection.
    conn, err := ln.Accept()
    if err != nil {
        fmt.Println("Error accepting connection:", err)
        ln.Close() // Close listener on error.
        return nil
    }


    // Update shared data with the new connection.
    externalConn.ConnectedConn[id] = true

    fmt.Println("Connected")
    return conn
}


func Listen_recive(receiver chan<- config.RemoteEvent,
	disconnected chan<- string,
	externalData *sharedData.SharedData,
	externalConn *sharedData.ExternalConn,
	aliveRecievd chan<- string,
	requestHallRequests chan<- string) {
	for _, id := range config.RemoteIDs{
		go Recive(receiver,id,disconnected,externalData,externalConn,aliveRecievd,requestHallRequests)
	}
}

var data = config.Elevator_data{Behavior: "doorOpen",Floor: 0,Direction: "down",CabRequests: []bool{true, false, false, false}}

func Recive(receiver chan<- config.RemoteEvent,
	id string,disconnected chan<- string,
	externalData *sharedData.SharedData,
	externalConn *sharedData.ExternalConn,
	aliveRecievd chan<- string,
	requestHallRequests chan<- string	){
	for {	
		if externalConn.ConnectedConn[id]{	
			decoder := gob.NewDecoder(externalConn.RemoteElevatorConnections[id])

			var typeID string
			err := decoder.Decode(&typeID) // Read type identifier to kono what type of data to decode next
			var netErr *net.OpError
			if errors.As(err, &netErr) { // check if it is a network-related error
				fmt.Println("Network error:", netErr)
				if externalConn.ConnectedConn[id]{ 
				
					disconnected<-id
				}
				return
			}
			if err != nil {
				fmt.Println("Error decoding type:", err)
				time.Sleep(1*time.Second)
				continue
			}
		
		
			switch typeID {//chooses what decoder to use based on what type that needs to be decoded 
			case "elevator_data":
				var data config.Elevator_data
		
				err = decoder.Decode(&data)
				if err != nil {
					fmt.Println("Error decoding Elevator_data:", err)
			
					return
				}
				// if data.Floor != -1 && !(data.Floor == 0 && data.Direction == "down") && !(data.Floor == 3 && data.Direction == "up") {//stops the elavator data form crashing the assigner 
				// 	externalData.RemoteElevatorData[id]=data
				// }
				// receiver<-config.Update{Floor: 0,ButtonType: 2,Value: false}//dummy update to trigger remote event
				
				event := config.RemoteEvent{
					EventType: "elevatorData",
					Id: id,
					ElevatorData: data,
				}

				receiver <- event

				//fmt.Println("Received Elevator_data:", data)
				
		
		
			case "Update":
				var Update config.Update
				err = decoder.Decode(&Update)
				if err != nil {
					fmt.Println("Error decoding int:", err)
					return
				}

				//receiver<-Update

				event := config.RemoteEvent{EventType: "update"}

				receiver<-event
		
			case "alive":
				aliveRecievd<-id
			
			case "RequestHallRequests":
				requestHallRequests<-id

			case "HallRequests":
				var hallRequests [][2]bool
		
				err = decoder.Decode(&hallRequests)
				if err != nil {
					fmt.Println("Error decoding Elevator_data:", err)
			
					return
				}
				//externalData.HallRequests=hallRequests
				//receiver<-config.Update{Floor: 0,ButtonType: 2,Value: false}//dummy update to trigger remote event
			
				event := config.RemoteEvent{
					EventType: "hallRequests",
					HallRequests: hallRequests}
			
				receiver <- event

			default:
				fmt.Println("Unknown type received:", typeID)
			}
		}else{
			return}

	}
}




