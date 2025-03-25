package backup

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



func Start_backup(elev *elevator.Elevator) {
	var backupDead = make(chan bool)
	var backupAlive = make(chan bool)

	fmt.Println("Starting backup in 5 sec...")
	time.Sleep(5 * time.Second)
	
	startBackupProcess(elev,backupAlive,backupDead)


	for {// makes sure that the backup is alive
		select{
			case<-backupAlive:
				Reset_timer()
			case<-backupDead:// Restarts the backup
				Stop_timer()
				startBackupProcess(elev,backupAlive,backupDead)
		}
	}
}

func startBackupProcess(elev *elevator.Elevator,backupAlive chan bool,backupDead chan bool) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", "go run backup_main.go")				
	
	case "windows":
		
		psCommand := "Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"go run backup_main.go\""

		cmd = exec.Command("powershell.exe", "-Command", psCommand)
	}
	cmd.Dir = "./backup/backup_main" // Directory to use when launching the backup process

	err := cmd.Run() // Starts the actual backup program

	if err != nil {
		fmt.Println("Error starting backup process:", err)
	}
	time.Sleep(2 * time.Second) // Waits for the backup to start listening

	StartConn()	
	go Start_timer(backupDead)
	go SendCabHartBeat(elev) 
	go HandleConnection(backupAlive)
}



