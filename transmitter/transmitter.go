package transmitter

import (
	"encoding/gob"
	"fmt"
	"root/sharedData"
	"root/config"
	"net"
	"sync"
	"time"
	"errors"
)


//var sendMu sync.Mutex 
var sendMu = make(map[string]*sync.Mutex)
var Disconnected chan<- string
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



func Send_Elevator_data(data config.Elevator_data,externalConn *sharedData.ExternalConn) {
	for _, id := range config.RemoteIDs{
		if externalConn.ConnectedConn[id] {
			go transmitt_Elevator_data(data,id,externalConn)
			
		}
	}
}

var prevData = config.Elevator_data{Behavior: "doorOpen",Floor: 1,Direction: "down",CabRequests: make([]bool, config.Num_floors)}

func transmitt_Elevator_data(data config.Elevator_data,id string,externalConn *sharedData.ExternalConn){

	var netErr *net.OpError

	if data.Floor == -1{
		data = prevData
	}else{
		prevData = data
	}
	for {
		sendMu[id].Lock() 	// Locking before sending
		if externalConn.ConnectedConn[id]{
			time.Sleep(7*time.Millisecond)
			encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
			err := encoder.Encode("elevator_data") // Type ID so the receiver kows what type of data to decode the next packat as 
			if errors.As(err, &netErr) { // check if it is a network-related error
				fmt.Println("Network error:", netErr)
				Disconnected<-id
				sendMu[id].Unlock() 
				continue
			}
			if err != nil {
				fmt.Println("Encoding error:", err)
				sendMu[id].Unlock() 
				return
			}
			time.Sleep(7*time.Millisecond)
			err = encoder.Encode(data) //sendes the Elevator_data
			if err != nil {	
				fmt.Println("Error encoding data:", err)
				sendMu[id].Unlock() 
				return
			}
			sendMu[id].Unlock()
			return
		}
		sendMu[id].Unlock()
		time.Sleep(1*time.Second)
	}
}

func Send_update(update config.Update,externalConn *sharedData.ExternalConn){
	for _, id := range config.RemoteIDs{
		if externalConn.ConnectedConn[id]{
			go transmitt_update(update,id,externalConn)
			
		}
	}
}

func transmitt_update(update config.Update, id string,externalConn *sharedData.ExternalConn){
	sendMu[id].Lock() // Locking before sending
	defer sendMu[id].Unlock() // Ensure to unlock after sending

	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
	err := encoder.Encode("Update") // Type ID so the receiver kows what type of data to decode the next packat as 
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
	time.Sleep(7*time.Millisecond)
	err = encoder.Encode(update) //sendes the update
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

}

func Send_alive(externalConn *sharedData.ExternalConn){
	for _, id := range config.RemoteIDs{
		go transmitt_alive(id,externalConn)
	}

}
func transmitt_alive(id string,externalConn *sharedData.ExternalConn){

	var netErr *net.OpError


	for {
		encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
		
		sendMu[id].Lock() // Locking before sending
		if externalConn.ConnectedConn[id]{
			err := encoder.Encode("alive")
			if errors.As(err, &netErr) { // check if it is a network-related error
				fmt.Println("Network error:", netErr)
				Disconnected<-id
				sendMu[id].Unlock() 
				time.Sleep(1*time.Second)
				
				continue
			}


			//fmt.Println("sent alive")
		}
		
		sendMu[id].Unlock() // Ensure to unlock after sending
		time.Sleep(time.Second*2)

	}
}

func RequestHallRequests(externalConn *sharedData.ExternalConn, id string){
	sendMu[id].Lock() // Locking before sending
	defer sendMu[id].Unlock() // Ensure to unlock after sending

	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
	err := encoder.Encode("RequestHallRequests") // Type ID so the receiver kows what type of data to decode the next packat as 
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
}

func Send_Hall_Requests(id string,externalConn *sharedData.ExternalConn,sharedData *sharedData.SharedData){
	sendMu[id].Lock() // Locking before sending
	defer sendMu[id].Unlock() // Ensure to unlock after sending

	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
	err := encoder.Encode("HallRequests") // Type ID so the receiver kows what type of data to decode the next packat as 
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
	time.Sleep(7*time.Millisecond)
	err = encoder.Encode(sharedData.HallRequests) //sendes the update
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

}