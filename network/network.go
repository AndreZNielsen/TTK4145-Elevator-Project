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

func StartPeerNetwork(remoteEvent chan<- [3]int,disconnected chan<- string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn){
	transmitter.InitDiscEventChan(disconnected)
	transmitter.InitMutex()



	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			counter +=1 
			continue // Local elevator not needed 
		}


		if counter%2 == 0 {

		externalConn.RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,externalConn)	

		}else{

		RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id, externalData)
		}
		counter +=1 
		
	}


	go reciver.Listen_recive(remoteEvent,disconnected,sharedData,externalConn)
	go transmitter.Send_alive(externalConn)
	

}

func ReconnectPeer(remoteEvent chan<- [3]int,disconnected chan<- string, reConnID string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn){


	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			counter +=1 
			continue // Local elevator not needed 
		}




		externalConn.RemoteElevatorConnections[reConnID] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,reConnID),config.Elevatoip[reConnID],reConnID,externalConn)	
		go reciver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn)


		RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,disconnected, externalData)	
		reciver.SetConn(externalData)
		transmitter.SetConn(externalData)
		go reciver.Recive(receiver,id,disconnected, externalData)

		}else if needReconnecting == id{


		externalConn.RemoteElevatorConnections[reConnID] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,reConnID),reConnID,externalConn)
		go reciver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn)


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