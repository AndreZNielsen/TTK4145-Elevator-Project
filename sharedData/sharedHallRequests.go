package sharedData

import (
	"net"
	"root/config"
	//"fmt"
)

// var RemoteElevatorConnections =  make(map[string]net.Conn)
// var Connected_conn = make(map[string]bool)


// var sharedHallRequests = make([][2]bool, config.Num_floors)
// //var RemoteElevatorData = Elevator_data{Behavior: "doorOpen",Floor: 0,Direction: "up",CabRequests: a}
// var RemoteElevatorData =  make(map[string]config.Elevator_data)
// om buttonEvent eller update skal sharedHallRequests oppdateres
//sharedHallRequests er input i assigner.go og setAllLights() (skal skru på/av hallLys)
//output fra assigner.go skal oppdatere elevator.requests 

//Vi må lage en update kanal-variabel/kanal som varsler hver gang det sendes eller mottas en TCP-melding 

//er det bedre/mulig å gjøre buttonEvent->sendTCPmessage til en update sånn at sharedHallRequests bare har ett input?
//Funksjonen som henter data til sharedHallRequests må hente data fra både når buttonEvent skjer(lokalt)
//og fra TCP-meldings-datastrukturen får en oppdatering. 


type SharedData struct {
	HallRequests [][2]bool
	RemoteElevatorData map[string]config.Elevator_data
	ObstrutedElevators map[string]bool
}

type ExternalConn struct {
	RemoteElevatorConnections map[string]net.Conn
	ConnectedConn map[string]bool
	 
}

func InitSharedData() *SharedData {
	return &SharedData{
		HallRequests:               make([][2]bool, config.Num_floors),
		RemoteElevatorData:         make(map[string]config.Elevator_data),
		//ObstrutedElevators:   		make(map[string]bool) ,	
	}
}


func InitExternalConn() *ExternalConn {
	return &ExternalConn{
		RemoteElevatorConnections:  make(map[string]net.Conn),
		ConnectedConn:              make(map[string]bool),
	}
}





	
// func GetsharedHallRequests()[][2]bool{
// 	return sharedHallRequests
// }
// func GetRemoteElevatorData()map[string]config.Elevator_data{
// 	return RemoteElevatorData
// }
// func UpdateRemoteElevatorData(NewRemoteElevatorData config.Elevator_data, id string){
// 	RemoteElevatorData[id] = NewRemoteElevatorData
// }
// func UpdateSharedHallRequests(NewSharedHallRequests [][2]bool){
// 	sharedHallRequests = NewSharedHallRequests
// }





