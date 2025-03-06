package elevator

import(
	"root/assigner"
	"root/SharedData"
	"root/transmitter"
)

func UpdatesharedHallRequests(update [3]int){
	sharedHallRequests := sharedData.GetsharedHallRequests()
	if update[2] == 1 && update[1] != 2{//igneores updates to cab requests(update[1] != 2)
		sharedHallRequests[update[0]][update[1]] = true
			
		}else if update[1] != 2{
		sharedHallRequests[update[0]][update[1]] = false
		}
	sharedData.ChangeSharedHallRequests(sharedHallRequests)
	}
func Transmitt_update_and_update_localHallRequests(update_val [3]int, elevatorData sharedData.Elevator_data){ //sends the hall requests update to the other elevator and updates the local hall requests
	UpdatesharedHallRequests(update_val)
	transmitter.Send_update(update_val)
}

func ChangeLocalHallRequests(){
	elevator.requests = makeRequests(assigner.Assigner(GetElevatordata(), sharedData.GetRemoteElevatorData(),sharedData.GetsharedHallRequests()),GetCabRequests(elevator.requests))
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

func convertSliceToArray(slice [][2]bool) [4][2]bool {
	var arr [4][2]bool // Fixed-size array

	// Copy elements from slice to array (up to 4 elements)
	for i := 0; i < len(slice) && i < 4; i++ {
		arr[i] = slice[i]
	}

	return arr
}
