package network

import(
	"root/reciver"
	"root/transmitter"
	"root/sharedData"
	"root/config"
	
	"fmt"
	

)

func Start_network(receiver chan<- [3]int,disconnected chan<- string,externalData *sharedData.SharedData,externalConn *sharedData.ExternalConn){


	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			continue // Local elevator not needed 
		}


		if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(id) {// the elavator with the lowest index will dial 

		externalConn.RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,disconnected)	
		}else{

		externalConn.RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id)
		}

		
	}

	go reciver.Listen_recive(receiver,disconnected)
	go transmitter.Send_alive(sharedData.SharedData)

	

}

func Network_reconnector(receiver chan<- [3]int,disconnected chan<- string, needReconnecting string,externalConn *sharedData.ExternalConn){


	for _, id := range config.PossibleIDs{
		if id == config.Elevator_id {
			continue // Local elevator not needed 
		}



		if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(id) && needReconnecting == id{

		externalConn.RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,disconnected)	
		go reciver.Recive(receiver,id,disconnected)

		}else if needReconnecting == id{

		externalConn.RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id)
		go reciver.Recive(receiver,id,disconnected)

		}
		
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