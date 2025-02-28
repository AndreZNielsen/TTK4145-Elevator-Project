package sharedData




var sharedHallRequests [4][2]bool // om buttonEvent eller update skal sharedHallRequests oppdateres
//sharedHallRequests er input i assigner.go og setAllLights() (skal skru på/av hallLys)
//output fra assigner.go skal oppdatere elevator.requests 

func UpdatesharedHallRequests(update [3]int){
	if update[2] == 1 && update[1] != 2{
	sharedHallRequests[update[0]][update[1]] = true
		
	}else if update[1] != 2{
	sharedHallRequests[update[0]][update[1]] = false
	}
}

func GetsharedHallRequests()[4][2]bool{
	return sharedHallRequests
}