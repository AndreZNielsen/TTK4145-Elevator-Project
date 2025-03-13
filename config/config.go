package config

var Num_floors = 4

var Elevator_id = "A"

var PossibleIDs = []string{"A", "B"}

var RemoteIDs = RemoveElement(PossibleIDs, Elevator_id)

var Elevatoip = map[string]string{
	"A": "localhost",
	"B": "localhost",
}


type Elevator_data struct {//data struct that contains all the data that the assigner needs to know about the elevator 
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`    
}



func RemoveElement(slice []string, element string) []string {
    // Create a copy of the slice to avoid modifying the original underlying array.
    copiedSlice := make([]string, len(slice))
    copy(copiedSlice, slice)

    for i, v := range copiedSlice {
        if v == element {
            return append(copiedSlice[:i], copiedSlice[i+1:]...)
        }
    }
    return copiedSlice
}