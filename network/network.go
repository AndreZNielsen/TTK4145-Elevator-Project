package network

import(
	"root/reciver"
	"root/transmitter"
	"root/SharedData"
	"net"
	//"fmt"
	

)
var elevatoip = map[string]string{
        "A": "localhost",
        "B": "localhost",
}

func Start_network(receiver chan<- [3]int){
	var counter int
	var RemoteElevatorConnections =  make(map[string]net.Conn)

	for _, id := range sharedData.GetPossibleIDs(){
		if id == sharedData.GetElevatorID() {
			counter +=1 
			continue // Local elevator not needed 
		}


		if counter%2 == 0 {

		RemoteElevatorConnections[id] = transmitter.Start_tcp_call(sharedData.PortGenerateor(sharedData.GetElevatorID(),id),elevatoip[id],id)	
		}else{

		RemoteElevatorConnections[id] = reciver.Start_tcp_listen(sharedData.PortGenerateor(sharedData.GetElevatorID(),id),id)
		}
		counter +=1 
		
	}
	sharedData.RemoteElevatorConnections = RemoteElevatorConnections
	reciver.SetConn()
	transmitter.SetConn()
	go network_reconnector(receiver)
	go reciver.Listen_recive(receiver)
	

}

func network_reconnector(receiver chan<- [3]int){
	var counter int
	var RemoteElevatorConnections =  make(map[string]net.Conn)
	for {
		NeedReconnecting := <- sharedData.Disconnected
		counter = 0

		for _, id := range sharedData.GetPossibleIDs(){
			if id == sharedData.GetElevatorID() {
				counter +=1 
				continue // Local elevator not needed 
			}


	
			if counter%2 == 0 && NeedReconnecting == id{

			RemoteElevatorConnections[id] = transmitter.Start_tcp_call(sharedData.PortGenerateor(sharedData.GetElevatorID(),id),elevatoip[id],id)	
			reciver.SetConn()
			transmitter.SetConn()
			go reciver.Recive(receiver,id)

			}else if NeedReconnecting == id{

			RemoteElevatorConnections[id] = reciver.Start_tcp_listen(sharedData.PortGenerateor(sharedData.GetElevatorID(),id),id)
			reciver.SetConn()
			transmitter.SetConn()
			go reciver.Recive(receiver,id)

			}
			counter +=1 
			
		}
		NeedReconnecting = "nil"
		reciver.SetConn()
		transmitter.SetConn()
	}
	}


