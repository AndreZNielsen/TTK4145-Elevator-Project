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
		// Initialize the mutex for a given id if it doesn't already exist.
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
	externalData.ConnectedConn[id]=true
	Disconnected = disconnected
	return conn_lift
}

func SetConn(externalData *sharedData.ExternalData){
	RemoteElevatorConn = externalData.RemoteElevatorConnections
}


func Send_Elevator_data(data config.Elevator_data, externalData *sharedData.ExternalData) {
	for _, id := range config.RemoteIDs{
		if externalData.ConnectedConn[id] {
			go transmitt_Elevator_data(data,id, externalData)
			
		}
	}


}

func transmitt_Elevator_data(data config.Elevator_data,id string, externalData *sharedData.ExternalData){

	var netErr *net.OpError


	sendMu[id].Lock() // Locking before sending
	defer sendMu[id].Unlock() // Ensure to unlock after sending

	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(RemoteElevatorConn[id])
	err := encoder.Encode("elevator_data") // Type ID so the receiver kows what type of data to decode the next packat as 
	if errors.As(err, &netErr) { // check if it is a network-related error
		fmt.Println("Network error:", netErr)
		fmt.Println("Trying to reconnect")
		externalData.ConnectedConn[id]=false
		Disconnected<-id
		Send_Elevator_data(data, externalData)
		fmt.Println("reconnect reconekted")

		time.Sleep(1*time.Second)
		return
	}
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
	time.Sleep(7*time.Millisecond)
	err = encoder.Encode(data) //sendes the Elevator_data
	if err != nil {	
		fmt.Println("Error encoding data:", err)
		return
	}
}

func Send_update(update [3]int, externalData *sharedData.ExternalData){
	for _, id := range config.RemoteIDs{
		if externalData.ConnectedConn[id]{
			go transmitt_update(update,id, externalData)
			
		}
	}
}


func transmitt_update(update [3]int, id string,externalConn *sharedData.ExternalConn){
	sendMu[id].Lock() // Locking before sending
	defer sendMu[id].Unlock() // Ensure to unlock after sending


	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(RemoteElevatorConn[id])
	err := encoder.Encode("int") // Type ID so the receiver kows what type of data to decode the next packat as 
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

func Send_alive(externalData *sharedData.ExternalData){
	for _, id := range config.RemoteIDs{
		go transmitt_alive(id, externalData)
	}

}
func transmitt_alive(id string, externalData *sharedData.ExternalData){
	SetConn(externalData)//Ensure conn is up-to-date
	var netErr *net.OpError


	
	for {
		SetConn(externalData)//Ensure conn is up-to-date
		encoder := gob.NewEncoder(RemoteElevatorConn[id])

		sendMu[id].Lock() // Locking before sending
		if externalConn.ConnectedConn[id]{

			err := encoder.Encode("alive")
			if errors.As(err, &netErr) { // check if it is a network-related error
				fmt.Println("Network error:", netErr)
				fmt.Println("Trying to reconnect")
				Disconnected<-id
				externalData.ConnectedConn[id] = false
				fmt.Println("reconnect reconekted")
				time.Sleep(1*time.Second)

				go Send_alive(externalConn)
				sendMu[id].Unlock() // Ensure to unlock after sending

				return
			}


			fmt.Println("sent alive")
		}
		
		sendMu[id].Unlock() // Ensure to unlock after sending
		time.Sleep(time.Second*2)

	}
}


