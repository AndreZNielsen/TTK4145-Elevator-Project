package main

import (
	"fmt"
	"time"
	"root/util"
	"os/exec"
	"strings"
	//"strconv"
	
)

var alive = make(chan []bool)
var dead = make(chan bool)

type Message struct {
	Type    string      `json:"type"`
	Content []bool `json:"content"`
}
var CabBackup []bool
func main() {
	go util.Start_timer(dead)
	util.StartTCPLis()
	go util.HandleConnection(alive)
	for{
		select {
		case CabBackup = <-alive:
			util.Reset_timer()
		case <-dead:
			time.Sleep(5 * time.Second)
			fmt.Println("Elevator is dead, restarting...")
			util.Conn.Close()
			restart_elavator()
			return
		}
	}
}


func restart_elavator(){
	strCabBackup := strings.Trim(fmt.Sprint(CabBackup), "[]")

	fmt.Println(strCabBackup)
	psCommand := fmt.Sprintf(
		"Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run main.go -isRestart=true -cabBackup='%s'\"",
		strCabBackup,
	)
	fmt.Println(psCommand)

 
	// Start PowerShell and execute the command
	cmd := exec.Command("powershell.exe", "-Command", psCommand)
	cmd.Dir = "../.."  
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting PowerShell:", err)
		return
	}
}
//psCommand := "Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run backup_main.go\""
