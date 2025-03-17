package backup

import(
	"fmt"
	"os/exec"
	"encoding/json"
	"time"
	"bufio"

)
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
func Start_backup(){

	psCommand := "cd './backup'; cd './backup_main'; go run backup_main.go"

	// Start PowerShell and execute the command
	cmd := exec.Command("powershell.exe", "-Command", psCommand)
	// Start the backup program 
	go send_to_backup(cmd)
	go read_from_backup(cmd)
	time.Sleep(2 * time.Second)
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting PowerShell:", err)
		return
	}
	// Wait for the child to exit		
	err = cmd.Wait()
	fmt.Println("Backup process exited, restarting in 2 seconds:", err)
	time.Sleep(5 * time.Second)
	
}


func send_to_backup(cmd *exec.Cmd){
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return 
	}

	for {

	message := Message{"message", "message recived"}
	jsonData, _ := json.Marshal(message)
	fmt.Fprintln(stdin, string(jsonData))
	fmt.Println("sendt")

	time.Sleep(5 * time.Second)
	}
}



func read_from_backup(cmd *exec.Cmd){
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return 
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var msg Message
		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			fmt.Println("Error decoding child message:", err)
			continue
		}
		fmt.Printf("Parent received: %+v\n", msg)
}
}