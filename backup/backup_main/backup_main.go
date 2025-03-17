package main

import(
	"fmt"
	"time"
	"os/exec"
	"root/util"
)

var alive = make(chan bool)
var dead = make(chan bool)

func main(){
	fmt.Println("yo")

	go util.Msg_reciver(alive)
	go util.Start_timer(alive,dead)
	for {
	select{
	case <-alive:	
	go util.Msg_transmitter()
	case <-dead:
		time.Sleep(time.Second*5)
		start_elavator()
		return
	}	
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