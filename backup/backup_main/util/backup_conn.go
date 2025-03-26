package util

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Type    string      `json:"type"`
	Content []bool `json:"content"`
}
var Conn net.Conn



func StartTCPLis() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started, waiting for connections...")
    Conn, err = listener.Accept()
    
    if err != nil {
        fmt.Println("Connection error:", err)
    
    }
}

func HandleConnection(alive chan []bool) {
	decoder := json.NewDecoder(Conn)

	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			return
		}

		fmt.Printf("Received: %+v\n", msg.Content)

		if msg.Type == "message" {// as long as it recives messages it will restart the backupAlive timer and store cab requests
			alive <- msg.Content
		}
	}
}

func SendHartBeat() {// sends a heartbeat to the elevator to show that the backup is alive
	encoder := json.NewEncoder(Conn)

	for {
		msg := Message{"message", make([]bool, 2)}
		err := encoder.Encode(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
		fmt.Println("Sent heartbeat")
		time.Sleep(3 * time.Second)
	}
}