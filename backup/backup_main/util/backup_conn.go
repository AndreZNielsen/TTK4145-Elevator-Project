package util

import (
	"bufio"
	"os"
	"fmt"
	"encoding/json"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func Msg_reciver(elvatorAlive chan bool) {
    scanner := bufio.NewScanner(os.Stdin) 
    fmt.Println("Starting message receiver...")
    for scanner.Scan() {
        msg := scanner.Text()
        var receivedMsg Message
        err := json.Unmarshal([]byte(msg), &receivedMsg)
        if err != nil {
            fmt.Println("Error decoding  message:", err)
            continue
        }
        fmt.Printf("Received: %+v\n", receivedMsg)
        // Signal that the elvator is alive if the message type is "message"
        if receivedMsg.Type == "message" {
            elvatorAlive <- true
        }
    }
}


// Send messages to elvator 
func Msg_transmitter() {
    message := Message{"message", "message received"}
    jsonData, _ := json.Marshal(message)
    fmt.Fprintln(os.Stdout, string(jsonData))
}
