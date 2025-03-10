package reciver

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
	"root/SharedData"
	"errors"

	
)
var ln net.Listener
var lis_lift1 net.Conn
//var lis_lift2 net.Conn
var RemoteElevatorConn =  make(map[string]net.Conn)

var port1 string 
var data = sharedData.Elevator_data{Behavior: "doorOpen",Floor: 0,Direction: "down",CabRequests: []bool{true, false, false, false}}
var Connected bool
func Start_tcp_listen(port string) {
	port1 = port

	if ln != nil {	// Close the previous listener if it's still open.
		ln.Close()
	}
	var err error
	ln, err = net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listen:", err)
		return
	}
	lis_lift1, err = ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	Connected = true
	fmt.Println("Connected")

}
func Start_tcp_listen2(port string) net.Conn{
	port1 = port

	if ln != nil {	// Close the previous listener if it's still open.
		ln.Close()
	}
	var err error
	ln, err = net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listen:", err)
		return nil
	}
	lis_lift1, err = ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return nil
	}
	Connected = true
	fmt.Println("Connected")
	return lis_lift1


}

func SetConn(){
	RemoteElevatorConn = sharedData.RemoteElevatorConnections
}

func Listen_recive(receiver chan<- [3]int) {
	for _, id := range sharedData.GetRemoteIDs(){
		fmt.Println("yoooooooooooo")
		go recive(receiver,id)
	}
}
func recive(receiver chan<- [3]int,id string){
	for {	
			Decode(receiver,id)

	}
}
func Decode(receiver chan<- [3]int,id string) {
	SetConn()//Ensure conn is up-to-date
	decoder := gob.NewDecoder(RemoteElevatorConn[id])

	var typeID string
	err := decoder.Decode(&typeID) // Read type identifier to kono what type of data to decode next
	var netErr *net.OpError
	if errors.As(err, &netErr) { // check if it is a network-related error
		fmt.Println("Network error:", netErr)
		go Start_tcp_listen(port1)
		Connected = false
		time.Sleep(1*time.Second)
		return
	}
	if err != nil {
		fmt.Println("Error decoding type:", err)
		time.Sleep(1*time.Second)
		return
	}


	switch typeID {//chooses what decoder to use based on what type that needs to be decoded 
	case "elevator_data":
		var data sharedData.Elevator_data

		err = decoder.Decode(&data)
		if err != nil {
			fmt.Println("Error decoding Elevator_data:", err)
	
			return
		}
		if data.Floor != -1 && !(data.Floor == 0 && data.Direction == "down") && !(data.Floor == 3 && data.Direction == "up") {//stops the elavator data form crashing the assigner 
		sharedData.ChangeRemoteElevatorData(data,sharedData.GetRemoteIDs()[0])
		}
			
		//fmt.Println("Received Elevator_data:", data)
		


	case "int":
		var num [3]int
		err = decoder.Decode(&num)
		if err != nil {
			fmt.Println("Error decoding int:", err)
			return
		}
		//fmt.Println("Received int:", num)
		receiver <- num //sends signal to main that hall requests have been updated and that the lights need to be updated

	case "alive":
		StartTimer()
		fmt.Println("StartTimer")
	
	default:
		fmt.Println("Unknown type received:", typeID)
	}
}



func Connection_lost(){
	Connected = false
}