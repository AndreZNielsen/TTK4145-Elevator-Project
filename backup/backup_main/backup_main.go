package main

import(
	"fmt"
	"time"
	"os/exec"
)

var alive = make(chan bool)
var dead = make(chan bool)

func main(){
	fmt.Println("yo")
	time.Sleep(time.Second*5)
	go Msg_reciver(alive)
	go Start_timer(alive,dead)
	select{
	case <-alive:	
	go Msg_transmitter()
	case <-dead:
		start_elavator()
	}
}

func start_elavator(){
psCommand := "Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"cd '..';cd '..'; go run -ldflags='-X root/config.Elevator_id=A' main.go\""
cmd := exec.Command("powershell.exe", "-Command", psCommand)
err := cmd.Start()
if err != nil {
	fmt.Println("Error starting PowerShell:", err)
	return
}
}