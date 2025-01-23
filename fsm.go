package main

import (
	"fmt"
	"time"
)

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ElevatorState int

const (
	EB_Idle ElevatorState = iota
	EB_Moving
	EB_DoorOpen
)

type Event struct {
	Type       string
	Floor      int
	ButtonType int
}

type Elevator struct {
	Floor     int
	Dirn      int
	Behaviour ElevatorState
	Requests  [N_FLOORS][N_BUTTONS]int
	DoorTimer *time.Timer
	Config    ElevatorConfig
}

type ElevatorConfig struct {
	DoorOpenDuration time.Duration
}

func (e *Elevator) OnRequestButtonPress(floor int, buttonType int) {
	fmt.Printf("Button pressed: Floor %d, ButtonType %d\n", floor, buttonType)

	switch e.Behaviour {
	case EB_DoorOpen:
		e.Requests[floor][buttonType] = 1
		e.resetDoorTimer()
	case EB_Moving:
		e.Requests[floor][buttonType] = 1
	case EB_Idle:
		e.Requests[floor][buttonType] = 1
		e.Dirn, e.Behaviour = e.chooseDirection()
		if e.Behaviour == EB_DoorOpen {
			e.openDoor()
		} else if e.Behaviour == EB_Moving {
			e.move()
		}
	}
	e.updateLights()
}

func (e *Elevator) OnFloorArrival(newFloor int) {
	fmt.Printf("Arrived at Floor: %d\n", newFloor)
	e.Floor = newFloor

	if e.Behaviour == EB_Moving && e.shouldStop() {
		e.stop()
		e.openDoor()
	}
	e.updateLights()
}

func (e *Elevator) OnDoorTimeout() {
	fmt.Println("Door timeout")
	if e.Behaviour == EB_DoorOpen {
		e.Dirn, e.Behaviour = e.chooseDirection()
		if e.Behaviour == EB_Moving {
			e.move()
		}
	}
	e.updateLights()
}

func (e *Elevator) move() {
	fmt.Printf("Moving in direction: %d\n", e.Dirn)
}

func (e *Elevator) stop() {
	fmt.Println("Stopping at current floor")
	e.Behaviour = EB_DoorOpen
}

func (e *Elevator) openDoor() {
	fmt.Println("Door opened")
	e.Behaviour = EB_DoorOpen
	e.resetDoorTimer()
}

func (e *Elevator) resetDoorTimer() {
	if e.DoorTimer != nil {
		e.DoorTimer.Stop()
	}
	e.DoorTimer = time.AfterFunc(e.Config.DoorOpenDuration, func() {
		e.OnDoorTimeout()
	})
}

func (e *Elevator) shouldStop() bool {
	return e.Requests[e.Floor][0] == 1 || e.Requests[e.Floor][1] == 1 || e.Requests[e.Floor][2] == 1
}

func (e *Elevator) chooseDirection() (int, ElevatorState) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[floor][btn] == 1 {
				if floor > e.Floor {
					return 1, EB_Moving
				} else if floor < e.Floor {
					return -1, EB_Moving
				}
				return 0, EB_DoorOpen
			}
		}
	}
	return 0, EB_Idle
}

func (e *Elevator) updateLights() {
	fmt.Println("Updating lights based on requests")
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[floor][btn] == 1 {
				fmt.Printf("Light ON: Floor %d, Button %d\n", floor, btn)
			}
		}
	}
}

func simulateEvents(eventChan chan Event) {
	go func() {
		time.Sleep(2 * time.Second)
		eventChan <- Event{Type: "buttonPress", Floor: 1, ButtonType: 0}

		time.Sleep(3 * time.Second)
		eventChan <- Event{Type: "floorArrival", Floor: 1}

		time.Sleep(5 * time.Second)
		eventChan <- Event{Type: "doorTimeout"}
	}()
}

func main() {
	elevator := &Elevator{
		Floor:     0,
		Dirn:      0,
		Behaviour: EB_Idle,
		Config:    ElevatorConfig{DoorOpenDuration: 3 * time.Second},
	}

	eventChan := make(chan Event)

	go func() {
		for event := range eventChan {
			switch event.Type {
			case "buttonPress":
				elevator.OnRequestButtonPress(event.Floor, event.ButtonType)
			case "floorArrival":
				elevator.OnFloorArrival(event.Floor)
			case "doorTimeout":
				elevator.OnDoorTimeout()
			}
		}
	}()

	simulateEvents(eventChan)

	select {} // Keep the program running
}
