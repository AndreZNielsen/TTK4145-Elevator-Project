package util

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
var Conn net.Conn
func Msg_transmitter() {
	var err error
    Conn, err = net.Dial("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}

	encoder := json.NewEncoder(Conn)

	for {
		msg := Message{"message", "message received"}
		err := encoder.Encode(msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
		fmt.Println("Sent heartbeat")
		time.Sleep(5 * time.Second)
	}
}


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

func HandleConnection(alive chan bool) {
	decoder := json.NewDecoder(Conn)

	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			return
		}

		fmt.Printf("Received: %+v\n", msg)

		if msg.Type == "message" {
			alive <- true
		}
	}
}