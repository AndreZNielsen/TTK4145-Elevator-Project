package backup

import (
	"encoding/json"
	"fmt"
	"root/elevator"
	"time"
	"net"
)
var Conn net.Conn


func StartConn(){
	if Conn !=nil{
		Conn.Close()
	}

	var err error
	Conn, err = net.Dial("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Failed to connect to backup server:", err)
		time.Sleep(5 * time.Second)
		
	}
}
func HandleConnection(backupAlive chan bool) {
	decoder := json.NewDecoder(Conn)

	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			return
		}

		//fmt.Printf("Received: %+v\n", msg.Content)

		if msg.Type == "message" {
			backupAlive<-true
		}
	}
}

func SendCabHartBeat(elev *elevator.Elevator){
	encoder := json.NewEncoder(Conn)

	// Send heartbeats
	for {
		msg := Message{"message", elevator.GetCabRequests(elev.Requests)}
		err := encoder.Encode(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
		//fmt.Println("Sent backup heartbeat")
		time.Sleep(2 * time.Second)
	}
	}

