package fsm

import (
	"fmt"
	"time"
	"elevio"
)

type Elevator struct {
	floor       int
	dirn        elevio.MotorDirection
	behaviour   Behaviour
	requests    [4][3]bool
	config      Config
}

type Config struct {
	doorOpenDuration_s float64
}

type Behaviour int

const (
	EB_Idle Behaviour = iota
	EB_Moving
	EB_DoorOpen
)

var elevator Elevator

func InitFSM() {
	// Initialize the elevator state
	elevator = Elevator{
		floor:     -1, // unknown floor
		dirn:      elevio.MD_Stop,
		behaviour: EB_Idle,
	}
	elevio.Init("127.0.0.1", 4) // Adjust IP and number of floors if needed
}

func setAllLights() {
	for floor := 0; floor < 4; floor++ {
		for btn := 0; btn < 3; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, elevator.requests[floor][btn])
		}
	}
}

func fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.dirn = elevio.MD_Down
	elevator.behaviour = EB_Moving
}

func fsm_onRequestButtonPress(btnFloor int, btnType elevio.ButtonType) {
	fmt.Printf("\n\nfsm_onRequestButtonPress(%d, %d)\n", btnFloor, btnType)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requestsShouldClearImmediately(btnFloor, btnType) {
			// Start door open timeout
			time.AfterFunc(time.Duration(elevator.config.doorOpenDuration_s)*time.Second, fsm_onDoorTimeout)
		} else {
			elevator.requests[btnFloor][btnType] = true
		}
	case EB_Moving:
		elevator.requests[btnFloor][btnType] = true
	case EB_Idle:
		elevator.requests[btnFloor][btnType] = true
		dirn, behaviour := chooseDirection()
		elevator.dirn = dirn
		elevator.behaviour = behaviour
		switch behaviour {
		case EB_DoorOpen:
			elevio.SetDoorOpenLamp(true)
			time.AfterFunc(time.Duration(elevator.config.doorOpenDuration_s)*time.Second, fsm_onDoorTimeout)
			clearAtCurrentFloor()
		case EB_Moving:
			elevio.SetMotorDirection(elevator.dirn)
		case EB_Idle:
			// Stay idle, do nothing
		}
	}

	setAllLights()
}

func fsm_onFloorArrival(newFloor int) {
	fmt.Printf("\n\nfsm_onFloorArrival(%d)\n", newFloor)
	elevator.floor = newFloor
	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if shouldStop() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			clearAtCurrentFloor()
			time.AfterFunc(time.Duration(elevator.config.doorOpenDuration_s)*time.Second, fsm_onDoorTimeout)
			setAllLights()
			elevator.behaviour = EB_DoorOpen
		}
	}
}

func fsm_onDoorTimeout() {
	fmt.Println("\n\nfsm_onDoorTimeout()")
	
	switch elevator.behaviour {
	case EB_DoorOpen:
		dirn, behaviour := chooseDirection()
		elevator.dirn = dirn
		elevator.behaviour = behaviour
		
		switch behaviour {
		case EB_DoorOpen:
			time.AfterFunc(time.Duration(elevator.config.doorOpenDuration_s)*time.Second, fsm_onDoorTimeout)
			clearAtCurrentFloor()
			setAllLights()
		case EB_Moving, EB_Idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevator.dirn)
		}
	}
}

func requestsShouldClearImmediately(floor int, btnType elevio.ButtonType) bool {
	// Implement logic for request clearing
	return false
}

func chooseDirection() (elevio.MotorDirection, Behaviour) {
	// Choose direction and behaviour based on requests
	return elevio.MD_Up, EB_Moving
}

func clearAtCurrentFloor() {
	// Clear requests at the current floor
}

func shouldStop() bool {
	// Check if the elevator should stop
	return false
}

