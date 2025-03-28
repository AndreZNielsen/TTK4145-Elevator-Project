package config

const Num_floors = 4

var Elevator_id = "A"
var PossibleIDs = []string{"A", "B", "C"}
var LocalElevatorServerPort = "localhost:15657"
var RemoteIDs = RemoveElement(PossibleIDs, Elevator_id)
var Elevators_ip = map[string]string{
	"A": "localhost",
	"B": "localhost",
	"C": "localhost",
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
