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

func start_tcp_call(port string, ip string) net.Conn {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		os.Exit(1)
	}
	// må huske defer conn.Close() etter den er brukt
	return conn
}
func start_tcp_listen(port string) net.Listener {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listen:", err)
		os.Exit(1)
	}
	// må huske defer defer ln.Close() etter den er brukt
	return ln
}

func send_tcp(conn net.Conn, data []int) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
}

func listen_recive(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		go decode(conn)
	}
}

func decode(conn net.Conn) {
	defer conn.Close()

	// Decode the received data
	var data []int
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	fmt.Println("Received data:", data)
}
