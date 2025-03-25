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
var Conn net.Conn
var BackupDead = make(chan bool)

func Start_backup(elev *elevator.Elevator) {
	var err error

	fmt.Println("Starting backup in 5 sec...")
	time.Sleep(5 * time.Second)
	
	for {
		if Conn !=nil{
			Conn.Close()
		}

		startBackupProcess()
		go Start_timer(BackupDead)
		time.Sleep(4 * time.Second)
		Conn, err = net.Dial("tcp", "localhost:5000")
		if err != nil {
			fmt.Println("Failed to connect to backup server:", err)
			time.Sleep(5 * time.Second)
			
		}
		go sendCabHartBeat(elev)
		go handleConnection()
		
		<-BackupDead
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


func handleConnection() {
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
			Reset_timer()
		}
	}
}

func sendCabHartBeat(elev *elevator.Elevator){
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
		time.Sleep(5 * time.Second)
	}
	}


var timer *time.Timer

func Start_timer(backupDead chan bool) {
	for{
		timer = time.NewTimer(10 * time.Second)
		<-timer.C
		fmt.Println("backup process not detected, restarting...")
		backupDead <- true
	}
}
	

func Reset_timer(){
	timer.Reset(10 * time.Second)
}