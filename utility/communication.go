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

var lis_lift1 net.Conn
//var lis_lift2 net.Conn


func Start_tcp_call(port string, ip string){
	var err error
	conn_lift1, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("Error connecting to pc:", ip, err)
		time.Sleep(5*time.Second)
		Start_tcp_call(port, ip)
	}

}
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

func Send_tcp(data [4][3]bool) {
	encoder := gob.NewEncoder(conn_lift1)
	err := encoder.Encode(data)
	if err != nil {	
		fmt.Println("Error encoding data:", err)
		return
	}
}

func Listen_recive() {
	for {
		Decode()
	}
}

func Decode() {

	// Decode the received data
	var data [4][3]bool
	decoder := gob.NewDecoder(lis_lift1)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	fmt.Println("Received data:", data)
}