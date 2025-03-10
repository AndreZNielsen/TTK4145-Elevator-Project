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
var conn_lift1 net.Conn
//var conn_lift2 net.Conn

var sendMu sync.Mutex 
var port1 string
var ip1 string
var transmitter__initialized_1 chan<- bool

func Start_tcp_call(port string, ip string, transmitter__initialized chan<- bool){
	var err error
	port1 = port
	ip1 = ip
	transmitter__initialized_1 = transmitter__initialized
	if conn_lift1 != nil {	// Close the previous listener if it's still open.
	conn_lift1.Close()
	}
	conn_lift1, err = net.Dial("tcp", ip+":"+port)//connects to the other elevatoe
	
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		time.Sleep(5*time.Second)
		Start_tcp_call(port, ip,transmitter__initialized)//trys again
		return
	}
	transmitter__initialized <- true
}


func Send_Elevator_data(data sharedData.Elevator_data) {
	var netErr *net.OpError

	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending
	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(conn_lift1)
	err := encoder.Encode("elevator_data") // Type ID so the receiver kows what type of data to decode the next packat as 
	if errors.As(err, &netErr) { // check if it is a network-related error
		fmt.Println("Network error:", netErr)
		fmt.Println("Trying to reconnect")
		Start_tcp_call(port1,ip1,nil)
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


	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending
	time.Sleep(7*time.Millisecond)
	encoder := gob.NewEncoder(conn_lift1)
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
	encoder := gob.NewEncoder(conn_lift1)
	var netErr *net.OpError
	
	for {
		sendMu.Lock() // Locking before sending
		defer sendMu.Unlock()
		err := encoder.Encode("alive")
		if errors.As(err, &netErr) { // check if it is a network-related error
			fmt.Println("Network error:", netErr)
			fmt.Println("Trying to reconnect")
			Start_tcp_call(port1,ip1,nil)
			fmt.Println("reconnect reconekted")
			time.Sleep(1*time.Second)
			go Send_alive()
			return
		}
		if err != nil {
			fmt.Println("Encoding error:", err)
			return
		}
		sendMu.Unlock() // Ensure to unlock after sending
		fmt.Println("sent alive")
		time.Sleep(time.Second)
	}
}






