package utility

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	//"time"
	"root/SharedData"
)

var lis_lift1 net.Conn
//var lis_lift2 net.Conn
type Elevator_data struct {
	Behavior    string 
	Floor       int
	Direction   string 
	CabRequests []bool 
	HallRequests [][2]bool        
}

func Start_tcp_listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listen:", err)
	}
	lis_lift1, err = ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}

}

func Listen_recive(receiver chan<- bool) {
	for {
		Decode(receiver)
	}
}

func Decode(receiver chan<- bool) {
	decoder := gob.NewDecoder(lis_lift1)

	var typeID string
	err := decoder.Decode(&typeID) // Read type identifier
	if err != nil {
		fmt.Println("Error decoding type:", err)
		time.Sleep(1*time.Second)
		return
	}

	switch typeID {
	case "elevator_data":
		var data Elevator_data

		err = decoder.Decode(&data)
		if err != nil {
			fmt.Println("Error decoding Elevator_data:", err)
	
			return
		}
		fmt.Println("Received Elevator_data:", data)
		receiver <- true

	case "int":
		var num [3]int
		err = decoder.Decode(&num)
		if err != nil {
			fmt.Println("Error decoding int:", err)
			return
		}
		fmt.Println("Received int:", num)
		sharedData.UpdatesharedHallRequests(num)
		receiver <- true

	default:
		fmt.Println("Unknown type received:", typeID)
	}
}
