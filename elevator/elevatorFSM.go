package elevator

import (
	//"fmt"
	"root/sharedData"
	Config "root/config"
	"root/elevio"
	//"runtime"
)

var (
	elevator Elevator
)

func MakeFsm() {
	elevator = MakeUninitializedelevator()
	//elevator parameters are set to default

	FsmOnInitBetweenFloors()
	//elevator is set to move down from the unknown start posistion
	//The elevator will now know what floor it is on, and will update its state accordingly
}

func GetElevatordata() Config.Elevator_data {
	return Config.Elevator_data{Behavior: EbToString(elevator.behaviour), Floor: elevator.floor, Direction: ElevioDirToString(elevator.direction), CabRequests: GetCabRequests(elevator.requests)}
}

func SetAllLights() {
	//Basically just takes the requests from the button presses and lights up the corresponding button lights
	requests := makeRequests(sharedData.GetsharedHallRequests(), GetCabRequests(elevator.requests))
	for floor := 0; floor < Num_floors; floor++ {
		for btn := 0; btn < Num_buttons; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, requests[floor][btn])
		}
	}
}

func FsmOnInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.direction = Dir_down
	elevator.behaviour = Behaviour_moving
}

func FsmOnRequestButtonPress(btn_floor int, btn_type Button) {
	var update [3]int

	switch elevator.behaviour {
	case Behaviour_door_open:
		if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
			StartTimer()
		} else {
			if btn_type == Btn_hallcab {
				elevator.requests[btn_floor][btn_type] = true
			}
	
			
			update = [3]int{btn_floor, int(btn_type), 1}
			go Transmitt_update_and_update_localHallRequests(update)
		}
	case Behaviour_moving:
		if btn_type == Btn_hallcab {
			elevator.requests[btn_floor][btn_type] = true
		}
		update = [3]int{btn_floor, int(btn_type), 1}
		go Transmitt_update_and_update_localHallRequests(update)


	case Behaviour_idle:
		if btn_type == Btn_hallup {
			elevator.requests[btn_floor][btn_type] = true
		}

		if elevator.floor == btn_floor {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = RequestsClearAtCurrentFloor(elevator)
			StartTimer()
			SetAllLights()
			elevator.behaviour = Behaviour_door_open
		}else {
			update = [3]int{btn_floor, int(btn_type), 1}
			go Transmitt_update_and_update_localHallRequests(update)
		}



	}
}

func FsmOnFloorArrival(newFloor int) {
	elevator.floor = newFloor

	elevio.SetFloorIndicator(elevator.floor)

	switch elevator.behaviour {
	case Behaviour_moving:
		if elevator.RequestsShouldStop() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			elevator = RequestsClearAtCurrentFloor(elevator)
			StartTimer()
			SetAllLights()
			elevator.behaviour = Behaviour_door_open
		}
	}
	go Send_Elevator_data(GetElevatordata())
	//fmt.Printf("\nNew state:\n")
	//elevator.print()
}

func FsmOnDoorTimeout() {
	switch elevator.behaviour {
	case Behaviour_door_open:
		pair := elevator.RequestsChooseDirection()
		elevator.direction = pair.dir
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case Behaviour_door_open:
			StartTimer()
			elevator = RequestsClearAtCurrentFloor(elevator)
			SetAllLights()
		case Behaviour_moving, Behaviour_idle:
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
		}
	}
	go Send_Elevator_data(GetElevatordata())
}

var doorObstructed bool

func DoorObstructed() {
	doorObstructed = true
	if elevator.behaviour == Behaviour_door_open {
		StartTimer()
	}
}

func DoorUnobstructed() {
	doorObstructed = false
	if elevator.behaviour == Behaviour_door_open {
		StartTimer()
	}
}

func IsDoorObstructed() bool {
	return doorObstructed
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

	// Extract columns 1 and 2 (index 0 and 1)
	for i := 0; i < len(matrix); i++ {
		newMatrix = append(newMatrix, [2]bool{matrix[i][0], matrix[i][1]})
	}
	return newMatrix
}



func makeRequests(HallRequests [][2]bool, GetCabRequests []bool) [Num_floors][3]bool {
    var result [Num_floors][3]bool

    for i := 0; i < Num_floors; i++ {
        result[i][0] = HallRequests[i][0]
        result[i][1] = HallRequests[i][1]
        result[i][2] = GetCabRequests[i]
    }
    return result
}

func GetElevator()Elevator{
	return elevator
}