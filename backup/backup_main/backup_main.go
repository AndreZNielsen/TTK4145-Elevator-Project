package main

import(
	"fmt"
	"time"
	"os/exec"
)

var alive = make(chan bool)
var dead = make(chan bool)

func main(){
	var parentCmd *exec.Cmd
	fmt.Println("yo")
	time.Sleep(time.Second*5)
	go Msg_reciver(alive)
	go Start_timer(alive,dead)
	select{
		case <-dead:
			fmt.Println("elavator process not detected, restarting")
			if parentCmd != nil {
				parentCmd.Process.Kill()
			}
			parentCmd = start_elavator()
	}
}

func start_elavator()*exec.Cmd{

}