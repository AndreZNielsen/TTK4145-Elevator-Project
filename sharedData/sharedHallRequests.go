package sharedData

import (
	"net"
	"root/config"
	"root/customStructs"
	//"fmt"
)




type SharedData struct {
	HallRequests customStructs.HallRequests
	RemoteElevatorData map[string]customStructs.Elevator_data
	ObstrutedElevators map[string]bool
}

type ExternalConn struct {
	RemoteElevatorConnections map[string]net.Conn
	ConnectedConn map[string]bool
	 
}

func InitSharedData() *SharedData {
	return &SharedData{
		HallRequests:               make(customStructs.HallRequests, config.Num_floors),
		RemoteElevatorData:         make(map[string]customStructs.Elevator_data),
		//ObstrutedElevators:   		make(map[string]bool) ,	
	}
}


func InitExternalConn() *ExternalConn {
	return &ExternalConn{
		RemoteElevatorConnections:  make(map[string]net.Conn),
		ConnectedConn:              make(map[string]bool),
	}
}








