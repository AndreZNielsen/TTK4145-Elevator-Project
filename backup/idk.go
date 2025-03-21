package backup

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

func Start_backup() {
	fmt.Println("Starting backup in 10 sec...")
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
			msg := Message{"message", "backup running"}
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
		cmd = exec.Command("xterm", "-e", "bash -c 'cd ./backup/backup_main && go run backup_main.go; exec bash'")
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
