package network

import(
	"root/reciver"
	"root/transmitter"
	"root/sharedData"
	"root/config"
	
	"fmt"
	

)

func StartPeerNetwork(remoteEvent chan<- config.Update,disconnected chan<- string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn){
	transmitter.InitDiscEventChan(disconnected)
	transmitter.InitMutex()


	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			continue // Local elevator not needed 
		}


		if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(id) {// the elavator with the lowest index will dial 

		externalConn.RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,externalConn)	
		}else{

		externalConn.RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id,externalConn)
		}

		
	}

	go reciver.Listen_recive(remoteEvent,disconnected,sharedData,externalConn)
	go transmitter.Send_alive(externalConn)
	

}

func ReconnectPeer(remoteEvent chan<- config.Update,disconnected chan<- string, reConnID string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn){



		if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(reConnID) {

		externalConn.RemoteElevatorConnections[reConnID] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,reConnID),config.Elevatoip[reConnID],reConnID,externalConn)	
		go reciver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn)

		}else{

		externalConn.RemoteElevatorConnections[reConnID] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,reConnID),reConnID,externalConn)
		go reciver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn)

		}
		
	}


func portGenerateor(localID, targetID string) string {
	localIndex := indexOfElevatorID(localID)
	targetIndex := indexOfElevatorID(targetID)
	port := 8000 + localIndex + targetIndex 


	return fmt.Sprintf("%d", port)
}

func indexOfElevatorID(target string) int {
    
	for i, v := range config.PossibleIDs {
        if v == target {
            return i
        }
    }
 	return -1 // if not in array 
}