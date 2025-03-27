package main

import (
	"fmt"
	"os/exec"
	"root/util"
	"runtime"
	"strings"
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
		case CabBackup = <-alive: // Stores the cab requests from the elevator and resets the reviver timer
			util.Reset_timer()
		case <-dead:
			fmt.Println("Elevator is dead, restarting...")
			util.Conn.Close()
			restartElavator()
			return // Kills the reviver after restarting the elevator

		}
	}
}

func restartElavator(){
	var cmd *exec.Cmd

	strCabBackup := strings.Trim(fmt.Sprint(CabBackup), "[]")

	switch runtime.GOOS {
		case "linux":
			gCommand := fmt.Sprintf(// Adds the cab info to the restart command
			"go run main.go -isRestart=true -cabBackup='%s'; exec bash",
			strCabBackup)

			cmd = exec.Command("gnome-terminal", "--", "bash", "-c", gCommand)

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


