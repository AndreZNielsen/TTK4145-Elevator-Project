package elevator

import(
	"root/assigner"
	"root/SharedData"
	"root/transmitter"
	"root/elevio"
	"fmt"
)

func UpdatesharedHallRequests(update [3]int){
	sharedHallRequests := sharedData.GetsharedHallRequests()
	if update[2] == 1 && update[1] != 2{//igneores updates to cab requests(update[1] != 2)
		sharedHallRequests[update[0]][update[1]] = true
			
		}else if update[1] != 2{
		sharedHallRequests[update[0]][update[1]] = false
		}
	sharedData.ChangeSharedHallRequests(sharedHallRequests)
	ChangeLocalHallRequests()
	}
func Transmitt_update_and_update_localHallRequests(update_val [3]int){ //sends the hall requests update to the other elevator and updates the local hall requests
	UpdatesharedHallRequests(update_val)
	transmitter.Send_update(update_val)
}

func ChangeLocalHallRequests(){
	fmt.Println(GetElevatordata())
	fmt.Println(sharedData.GetRemoteElevatorData())
	if GetElevatordata().Floor != -1 && !(GetElevatordata().Floor == 0 && GetElevatordata().Direction == "down") && !(GetElevatordata().Floor == 3 && GetElevatordata().Direction == "up") {//stops the elavator data form crashing the assigner 

	elevator.requests = makeRequests(assigner.Assigner(GetElevatordata(), sharedData.GetRemoteElevatorData(),sharedData.GetsharedHallRequests()),GetCabRequests(elevator.requests))
	Start_if_idle()
	SetAllLights()
	elevator.print()
}
}

func Send_Elevator_data( elevatorData sharedData.Elevator_data){
	transmitter.Send_Elevator_data(elevatorData)
}
func Start_if_idle(){
	switch elevator.behaviour{
	case BEHAVIOUR_IDLE:	
		pair := elevator.RequestsChooseDirection()
		elevator.direction = pair.dir
		elevator.behaviour = pair.behaviour
		elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
	case BEHAVIOUR_DOOR_OPEN:
		StartTimer()
		elevio.SetDoorOpenLamp(true)

}
	}


