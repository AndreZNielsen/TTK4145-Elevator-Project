package transmitter

import (
	"encoding/gob"
	"fmt"
	sharedData "root/SharedData"
	"net"
	"sync"
	"time"
	"errors"
)
var conn_lift net.Conn
//var conn_lift2 net.Conn

var sendMu sync.Mutex 

var RemoteElevatorConn =  make(map[string]net.Conn)

func Start_tcp_call(port string, ip string, id string)net.Conn{
	
	var err error
	if conn_lift != nil {	// Close the previous listener if it's still open.
	conn_lift.Close()
	}
	conn_lift, err = net.Dial("tcp", ip+":"+port)//connects to the other elevatoe
	
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		time.Sleep(5*time.Second)
		conn_lift = Start_tcp_call(port, ip,id)//trys again
		return conn_lift
	}
	sharedData.Connected_conn[id]=true
	return conn_lift
}

func SetConn(){
	RemoteElevatorConn = sharedData.RemoteElevatorConnections
}


func Send_Elevator_data(data sharedData.Elevator_data) {
	for _, id := range sharedData.GetRemoteIDs(){
		if sharedData.Connected_conn[id] {
			go transmitt_Elevator_data(data,id)
			
		}
	}


}

func transmitt_Elevator_data(data sharedData.Elevator_data,id string){

	var netErr *net.OpError

	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending
	SetConn()//Ensure conn is up-to-date
	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(RemoteElevatorConn[id])
	err := encoder.Encode("elevator_data") // Type ID so the receiver kows what type of data to decode the next packat as 
	if errors.As(err, &netErr) { // check if it is a network-related error
		fmt.Println("Network error:", netErr)
		fmt.Println("Trying to reconnect")
		sharedData.Connected_conn[id]=false
		sharedData.Disconnected<-id
		Send_Elevator_data(data)
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

func Send_update(update [3]int){
	for _, id := range sharedData.GetRemoteIDs(){
		if sharedData.Connected_conn[id]{
			go transmitt_update(update,id)
			
		}
	}
}

func transmitt_update(update [3]int, id string){
	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending
	SetConn()//Ensure conn is up-to-date

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

func Send_alive(){
	for _, id := range sharedData.GetRemoteIDs(){
		go transmitt_alive(id)
	}

}
func transmitt_alive(id string){
	SetConn()//Ensure conn is up-to-date
	var netErr *net.OpError


	
	for {
		SetConn()//Ensure conn is up-to-date
		encoder := gob.NewEncoder(RemoteElevatorConn[id])
		
		sendMu.Lock() // Locking before sending
		if sharedData.Connected_conn[id]{
			err := encoder.Encode("alive")
			if errors.As(err, &netErr) { // check if it is a network-related error
				fmt.Println("Network error:", netErr)
				fmt.Println("Trying to reconnect")
				sharedData.Disconnected<-id
				sharedData.Connected_conn[id] = false
				fmt.Println("reconnect reconekted")
				time.Sleep(1*time.Second)
				go Send_alive()
				sendMu.Unlock() // Ensure to unlock after sending
				return
			}


			fmt.Println("sent alive")
		}
		
		sendMu.Unlock() // Ensure to unlock after sending
		time.Sleep(time.Second*2)

	}
}


