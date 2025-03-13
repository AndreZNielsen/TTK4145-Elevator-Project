package assigner

import (
	"root/config"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	

)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]config.Elevator_data `json:"states"`
}



func Assigner(localelvator config.Elevator_data,RemoteElevatorData map[string]config.Elevator_data, hallRequests [][2]bool) [][2]bool{
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


	states := map[string]config.Elevator_data{
		config.Elevator_id: localelvator,
	}

	// List of all possible elevator IDs.
	possibleIDs := []string{"A", "B", "C"}

	// Loop over possible IDs and add remote data if available.
	for _, id := range possibleIDs {
		if id == config.Elevator_id {
			continue // Local elevator already added.
		}
		// Only add the remote elevator if its data exists.
		if remote, ok := RemoteElevatorData[id]; ok {
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
	
	fmt.Printf("output: \n")
	for k, v := range output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	
	return output[config.Elevator_id]
}

