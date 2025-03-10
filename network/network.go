package network

import(
	"root/reciver"
	"root/transmitter"
	"root/SharedData"
	"net"
	"fmt"
)
var elevator_1_ip = "localhost"


func Start_network(){
	var counter int
	var RemoteElevatorConnections =  make(map[string]net.Conn)

	for _, id := range sharedData.GetPossibleIDs(){
		if id == sharedData.GetElevatorID() {
			fmt.Println(counter)	
			counter +=1 
			
			continue // Local elevator not needed 
		}
		fmt.Println(id)
		fmt.Println(counter)

		if counter%2 == 0 {
		RemoteElevatorConnections[id] = transmitter.Start_tcp_call2("8080",elevator_1_ip)	
		}else{
		RemoteElevatorConnections[id] = reciver.Start_tcp_listen2("8080")
		}
		counter +=1 
		
	}
	sharedData.RemoteElevatorConnections = RemoteElevatorConnections
	reciver.SetConn()
	transmitter.SetConn()
}