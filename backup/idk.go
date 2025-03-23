package backup

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"root/elevator"
	"runtime"
	"time"
)

type Message struct {
	Type    string      `json:"type"`
	Content []bool `json:"content"`
}

func Start_backup(elev *elevator.Elevator) {
	fmt.Println("Starting backup in 5 sec...")
	time.Sleep(5 * time.Second)

	
		startBackupProcess()
		time.Sleep(4 * time.Second)

		conn, err := net.Dial("tcp", "localhost:5000")
		if err != nil {
			fmt.Println("Failed to connect to backup server:", err)
			time.Sleep(5 * time.Second)
			
		}
		defer conn.Close()

		encoder := json.NewEncoder(conn)

		// Send heartbeats
		for {
			msg := Message{"message", elevator.GetCabRequests(elev.Requests)}
			err := encoder.Encode(msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
				break
			}
			fmt.Println("Sent backup heartbeat")
			time.Sleep(5 * time.Second)
		}
	
}

func startBackupProcess() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", "cd ./backup/backup_main && go run backup_main.go")				
	case "windows":
		psCommand := "Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run backup_main.go\""
	
		// Start PowerShell and execute the command
		cmd = exec.Command("powershell.exe", "-Command", psCommand)
		cmd.Dir = "./backup/backup_main"

	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting backup process:", err)
	}
}
