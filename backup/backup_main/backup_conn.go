package main

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

func Msg_reciver(parentAlive chan bool){
	scanner := bufio.NewScanner(os.Stdout)
	for scanner.Scan() {
		msg := scanner.Text()
		var receivedMsg Message
		err := json.Unmarshal([]byte(msg), &receivedMsg)
		if err != nil {
			fmt.Println("Error decoding parent message:", err)
			continue
		}

		fmt.Printf("Child received: %+v\n", receivedMsg)

		// Check if the parent is alive
		if receivedMsg.Type == "message" {
			parentAlive <- true
		}
	}
}	

func Msg_transmitter(){
		message := Message{"message", "message recived"}
		jsonData, _ := json.Marshal(message)
		fmt.Fprintln(os.Stdin, string(jsonData))
}

