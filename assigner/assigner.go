package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/config"
	"root/sharedData"
	"runtime"
	"os"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]config.Elevator_data `json:"states"`
}



func Assigner(localelvator config.Elevator_data,RemoteElevatorData map[string]config.Elevator_data, hallRequests [][2]bool, externalConn sharedData.ExternalConn) [][2]bool{
	var input HRAInput
	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	// added this so that you dont have to run chmod +x on the executable
	err := os.Chmod("assignerExecutables/" + hraExecutable, 0755)
	if err != nil {
		fmt.Println("os.Chmod error: ", err)
		return nil
	}
	
	if localelvator.Floor == -1 { // kan kanje fjernes
		return make([][2]bool,config.Num_floors)
	}

	states := map[string]config.Elevator_data{//adds the local elevator to the states
		config.Elevator_id: localelvator,
	}



	// Loop over Remote IDs and add remote data if available.
	for _, id := range config.RemoteIDs {

		// Only add the remote elevator if its data exists and the elavator is not obstructed and the elavator is in network.
		if remote, ok := RemoteElevatorData[id]; ok && !remote.Obstructed && externalConn.ConnectedConn[id] {
			states[id] = remote
		}
	}

	input = HRAInput{
		HallRequests: hallRequests,
		States:       states,
	}


	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return nil
	}

	ret, err := exec.Command("assignerExecutables/" + hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return nil
	}

	output := make((map[string][][2]bool))
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return nil
	}
	
	//fmt.Printf("output: \n")
	//for k, v := range output {
		//fmt.Printf("%6v :  %+v\n", k, v)
	//}
	
	return output[config.Elevator_id] // returns the hall requests for the local elevator
}

