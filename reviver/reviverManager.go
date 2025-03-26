package reviver

import (
	"fmt"
	"os/exec"
	"root/elevator"
	"runtime"
	"time"
)

type Message struct {
	Type    string      `json:"type"`
	Content []bool `json:"content"`
}



func StartReviver(elev *elevator.Elevator) {
	var reviverDead = make(chan bool)
	var reviverAlive = make(chan bool)

	fmt.Println("Starting reviver in 5 sec...")
	time.Sleep(5 * time.Second)
	
	startReviverProcess(elev,reviverAlive,reviverDead)


	for {// makes sure that the reviver is alive
		select{
			case<-reviverAlive:
				Reset_timer()
			case<-reviverDead:// Restarts the reviver if it dies
				Stop_timer()
				startReviverProcess(elev,reviverAlive,reviverDead)
		}
	}
}

func startReviverProcess(elev *elevator.Elevator,reviverAlive chan bool,reviverDead chan bool) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", "go run reviver_main.go")				
	
	case "windows":
		
		psCommand := "Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run reviver_main.go\""

		cmd = exec.Command("powershell.exe", "-Command", psCommand)
	}
	cmd.Dir = "./reviver/reviver_main" // Directory to use when launching the reviver process

	err := cmd.Run() // Starts the actual reviver program

	if err != nil {
		fmt.Println("Error starting reviver process:", err)
	}
	time.Sleep(2 * time.Second) // Waits for the reviver to start listening

	StartConn()	
	go Start_timer(reviverDead)
	go SendCabHartBeat(elev) 
	go HandleConnection(reviverAlive)
}



