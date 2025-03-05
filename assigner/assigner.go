package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	sharedData "root/SharedData"
	"root/utility"
	"runtime"
	"root/elevator"
	"os"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase



type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]utility.Elevator_data `json:"states"`
}

func convertArrayToSlice(arr [4][2]bool) [][2]bool {
	// Create a slice to hold the converted data
	slice := make([][2]bool, len(arr))

	// Copy elements from the array to the slice
	for i, v := range arr {
		slice[i] = v
	}

	return slice
}
[[false false] [false true] [true true] [false false]]
     B :  [[false false] [false false] [false false] [false true]]

func Assigner() {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := HRAInput{
		HallRequests: convertArrayToSlice(sharedData.GetsharedHallRequests()),
		States: map[string]utility.Elevator_data{
			"A": elevator.GetElevatordata(),
			"B": utility.GetRemoteElevatorData(),
		},
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	err = os.Chmod("assignerExecutables/" + hraExecutable, 0755)
	if err != nil {
		fmt.Println("Error setting executable permissions:", err)
		return
	}

	ret, err := exec.Command("assignerExecutables/" + hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	return (*output)["A"]
}
