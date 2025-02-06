package main

import (
	"fmt"
	"os"
	"os/exec"
	"net"
	"time"
	"bufio"
)

func main() {
	// Start the PowerShell window that runs other.go
	psCommand := "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'go run A.go'"

	// Start PowerShell and execute the command
	cmd := exec.Command("powershell.exe", "-Command", psCommand)

	// Start the PowerShell window in parallel
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting PowerShell:", err)
		return
	}

	// Give the PowerShell some time to start up and listen on the TCP port
	time.Sleep(time.Second*5)
	// Connect to the TCP server (other.go)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()
	time.Sleep(time.Second*5)
	// Send a message to the server (other.go)
	message := "Hello from main.go!\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	// Use goroutine to send the "I am alive" message repeatedly
	alive(conn)

	// Read the response from the server (other.go)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("Received from other.go:", scanner.Text())
		if scanner.Text() == "other.go has finished." {
			break
		}
	}

	// Wait for the PowerShell window (other.go) to finish its task
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for PowerShell:", err)
	}

	// Continue with main program work after communication is done
	fmt.Println("Main program has finished.")
}

func alive(conn net.Conn) {
	var err error
	for i := 0; i < 5; i++ {
		message := "I am alive\n"
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		time.Sleep(time.Second)
	}
}
