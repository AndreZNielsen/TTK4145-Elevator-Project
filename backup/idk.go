package backup

import(
	"fmt"
	"os/exec"
	"io"
	"encoding/json"
	"time"
	"bufio"
	"runtime"

)
type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
func Start_backup(){
for{
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
   		cmd = exec.Command("go", "run", "backup_main.go")
		cmd.Dir = "./backup/backup_main"
	case "windows":
	cmd = exec.Command("go", "run", "backup_main.go")
	cmd.Dir = "./backup/backup_main"
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return 
	}
	go read_from_backup(stdout)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return 
	}
	go send_to_backup(stdin)

	// Start the backup program 
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting PowerShell:", err)
		return
	}
	// Wait for the backup to exit		
	err = cmd.Wait()
	fmt.Println("Backup process exited, restarting in 5 seconds:", err)
	time.Sleep(5 * time.Second)
}
}


func send_to_backup(stdin io.WriteCloser){//sends the data to be backupt 
	for {
	message := Message{"message", "message recived"}
	jsonData, _ := json.Marshal(message)
	fmt.Fprintln(stdin, string(jsonData))
	fmt.Println("sendt")
	time.Sleep(5 * time.Second)
	}
}



func read_from_backup(stdout io.ReadCloser){ // kan nok fjernes 
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var msg Message
		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}
		fmt.Printf("Received: %+v\n", msg)
}
}