package customStructs


type Update struct { // This type will allow us to improve a few functions that use updates
    Floor       int
    ButtonType  int
    Value       bool
}



type Elevator_data struct {//data struct that contains all the data that the assigner needs to know about the elevator 
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"` 
    Obstructed  bool  
    Stuck       bool

}

type HallRequests [][2]bool

type RemoteEvent struct {
	EventType     string
	Update        Update
	ElevatorData  Elevator_data
	HallRequests  HallRequests
	Id            string
}