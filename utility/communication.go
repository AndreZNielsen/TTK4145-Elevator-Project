package utility

import (
	"encoding/gob"
	"fmt"
	"os"

	//"os/exec"
	"net"
	//"time"
	//"bufio"
)

func Start_tcp_call(port string, ip string) net.Conn {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		os.Exit(1)
	}
	// må huske defer conn.Close() etter den er brukt
	return conn
}
func Start_tcp_listen(port string) net.Conn {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listen:", err)
		os.Exit(1)
	}
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}
	// må huske defer defer con.Close() etter den er brukt
	return conn
}

func Send_tcp(conn net.Conn, data [4][3]bool) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(data)
	if err != nil {	
		fmt.Println("Error encoding data:", err)
		return
	}
}

func Listen_recive(conn net.Conn) {
	for {
		Decode(conn)
	}
}

func Decode(conn net.Conn) {

	// Decode the received data
	var data [4][3]bool
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	fmt.Println("Received data:", data)
}
