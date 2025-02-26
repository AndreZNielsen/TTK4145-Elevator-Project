package utility

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
	
)

var lis_lift1 net.Conn
//var lis_lift2 net.Conn


func Start_tcp_listen(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listen:", err)
	}
	lis_lift1, err = ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}

}

func Listen_recive() {
	for {
		Decode()
	}
}

func Decode() {

	// Decode the received data
	var data [4][6]int8
	decoder := gob.NewDecoder(lis_lift1)
	
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	fmt.Println("Received data:", data)
}