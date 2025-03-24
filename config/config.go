package config

const Num_floors = 4

var Elevator_id = "A"


var PossibleIDs = []string{"A","B","C"}


var RemoteIDs = RemoveElement(PossibleIDs, Elevator_id)

var Elevatoip = map[string]string{
	"A": "localhost",
	"B": "localhost",
    "C": "localhost",
}

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
    Obstructed bool  
}

type HallRequests [][2]bool

type RemoteEvent struct {
	EventType     string
	Update        Update
	ElevatorData  Elevator_data
	HallRequests  HallRequests
	Id            string
}


func RemoveElement(slice []string, element string) []string {
    copiedSlice := make([]string, len(slice))
    copy(copiedSlice, slice)

    for i, v := range copiedSlice {
        if v == element {
            return append(copiedSlice[:i], copiedSlice[i+1:]...)
        }
    }
    return copiedSlice
}