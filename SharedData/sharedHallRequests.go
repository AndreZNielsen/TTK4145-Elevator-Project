package sharedData

import (
	"net"
	"root/assigner"
	"sort"
	"fmt"
)

type Elevator_data = assigner.Elevator_data
var elevator_id = assigner.GetElevatorID()
var RemoteElevatorConnections =  make(map[string]net.Conn)
var Disconnected = make(chan string)
var Connected_conn = make(map[string]bool)


var NUM_FLOORS = 4
var sharedHallRequests = make([][2]bool, NUM_FLOORS)
//var RemoteElevatorData = Elevator_data{Behavior: "doorOpen",Floor: 0,Direction: "up",CabRequests: a}
var RemoteElevatorData =  make(map[string]Elevator_data)
// om buttonEvent eller update skal sharedHallRequests oppdateres
//sharedHallRequests er input i assigner.go og setAllLights() (skal skru på/av hallLys)
//output fra assigner.go skal oppdatere elevator.requests 

//Vi må lage en update kanal-variabel/kanal som varsler hver gang det sendes eller mottas en TCP-melding 

//er det bedre/mulig å gjøre buttonEvent->sendTCPmessage til en update sånn at sharedHallRequests bare har ett input?
//Funksjonen som henter data til sharedHallRequests må hente data fra både når buttonEvent skjer(lokalt)
//og fra TCP-meldings-datastrukturen får en oppdatering. 

var possibleIDs = []string{"A", "B"}
var remoteIDs = RemoveElement(possibleIDs, elevator_id)
	
func GetsharedHallRequests()[][2]bool{
	return sharedHallRequests
}
func GetRemoteElevatorData()map[string]Elevator_data{
	return RemoteElevatorData
}
func ChangeRemoteElevatorData(NewRemoteElevatorData Elevator_data, id string){
	RemoteElevatorData[id] = NewRemoteElevatorData
}
func ChangeSharedHallRequests(NewSharedHallRequests [][2]bool){
	sharedHallRequests = NewSharedHallRequests
}
func GetElevatorID() string{
	return elevator_id
}
func GetPossibleIDs()[]string{
	return possibleIDs
}
func GetRemoteIDs()[]string{
	return remoteIDs
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

func PortGenerateor(localID, targetID string) string {
	// Combine the two IDs in a deterministic order
	ids := []string{localID, targetID}
	sort.Strings(ids) // ensures the order is consistent regardless of input order
	combined := ids[0] + ids[1]

	// makes a hash
	var hash int
	for _, ch := range combined {
		hash += int(ch)
	}

	//we choose 8000 as the base to make it in typical port range
	port := 8000 + (hash % 1000)
	return fmt.Sprintf("%d", port)
}

