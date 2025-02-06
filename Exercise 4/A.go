// other.go (TCP Server)
package main

import (
	"fmt"
	"net"
	"os"
	"time"
	//"time"
)
var time_since time.Time

func main() {
	// Listen on TCP port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Listening on port 8080...")
	time_since = time.Now()
	go is_alive()
	for {
		// Accept an incoming connection
		conn, err := ln.Accept()
		if err != nil {
		fmt.Println("Error accepting connection:", err)
		}
		// Handle the connection in a new goroutine
		go handleRequest(conn)

	}
}

func handleRequest(conn net.Conn) {
	// Close the connection when the function finishes
	defer conn.Close()

	// Send a welcome message to the client (main.go)
	conn.Write([]byte("Hello from other.go!\n"))

	// Read data sent by the client
	buffer := make([]byte, 1024)
	for {
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return
	}
	fmt.Println("Received from main.go:", string(buffer))
	time_since = time.Now()
	
}

}
func is_alive(){
	for {
	fmt.Println(time.Since(time_since) > time.Second *10)
	if time.Since(time_since) > time.Second *10 {
		fmt.Println("B not alive")
		break
	}
	time.Sleep(time.Second*2)
	}
	os.Exit(0)
}