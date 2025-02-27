package sharedData

import (
	elevalgo "root/elevator"
)




 var sharedHallRequests [4][2]bool // om buttonEvent eller update skal sharedHallRequests oppdateres
 //sharedHallRequests er input i assigner.go og setAllLights() (skal skru på/av hallLys)
 //output fra assigner.go skal oppdatere elevator.requests 
