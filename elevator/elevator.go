package elevator

import (
	"fmt"
	"time"
)

type ElevatorBehaviour int

const (
	Behaviour_idle = iota
	Behaviour_door_open
	Behaviour_moving
)

type ClearRequestVariant int

const (
	CV_All = iota

	CV_InDirn
)

type Elevator struct {
	floor     int
	direction Dir
	requests  [Num_floors][Num_buttons]bool
	behaviour ElevatorBehaviour
	config    config
}

type config struct {
	clearRequestVariation ClearRequestVariant
	doorOpenDuration      time.Duration
}

type DirBehaviourPair struct {
	dir Dir
	//direction of the elevator: DIR_UP, DIR_DOWN, DIR_STOP

	behaviour ElevatorBehaviour
	//states of the elevator: BEHAVIOUR_IDLE, BEHAVIOUR_DOOR_OPEN, BEHAVIOUR_MOVING
}

const (
	Num_floors  = 4
	Num_buttons = 3
)

type Dir int

const (
	Dir_down Dir = iota - 1
	Dir_stop
	Dir_up
)

type Button int

const (
	Btn_hallup Button = iota
	Btn_halldown
	Btn_hallcab
)


func ElevioDirToString(d Dir) string {
	switch d {
	case Dir_up:
		return "up"
	case Dir_down:
		return "down"
	case Dir_stop:
		return "stop"
	default:
		return "udefined"
	}
}

func ElevioButtonToString(b Button) string {
	switch b {
	case Btn_hallup:
		return "HallUp"
	case Btn_halldown:
		return "HallDown"
	case Btn_hallcab:
		return "Cab"
	default:
		return "undefined"
	}
}

func EbToString(behaviour ElevatorBehaviour) string {
	switch behaviour {
	case Behaviour_idle:
		return "idle"
	case Behaviour_door_open:
		return "doorOpen"
	case Behaviour_moving:
		return "moving"
	default:
		return "undefined"
	}
}

//this function just prints the current elevator status in the terminal
//If the code works properly at some point, any changes in the terminal that the simulator is run in
//should be visible in the terminal that the go-program is run in as well ;)

func (e *Elevator) print() {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |floor = %-2d          |\n", e.floor)
	fmt.Printf("  |dirn  = %-12.12s|\n", ElevioDirToString(e.direction))
	fmt.Printf("  |behav = %-12.12s|\n", EbToString(e.behaviour))

	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := Num_floors - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < Num_buttons; btn++ {
			if (f == Num_floors-1 && btn == int(Btn_hallup)) || (f == 0 && btn == int(Btn_halldown)) {
				fmt.Print("|     ")
			} else {
				if e.requests[f][btn] {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
					}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")


}

//Defalult elevator that starts in floor: -1, this doesnt make sense, but it does
//We cant initialize the elevator in a spesific floor, and PollFloorSensor() will update the variable to the correct
//floor as soon as the elevator starts moving i think

func MakeUninitializedelevator() Elevator {
	return Elevator{
		floor:     -1,
		direction: Dir_stop,
		behaviour: Behaviour_idle,
		config: config{
			clearRequestVariation: CV_InDirn,
			doorOpenDuration:      3.0,
		},
	}
}




