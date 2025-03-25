package network

import(
	"root/reciver"
	"root/transmitter"
	"root/sharedData"
	"root/config"
	"root/elevator"
	"fmt"
	

)
var aliveRecievd = make(chan string)
var aliveTimeOut = make(chan string)
var requestHallRequests = make(chan string)

func StartPeerNetwork(remoteEvent chan<- config.RemoteEvent,disconnected chan<- string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn){
	transmitter.InitDiscEventChan(disconnected)
	transmitter.InitMutex()
	//InitAlive()

	for _, id := range config.RemoteIDs{


		if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(id) {// the elavator with the lowest index will dial 

		externalConn.RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevatoip[id],id,externalConn)	
		}else{

		externalConn.RemoteElevatorConnections[id] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id,externalConn)
		}
		go StartAliveTimer(aliveTimeOut,id)
		
	}
	go handleAliveTimer(aliveRecievd,aliveTimeOut,externalConn,disconnected)
	go reciver.Listen_recive(remoteEvent,disconnected,sharedData,externalConn,aliveRecievd,requestHallRequests)
	go transmitter.Send_alive(externalConn)
	go handleRequestHallRequests(requestHallRequests,externalConn,sharedData)
	

}

func ReconnectPeer(remoteEvent chan<- config.RemoteEvent,disconnected chan<- string, reConnID string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn,elev *elevator.Elevator){

	totalDicvonnect := allFalse(externalConn.ConnectedConn)

	if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(reConnID) {

		externalConn.RemoteElevatorConnections[reConnID] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,reConnID),config.Elevatoip[reConnID],reConnID,externalConn)	
		go reciver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn,aliveRecievd,requestHallRequests)

	}else{

		externalConn.RemoteElevatorConnections[reConnID] = reciver.Start_tcp_listen(portGenerateor(config.Elevator_id,reConnID),reConnID,externalConn)
		go reciver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn,aliveRecievd,requestHallRequests)

	}

	StartAliveTimer(aliveTimeOut,reConnID)

	if(totalDicvonnect){
		transmitter.RequestHallRequests(externalConn,reConnID)
	}
	transmitter.Send_Elevator_data(elevator.GetElevatorData(elev), externalConn) 
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

func allFalse(m map[string]bool) bool {
	for _, v := range m {
		if v {
			return false
		}
	}
	return true
}

func handleAliveTimer(aliveRecievd chan string,aliveTimeOut chan string,externalConn *sharedData.ExternalConn,disconnected chan<- string){
	for{
		select{
		case id := <-aliveRecievd:
			ResetAliveTimer(id)

		case id := <-aliveTimeOut:
			if externalConn.ConnectedConn[id]{
				fmt.Println("handleAliveTimer triggered disconnect")
				disconnected <- id
			}
		}

	}
}
func handleRequestHallRequests(requestHallRequests chan string,externalConn *sharedData.ExternalConn,sharedData *sharedData.SharedData){
	for{
		id := <-requestHallRequests
		transmitter.Send_Hall_Requests(id,externalConn,sharedData)
	}
}

