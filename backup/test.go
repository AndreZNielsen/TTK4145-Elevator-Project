package backup

import(
	"fmt"
	"os/exec"
)


func main(){//skal fjernes
	Start_backup()
}
func start_backup(){
	psCommand := "Start-Process powershell -ArgumentList \"-NoExit\", \"-Command\", \"cd './backup_main'; go run backup_main.go\""

	// Start PowerShell and execute the command
	cmd := exec.Command("powershell.exe", "-Command", psCommand)
	// Start the backup program 
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting PowerShell:", err)
		return
	}
}

