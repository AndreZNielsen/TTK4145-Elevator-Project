package elevator

import (
	"root/elevio"
	
)


func DoorOpen(elevator *Elevator) {
	StartTimer()
	elevio.SetDoorOpenLamp(true)
	elevator.behaviour = Behaviour_door_open
	StopStuckTimer(elevator)
}


func DoorClose(elevator *Elevator) {
	StopTimer()
	elevio.SetDoorOpenLamp(false)
	elevator.behaviour = Behaviour_idle
}
