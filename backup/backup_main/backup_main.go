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

	go util.Msg_reciver(alive)
	go util.Start_timer(alive,dead)
	for {
	select{
	case <-alive:	
	go util.Msg_transmitter()//kan kanje fjernes 
	case <-dead:
		time.Sleep(time.Second*5)
		restart_elavator()
		return
	}	
	}


}

func restart_elavator(){
	cmd := exec.Command("go", "run", "-ldflags=-X root/config.Elevator_id=A", "main.go")
	cmd.Dir = "../.."  
err := cmd.Start()
if err != nil {
	fmt.Println("Error starting PowerShell:", err)
	return
}
}