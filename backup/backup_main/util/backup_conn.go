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

func Msg_reciver(parentAlive chan bool) {
    scanner := bufio.NewScanner(os.Stdin) // Change to os.Stdin
    fmt.Println("Starting message receiver...")
    for scanner.Scan() {
        msg := scanner.Text()
        var receivedMsg Message
        err := json.Unmarshal([]byte(msg), &receivedMsg)
        if err != nil {
            fmt.Println("Error decoding parent message:", err)
            continue
        }
        fmt.Printf("Child received: %+v\n", receivedMsg)
        // Signal that the parent is alive if the message type is "message"
        if receivedMsg.Type == "message" {
            parentAlive <- true
        }
    }
}


// Send messages to parent via child's standard output
func Msg_transmitter() {
    message := Message{"message", "message received"}
    jsonData, _ := json.Marshal(message)
    fmt.Fprintln(os.Stdout, string(jsonData)) // Write to os.Stdout
}
