package network

import(
	"root/reciver"
	"root/transmitter"
	"root/sharedData"
	"root/config"
	"net"
	"sort"
	"fmt"
	

)

func Start_network(receiver chan<- [3]int,disconnected chan<- string, externalData *sharedData.ExternalData){
	var counter int
	var RemoteElevatorConnections =  make(map[string]net.Conn)

	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			counter +=1 
			continue // Local elevator not needed 
		}


		if counter%2 == 0 {

		RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,disconnected, externalData)	
		}else{

		RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id, externalData)
		}
		counter +=1 
		
	}
	externalData.RemoteElevatorConnections = RemoteElevatorConnections
	reciver.SetConn(externalData)
	transmitter.SetConn(externalData)
	go reciver.Listen_recive(receiver,disconnected, externalData)
	go transmitter.Send_alive(externalData)

	

}

func Network_reconnector(receiver chan<- [3]int,disconnected chan<- string, needReconnecting string, externalData *sharedData.ExternalData){
	var counter int
	var RemoteElevatorConnections =  make(map[string]net.Conn)
		counter = 0

	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			counter +=1 
			continue // Local elevator not needed 
		}



		if counter%2 == 0 && needReconnecting == id{

		RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,disconnected, externalData)	
		reciver.SetConn(externalData)
		transmitter.SetConn(externalData)
		go reciver.Recive(receiver,id,disconnected, externalData)

		}else if needReconnecting == id{

		RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id, externalData)
		reciver.SetConn(externalData)
		transmitter.SetConn(externalData)
		go reciver.Recive(receiver,id,disconnected, externalData)

		}
		counter +=1 
		
	}
	reciver.SetConn(externalData)
	transmitter.SetConn(externalData)
	}


func portGenerateor(localID, targetID string) string {
	// Combine the two IDs in a deterministic order
	ids := []string{localID, targetID}
	sort.Strings(ids) // ensures the order is consistent regardless of input order
	combined := ids[0] + ids[1]

	// makes a hash
	var hash int
	for _, ch := range combined {
		hash += int(ch)
	}

	//we choose 8000 as the base to make it in typical port range
	port := 8000 + (hash % 1000)
	return fmt.Sprintf("%d", port)
}