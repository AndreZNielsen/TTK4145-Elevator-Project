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




func Start_tcp_listen(port string, id string) net.Conn {
    // If there is an existing connection for this id, close it.
    if existingConn := sharedData.RemoteElevatorConnections[id]; existingConn != nil {
        existingConn.Close()
    }

    // Start a new listener locally.
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Error starting listener:", err)
        return nil
    }
    // Accept a new connection.
    conn, err := ln.Accept()
    if err != nil {
        fmt.Println("Error accepting connection:", err)
        ln.Close() // Close listener on error.
        return nil
    }

    ln.Close()

    // Update shared data with the new connection.
    sharedData.Connected_conn[id] = true

    fmt.Println("Connected")
    return conn
}


func Listen_recive(receiver chan<- [3]int,disconnected chan<- string) {
	for _, id := range config.RemoteIDs{
		go Recive(receiver,id,disconnected)
	}
}
func Recive(receiver chan<- [3]int,id string,disconnected chan<- string){
	for {	
		if sharedData.Connected_conn[id]{	
			Decode(receiver,id,disconnected)
		}else{
			return}

	}
}

var data = config.Elevator_data{Behavior: "doorOpen",Floor: 0,Direction: "down",CabRequests: []bool{true, false, false, false}}



func Decode(receiver chan<- [3]int,id string,disconnected chan<- string) {
	decoder := gob.NewDecoder(sharedData.RemoteElevatorConnections[id])

	var typeID string
	err := decoder.Decode(&typeID) // Read type identifier to kono what type of data to decode next
	var netErr *net.OpError
	if errors.As(err, &netErr) { // check if it is a network-related error
		fmt.Println("Network error:", netErr)
		sharedData.Connected_conn[id] = false
		disconnected<-id
		return
	}
	if err != nil {
		fmt.Println("Error decoding type:", err)
		time.Sleep(1*time.Second)
		return
	}


	switch typeID {//chooses what decoder to use based on what type that needs to be decoded 
	case "elevator_data":
		var data config.Elevator_data

		err = decoder.Decode(&data)
		if err != nil {
			fmt.Println("Error decoding Elevator_data:", err)
	
			return
		}
		if data.Floor != -1 && !(data.Floor == 0 && data.Direction == "down") && !(data.Floor == 3 && data.Direction == "up") {//stops the elavator data form crashing the assigner 
		sharedData.ChangeRemoteElevatorData(data,id)
		}
			
		//fmt.Println("Received Elevator_data:", data)
		


	case "int":
		var num [3]int
		err = decoder.Decode(&num)
		if err != nil {
			fmt.Println("Error decoding int:", err)
			return
		}
		receiver<-num
		//fmt.Println("Received int:", num)

	case "alive":
		StartTimer()
		fmt.Println("StartTimer")
	
	default:
		fmt.Println("Unknown type received:", typeID)
	}
}

