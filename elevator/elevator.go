package elevator

import (
	//"fmt"
	"time"
	"root/sharedData"
	"root/elevio"
	Config "root/config"
	"root/customStructs"
	"strings"
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

	floor     	int
	direction 	Dir
	Requests  	[Num_floors][Num_buttons]bool
	behaviour 	ElevatorBehaviour
	config   	elevatorConfig
	obstructed 	bool
	Stuck	   	bool

}

type elevatorConfig struct {
	clearRequestVariation ClearRequestVariant
	doorOpenDuration      time.Duration
}

type DirBehaviourPair struct {
	dir Dir
	behaviour ElevatorBehaviour
}

const (
	Num_floors  = Config.Num_floors
	Num_buttons = 3
)

type Dir int

const (

	Dir_down Dir = - 1
	Dir_stop Dir = 0
	Dir_up   Dir = 1

)

type Button int

const (

	Btn_hallup   Button = 0
	Btn_halldown Button = 1
	Btn_hallcab  Button = 2

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

func MakeUninitializedelevator() Elevator {
	return Elevator{
		floor:     -1,
		direction: Dir_stop,
		behaviour: Behaviour_idle,
		obstructed: false,

		Stuck: false,

		config: elevatorConfig{
			clearRequestVariation: CV_InDirn,
			doorOpenDuration:      3.0,

		},
	}
}

func IsDoorObstructed(elevator *Elevator) bool {
    return elevator.obstructed
}

func GetCabRequests(matrix [Num_floors][3]bool) []bool {
    var column []bool
    for i := 0; i < len(matrix); i++ {
        column = append(column, matrix[i][2])
    }
    return column
}

func GetHallRequests(matrix [Num_floors][3]bool) [][2]bool {
    var newMatrix [][2]bool

    for i := 0; i < len(matrix); i++ {
        newMatrix = append(newMatrix, [2]bool{matrix[i][0], matrix[i][1]})
    }
    return newMatrix
}

func MakeRequests(HallRequests [][2]bool, CabRequests []bool) [Num_floors][3]bool {
    var result [Num_floors][3]bool

    for i := 0; i < Num_floors; i++ {
        result[i][0] = HallRequests[i][0]
        result[i][1] = HallRequests[i][1]
        result[i][2] = CabRequests[i]
    }
    return result
}

func GetElevator(elevator *Elevator) Elevator {
    return *elevator
}

func GetElevatorData(elevator *Elevator) customStructs.Elevator_data {
    return customStructs.Elevator_data{
        Behavior:    EbToString(elevator.behaviour), 
        Floor:       elevator.floor, 
        Direction:   ElevioDirToString(elevator.direction), 
        CabRequests: GetCabRequests(elevator.Requests),
		Obstructed:  elevator.obstructed,
		Stuck:       elevator.Stuck,
    }
}

func SetAllLights(elevator *Elevator, SharedData *sharedData.SharedData) {
    requests := MakeRequests(SharedData.HallRequests, GetCabRequests(elevator.Requests))
    for floor := 0; floor < Num_floors; floor++ {
        for btn := 0; btn < Num_buttons; btn++ {
            elevio.SetButtonLamp(elevio.ButtonType(btn), floor, requests[floor][btn])
        }
    }
}

func RestorCabRequests(elevator *Elevator, cabBackup string){
	var cabBackupBool []bool

	// Split the string by space
	values := strings.Split(cabBackup, " ")
	// Convert each string into bool and append to the slice
	for _, v := range values {
		if v == "true" {
			cabBackupBool = append(cabBackupBool, true)
		} else if v == "false" {
			cabBackupBool = append(cabBackupBool, false)
		}
	}
	elevator.Requests = MakeRequests(GetHallRequests(elevator.Requests),cabBackupBool)
	Start_if_idle(elevator)
}
