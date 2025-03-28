package network

import(
	"root/receiver"
	"root/transmitter"
	"root/sharedData"
	"root/config"
	"root/elevator"
	"root/customStructs"
	"fmt"
	

)
var aliveRecievd = make(chan string)
var aliveTimeOut = make(chan string)

var requestHallRequests = make(chan customStructs.HallRequests)

func StartPeerNetwork(remoteEvent chan<- customStructs.RemoteEvent,disconnected chan<- string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn,elev *elevator.Elevator){
	transmitter.InitDiscEventChan(disconnected)
	transmitter.InitMutex()
	InitAliveTimer()

	for _, id := range config.RemoteIDs{

		if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(id) {// the elavator with the lowest index will dial 

		externalConn.RemoteElevatorConnections[id] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,id),config.Elevators_ip[id],id,externalConn)	

		}else{

		externalConn.RemoteElevatorConnections[id] = receiver.Start_tcp_listen(portGenerateor(config.Elevator_id,id),id,externalConn)
		}
		go StartAliveTimer(aliveTimeOut,id)
		
	}
	go handleAliveTimer(aliveRecievd,aliveTimeOut,externalConn,disconnected)
	go receiver.Listen_recive(remoteEvent,disconnected,sharedData,externalConn,aliveRecievd,requestHallRequests)
	go transmitter.Send_alive(externalConn)
	go handleRequestHallRequests(requestHallRequests,externalConn,sharedData)
	transmitter.Send_Elevator_data(elevator.GetElevatorData(elev), externalConn) 
	
}

func ReconnectPeer(remoteEvent chan<- customStructs.RemoteEvent,disconnected chan<- string, reConnID string,sharedData *sharedData.SharedData,externalConn *sharedData.ExternalConn,elev *elevator.Elevator){

	totalDisconnect := allFalse(externalConn.ConnectedConn)//checks if all the connections are down

	if indexOfElevatorID(config.Elevator_id)< indexOfElevatorID(reConnID) {// the elavator with the lowest index will dial 

		externalConn.RemoteElevatorConnections[reConnID] = transmitter.Start_tcp_call(portGenerateor(config.Elevator_id,reConnID),config.Elevators_ip[reConnID],reConnID,externalConn)	

		go receiver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn,aliveRecievd,requestHallRequests)

	}else{

		externalConn.RemoteElevatorConnections[reConnID] = receiver.Start_tcp_listen(portGenerateor(config.Elevator_id,reConnID),reConnID,externalConn)
		go receiver.Recive(remoteEvent,reConnID,disconnected,sharedData,externalConn,aliveRecievd,requestHallRequests)

	}

	go StartAliveTimer(aliveTimeOut,reConnID)


	if(totalDisconnect){//when it reenters the network it will request the hall requests from the first elevator in the list
		transmitter.MergeHallRequests(externalConn, sharedData.HallRequests, reConnID)
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
				disconnected <- id
			}
		}

	}
}

func handleRequestHallRequests(requestHallRequests chan customStructs.HallRequests,externalConn *sharedData.ExternalConn,sharedData *sharedData.SharedData){
	for{

		remoteHallRequests := <-requestHallRequests
		sharedData.HallRequests = mergeHallRequests(remoteHallRequests,sharedData.HallRequests)
		transmitter.Send_Hall_Requests(externalConn,sharedData.HallRequests)
	}
}

// mergeHallRequests merges two HallRequests into one new HallRequests
func mergeHallRequests(a, b customStructs.HallRequests) customStructs.HallRequests {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	// Create result slice
	result := make(customStructs.HallRequests, maxLen)

	for i := 0; i < maxLen; i++ {
		var aVal, bVal [2]bool

		if i < len(a) {
			aVal = a[i]
		}
		if i < len(b) {
			bVal = b[i]
		}
		result[i] = [2]bool{aVal[0] || bVal[0], aVal[1] || bVal[1]}
	}

	return result

}

