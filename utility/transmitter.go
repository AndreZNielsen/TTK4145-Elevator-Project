package utility

import (
	"encoding/gob"
	"fmt"
	//"os"
	//"os/exec"
	"net"
	"time"
	//"bufio"
)
var conn_lift1 net.Conn
//var conn_lift2 net.Conn



func Start_tcp_call(port string, ip string){
	var err error
	conn_lift1, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		time.Sleep(5*time.Second)
		Start_tcp_call(port, ip)
	}

}


func Send_Elevator_data(data Elevator_data) {
	encoder := gob.NewEncoder(conn_lift1)
	err := encoder.Encode("elevator_data") // Type ID
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
	err = encoder.Encode(data)
	if err != nil {	
		fmt.Println("Error encoding data:", err)
		return
	}
}

func Send_update(update [3]int){
	encoder := gob.NewEncoder(conn_lift1)
	err := encoder.Encode("int") // Type ID
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
	err = encoder.Encode(update)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
}



