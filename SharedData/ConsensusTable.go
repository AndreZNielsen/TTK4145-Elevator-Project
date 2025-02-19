package sharedData
import (
	"root/elevator"
	"root/elevio"
	)

type LocalData struct {
	CabRequests [elevio.NUM_FLOORS]bool
	state elevator.Elevator
	FloorRequests [elevio.NUM_FLOORS][elevio.NUM_BUTTONS]bool
}	


type SharedData struct {
		hallRequests [elevio.NUM_FLOORS][elevio.NUM_BUTTONS]bool
		
}