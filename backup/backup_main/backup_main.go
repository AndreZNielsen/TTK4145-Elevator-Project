package main

import (
	"fmt"
	//"time"
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
	util.StartTCPLis()
	go util.Start_timer(dead)
	go util.SendHartBeat()
	go util.HandleConnection(alive)
	for{
		select {
		case CabBackup = <-alive:
			util.Reset_timer()
		case <-dead:
			fmt.Println("Elevator is dead, restarting...")
			util.Conn.Close()
			restart_elavator()
			return // Kills the backup after restarting the elevator

		}
	}
}


func restart_elavator(){
	var cmd *exec.Cmd

	strCabBackup := strings.Trim(fmt.Sprint(CabBackup), "[]")

	switch runtime.GOOS {
		case "linux":
			gCommand := fmt.Sprintf(// Adds the cab info to the restart command
			"go run main.go -isRestart=true -cabBackup='%s'",
			strCabBackup)

			cmd = exec.Command("gnome-terminal", "--", "bash", "-c", gCommand)

			fmt.Println(gCommand)
		case "windows":

			psCommand := fmt.Sprintf( // Adds the cab info to the restart command
			"Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run main.go -isRestart=true -cabBackup='%s'\"",
			strCabBackup)

			cmd = exec.Command("powershell.exe", "-Command", psCommand)

	}
	cmd.Dir = "../.." // Directory to use when launching the elevator process
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting :", err)
		return
	}
}


