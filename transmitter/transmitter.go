package transmitter

import (
	"fmt"
	"root/sharedData"
	"root/config"
	"net"
	"sync"
	"time"
	"errors"
	"encoding/json"
    "root/customStructs"

)

type Message struct {
    TypeID string      `json:"typeID"`
    Data   interface{} `json:"data"`
}


var sendMu = make(map[string]*sync.Mutex)
var Disconnected chan<- string
var netErr *net.OpError


func InitMutex(){
	for _,id := range(config.RemoteIDs){
		// Initialize the mutex if it doesn't already exist
		if _, exists := sendMu[id]; !exists {
			sendMu[id] = &sync.Mutex{}
		}
	}
}

func InitDiscEventChan(disconnected chan<- string){
	Disconnected = disconnected
}


func Start_tcp_call(port string, ip string, id string,externalConn *sharedData.ExternalConn)net.Conn{
	for{
		if existingConn := externalConn.RemoteElevatorConnections[id]; existingConn != nil {// Close the previous listener if it's still open.
			existingConn.Close()
		}
		
		conn_lift, err := net.Dial("tcp", ip+":"+port)//connects to the other elevatoe
		
		if err != nil {
			fmt.Println("Error connecting to pc:", ip, err)
			time.Sleep(5*time.Second)
			continue //trys again
		}	
		externalConn.ConnectedConn[id]=true

		return conn_lift
	}
}



func Send_Elevator_data(data customStructs.Elevator_data,externalConn *sharedData.ExternalConn) {
	for _, id := range config.RemoteIDs{//sends the elevator data to all the other elevators
		if externalConn.ConnectedConn[id] {
			go transmitt_Elevator_data(data,id,externalConn)
		}
	}
}

func transmitt_Elevator_data(data customStructs.Elevator_data, id string, externalConn *sharedData.ExternalConn) {
    for {
        sendMu[id].Lock()  // Locking before sending
        if externalConn.ConnectedConn[id] {

            message := Message{
                TypeID: "elevator_data", 
                Data:   data,             
            }
            
            encoder := json.NewEncoder(externalConn.RemoteElevatorConnections[id])

            for i := 0; i < 10; i++ { 
                err := encoder.Encode(message) // Sends the message as a JSON package
                if err != nil {

                    if errors.As(err, &netErr) {
                        fmt.Println("Network error while encoding update:", netErr)
                        Disconnected <- id
                    } else {
                        fmt.Println("Error encoding data:", err)                       
                    }

                    sendMu[id].Unlock()
                    return
                }
            }
            sendMu[id].Unlock()
            return
        }
        sendMu[id].Unlock()
        time.Sleep(1 * time.Second)
    }
}


func Send_update(update customStructs.Update,externalConn *sharedData.ExternalConn){
	for _, id := range config.RemoteIDs{//sends the update to all the other elevators
		if externalConn.ConnectedConn[id]{
			go transmitt_update(update,id,externalConn)
		}
	}
}

func transmitt_update(update customStructs.Update, id string, externalConn *sharedData.ExternalConn) {
    sendMu[id].Lock() // Lock for this specific id to ensure only one thread sends at a time
    defer sendMu[id].Unlock() // Ensure the mutex is unlocked after sending

	message:= Message{
        TypeID: "Update",   // Type ID to identify the message as an update
        Data:   update,     // contains the update data
    }

    encoder := json.NewEncoder(externalConn.RemoteElevatorConnections[id])

    for i := 0; i < 10; i++{

        err := encoder.Encode(message) // Sends the message as a JSON package
        if  err != nil {

            if errors.As(err, &netErr) {
                fmt.Println("Network error while encoding update:", netErr)
                Disconnected <- id
            } else {
                fmt.Println("Error encoding update:", err)
   
            }
            return
        } 
        
    }
}

func Send_alive(externalConn *sharedData.ExternalConn){
	for _, id := range config.RemoteIDs{
		go transmitt_alive(id,externalConn)
	}
}

func transmitt_alive(id string, externalConn *sharedData.ExternalConn) {
    for {
        sendMu[id].Lock() 
        if externalConn.ConnectedConn[id] {
            
            message := Message{
                TypeID: "alive", 
                Data:   "alive", 
            }

            encoder := json.NewEncoder(externalConn.RemoteElevatorConnections[id])
            err := encoder.Encode(message) 

            if err != nil {
                if errors.As(err, &netErr) { 
                    fmt.Println("Network error while encoding alive:", netErr)

                    Disconnected <- id
                } else {
                    fmt.Println("Error encoding alive message:", err)
                }
                sendMu[id].Unlock()
                time.Sleep(1 * time.Second)
                continue
            }
        }
        sendMu[id].Unlock() 
        time.Sleep(1 * time.Second) //this can be adjusted to lower risk of case: disconnect because of packetloss
    }
}



func MergeHallRequests(externalConn *sharedData.ExternalConn, hallRequests customStructs.HallRequests, id string) {
    for {
        if !externalConn.ConnectedConn[id]{
            time.Sleep(1*time.Second)
            continue
        }
        sendMu[id].Lock() 
        defer sendMu[id].Unlock() 

        message := Message{
            TypeID: "mergeHallRequests",
            Data:   hallRequests,   
        }

        encoder := json.NewEncoder(externalConn.RemoteElevatorConnections[id])
        err := encoder.Encode(message) 
        if err != nil {

            if errors.As(err, &netErr) {
                fmt.Println("Network error while encoding mergeHallRequests:", netErr)
                Disconnected <- id
            } else {
                fmt.Println("Error encoding mergeHallRequests:", err)
            }
            return
        }
        return
    }
}

func Send_Hall_Requests(externalConn *sharedData.ExternalConn, hallRequests customStructs.HallRequests) {

    message := Message{
        TypeID: "HallRequests",
        Data:   hallRequests,  
    }

    for _, id := range config.RemoteIDs{
        sendMu[id].Lock() 
		if externalConn.ConnectedConn[id]{
			encoder := json.NewEncoder(externalConn.RemoteElevatorConnections[id])
            err := encoder.Encode(message) 
            if err != nil {
                if errors.As(err, &netErr) {
                    fmt.Println("Network error while encoding update:", netErr)
                    Disconnected <- id
                } else {
                    fmt.Println("Error encoding HallRequests:", err)
                }
                sendMu[id].Unlock() 
                continue
            }
		}
        sendMu[id].Unlock() 
	}    


}
