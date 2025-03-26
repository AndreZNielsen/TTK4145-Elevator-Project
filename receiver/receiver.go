package receiver

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"root/sharedData"
	"root/config"
	"root/customStructs"


	"errors"

    "bufio"
)

var netErr *net.OpError




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


func Listen_recive(receiver chan<- customStructs.RemoteEvent,
	disconnected chan<- string,
	externalData *sharedData.SharedData,
	externalConn *sharedData.ExternalConn,
	aliveRecievd chan<- string,

	requestHallRequests chan<- customStructs.HallRequests) {

	for _, id := range config.RemoteIDs{
		go Recive(receiver,id,disconnected,externalData,externalConn,aliveRecievd,requestHallRequests)
	}
}


func Recive(receiver chan<- customStructs.RemoteEvent,
    id string,
    disconnected chan<- string,
    externalData *sharedData.SharedData,
    externalConn *sharedData.ExternalConn,
    aliveRecievd chan<- string,


    requestHallRequests chan<- customStructs.HallRequests) {


    scann := bufio.NewScanner(externalConn.RemoteElevatorConnections[id])
    for scann.Scan(){
        if externalConn.ConnectedConn[id] {
            //decoder := json.NewDecoder(externalConn.RemoteElevatorConnections[id])

            var message struct {
                TypeID string          `json:"typeID"` 
                Data   json.RawMessage `json:"data"`   
            }

            
            //err := decoder.Decode(&message)
            err := json.Unmarshal(scann.Bytes(),&message)



            if err != nil {
                if errors.As(err, &netErr) { 
                    fmt.Println("Network error while encoding alive:", netErr)
                    disconnected <- id
                } else {
                    fmt.Println("Error decoding message:", err)
                }


                time.Sleep(1 * time.Second)
                continue
            }

          
            switch message.TypeID {
            case "elevator_data":
                var data customStructs.Elevator_data
                err := json.Unmarshal(message.Data, &data) 
                if err != nil {
                    fmt.Println("Error decoding Elevator_data:", err)
                    return
                }

                event := customStructs.RemoteEvent{
                    EventType:    "elevatorData",
                    Id:           id,
                    ElevatorData: data,
                }
                receiver <- event

            case "Update":
                var update customStructs.Update
                err := json.Unmarshal(message.Data, &update) 
                if err != nil {
                    fmt.Println("Error decoding Update:", err)
                    return
                }

                event := customStructs.RemoteEvent{
                    EventType: "update",
                    Update:    update,
                }
                receiver <- event

            case "alive":
                aliveRecievd <- id

            case "RequestHallRequests":


                 var hallRequests customStructs.HallRequests
                err := json.Unmarshal(message.Data, &hallRequests) 
                if err != nil {
                    fmt.Println("Error decoding Update:", err)
                    return
                }
                requestHallRequests <- hallRequests


            case "HallRequests":
                var hallRequests customStructs.HallRequests
                err := json.Unmarshal(message.Data, &hallRequests) 
                if err != nil {
                    fmt.Println("Error decoding HallRequests:", err)
                    return
                }

                event := customStructs.RemoteEvent{
                    EventType:  "hallRequests",
                    HallRequests: hallRequests,
                }
                receiver <- event

            default:
                fmt.Println("Unknown type received:", message.TypeID)
            }
            aliveRecievd <- id  //every message restarts alivetimer
        } else {
            return
        }
    }
}





