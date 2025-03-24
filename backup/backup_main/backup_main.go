package main

import (
	"fmt"
	"time"
	"root/util"
	"os/exec"
	"strings"
	"runtime"
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
	go util.Msg_transmitter()
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
			time.Sleep(10*time.Second)
			return
		}
	}
}


func restart_elavator(){
	var cmd *exec.Cmd

	strCabBackup := strings.Trim(fmt.Sprint(CabBackup), "[]")

	switch runtime.GOOS {
		case "linux":
			gCommand := fmt.Sprintf("cd ../.. && go run main.go -isRestart=true -cabBackup='%s'", strCabBackup)
			cmd = exec.Command("gnome-terminal", "--", "bash", "-c", gCommand)

			fmt.Println(gCommand)
		case "windows":
		
			psCommand := fmt.Sprintf(
			"Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run main.go -isRestart=true -cabBackup='%s'\"",
			strCabBackup)
			cmd = exec.Command("powershell.exe", "-Command", psCommand)
			cmd.Dir = "../.." 

			 
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting PowerShell:", err)
		return
	}
}


