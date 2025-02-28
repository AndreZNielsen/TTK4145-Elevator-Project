package utility

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"root/SharedData"
)

var mu sync.Mutex
var a sync.Mutex   // Mutex to protect the sending of data to the receiver channel
var lis_lift1 net.Conn



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
		mu.Lock()
		Decode(receiver)
		mu.Unlock()
	}
}

func Decode(receiver chan<- bool) {
	decoder := gob.NewDecoder(lis_lift1)

	var typeID string
	err := decoder.Decode(&typeID) // Read type identifier
	if err != nil {
		fmt.Println("Error decoding type:", err)
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
		
		// Protecting the sending operation with a mutex to ensure that only one goroutine can send at a time
		a.Lock()
		receiver <- true
		a.Unlock()

	case "int":
		var num [3]int
		err = decoder.Decode(&num)
		if err != nil {
			fmt.Println("Error decoding int:", err)
			return
		}
		fmt.Println("Received int:", num)
		
		// Protecting the sending operation with a mutex to ensure that only one goroutine can send at a time
		a.Lock()
		sharedData.UpdatesharedHallRequests(num)
		receiver <- true
		a.Unlock()

	default:
		fmt.Println("Unknown type received:", typeID)
	}
}
