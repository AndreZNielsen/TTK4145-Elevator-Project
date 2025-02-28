package utility

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"
)

var encoder *gob.Encoder
var sendMu sync.Mutex // Mutex to protect the sending of data

type Elevator_data struct {
	Behavior    string
	Floor       int
	Direction   string
	CabRequests []bool
	HallRequests [][2]bool
}
func InitEncoder(conn net.Conn) {
	encoder = gob.NewEncoder(conn)
}
func Start_tcp_call(port string, ip string) {
	var err error
	conn_lift1, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		time.Sleep(5 * time.Second)
		Start_tcp_call(port, ip) // Retry connection
	} else {
		InitEncoder(conn net.Conn)
	}

}

func Send_Elevator_data(data Elevator_data) {
	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending
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

func Send_update(update [3]int) {
	sendMu.Lock() // Locking before sending
	defer sendMu.Unlock() // Ensure to unlock after sending

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

