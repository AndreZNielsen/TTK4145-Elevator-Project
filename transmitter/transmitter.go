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


var sendMu sync.Mutex 

var Disconnected chan<- string

func Start_tcp_call(port string, ip string, id string,disconnected chan<- string,externalConn *sharedData.ExternalConn)net.Conn{
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
		Disconnected = disconnected
	
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

func transmitt_Elevator_data(data config.Elevator_data,id string,externalConn *sharedData.ExternalConn){

	var netErr *net.OpError

	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending
	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
	err := encoder.Encode("elevator_data") // Type ID so the receiver kows what type of data to decode the next packat as 
	if errors.As(err, &netErr) { // check if it is a network-related error
		fmt.Println("Network error:", netErr)
		fmt.Println("Trying to reconnect")
		externalConn.ConnectedConn[id]=false
		Disconnected<-id
		Send_Elevator_data(data,externalConn)
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

func Send_update(update [3]int,externalConn *sharedData.ExternalConn){
	for _, id := range config.RemoteIDs{
		if externalConn.ConnectedConn[id]{
			go transmitt_update(update,id,externalConn)
			
		}
	}
}

func transmitt_update(update [3]int, id string,externalConn *sharedData.ExternalConn){
	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending

	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
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

func Send_alive(externalConn *sharedData.ExternalConn){
	for _, id := range config.RemoteIDs{
		go transmitt_alive(id,externalConn)
	}

}
func transmitt_alive(id string,externalConn *sharedData.ExternalConn){

	var netErr *net.OpError


	for {
		encoder := gob.NewEncoder(externalConn.RemoteElevatorConnections[id])
		
		sendMu.Lock() // Locking before sending
		if externalConn.ConnectedConn[id]{
			err := encoder.Encode("alive")
			if errors.As(err, &netErr) { // check if it is a network-related error
				fmt.Println("Network error:", netErr)
				fmt.Println("Trying to reconnect")
				Disconnected<-id
				externalConn.ConnectedConn[id] = false
				fmt.Println("reconnect reconekted")
				time.Sleep(1*time.Second)
				go Send_alive(externalConn)
				sendMu.Unlock() // Ensure to unlock after sending
				return
			}


			fmt.Println("sent alive")
		}
		
		sendMu.Unlock() // Ensure to unlock after sending
		time.Sleep(time.Second*2)

	}
}

