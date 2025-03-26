package elevator

import (
	"fmt"
	"root/elevio"
	"root/sharedData"
	"root/transmitter"
	"root/customStructs"
)

type LocalEvent struct {
	EventType  string
	Button     elevio.ButtonEvent
	Floor      int
	Obstructed bool
	Stuck 	   bool
}

func FSM_MakeElevator(elevator *Elevator, elevator_ip string, Num_floors int) {
	elevio.Init(elevator_ip, Num_floors)
	*elevator = MakeUninitializedelevator()
	FSM_InitBetweenFloors(elevator)	
}


func FSM_InitBetweenFloors(elevator *Elevator) { // Create Move-down function
	if elevio.GetFloor() == -1 {
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.direction = Dir_down
	elevator.behaviour = Behaviour_moving
	} else {
		elevator.floor = elevio.GetFloor()
	}
}

func FSM_HandleButtonPress(elevator *Elevator, btn_floor int, btn_type Button, SharedData *sharedData.SharedData) []customStructs.Update {
	
	updates := []customStructs.Update{}

	if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
		DoorOpen(elevator) 
		return updates
	}
	if btn_type == Btn_hallcab {
		elevator.Requests[btn_floor][btn_type] = true
	}
	
	update := customStructs.Update{
		Floor:		btn_floor, 
		ButtonType: int(btn_type), 
		Value: 		true,
	}
	return append(updates, update)
}

func FSM_HandleFloorArrival(elevator *Elevator, newFloor int, SharedData *sharedData.SharedData) []customStructs.Update {
	
	updates := []customStructs.Update{}

	elevator.floor = newFloor
	elevio.SetFloorIndicator(elevator.floor)

	if elevator.behaviour != Behaviour_moving {
		return updates
	}

	if elevator.ShouldStop() {
		elevio.SetMotorDirection(elevio.MD_Stop)
		updates = elevator.RequestsClearAtCurrentFloor(SharedData)
		DoorOpen(elevator)

	} else{
		StartStuckTimer()
	}	
	return updates
}

func FSM_startNextRequest(elevator *Elevator, SharedData *sharedData.SharedData, externalConn *sharedData.ExternalConn) {
	DoorClose(elevator)

	nextBehaviourPair := elevator.SelectNextDirection()
	elevator.direction = nextBehaviourPair.dir
	elevator.behaviour = nextBehaviourPair.behaviour

	switch elevator.behaviour {
	case Behaviour_door_open: // if next state is "door_open", the door opens
		DoorOpen(elevator)
		updates := elevator.RequestsClearAtCurrentFloor(SharedData)
		if len(updates) == 0 {
			return
		}
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, SharedData, updates[i])
			transmitter.Send_update(updates[i], externalConn)
		}
	case Behaviour_moving:
		elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
		StartStuckTimer()
	
	case Behaviour_idle:
		elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))		
	}
}



func FSM_HandleLocalEvent(elevator *Elevator, event LocalEvent, SharedData *sharedData.SharedData, externalConn *sharedData.ExternalConn) {
	switch event.EventType {
	case "button":
		updates := FSM_HandleButtonPress(elevator, event.Button.Floor, Button(event.Button.Button), SharedData)
		
		if len(updates) == 0 {
			return
		}
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, SharedData, updates[i])
			transmitter.Send_update(updates[i], externalConn)
		}

		SetAllLights(elevator, SharedData)
		AssignLocalHallRequests(elevator, SharedData, *externalConn)
		Start_if_idle(elevator)
		

	case "floor":
		updates := FSM_HandleFloorArrival(elevator, event.Floor, SharedData)

		if len(updates) == 0 {
			return
		}
		
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, SharedData, updates[i])
			transmitter.Send_update(updates[i], externalConn)
		}
		
		AssignLocalHallRequests(elevator, SharedData, *externalConn)
		SetAllLights(elevator, SharedData)

	case "obstructed":
		FSM_HandleObstruction(elevator, event.Obstructed)
		//send_elevator_data for å sende obstruction
	case "stuck":
		elevator.Stuck = event.Stuck
		
	case "timer":
		if IsDoorObstructed(elevator) {
			DoorOpen(elevator,) // Door is kept open if it is obstructed						
		} else {
			FSM_startNextRequest(elevator, SharedData, externalConn) 
			SetAllLights(elevator, SharedData)
		}
	}
}

func FSM_HandleRemoteEvent(elevator *Elevator, SharedData *sharedData.SharedData, event customStructs.RemoteEvent, externalConn sharedData.ExternalConn) { // Ideally this should say RemoteEvent, instead of [3]int, maybe fix this later

	switch event.EventType {
	case "update":
		UpdatesharedHallRequests(elevator, SharedData, event.Update)
		
	case "elevatorData":
		SharedData.RemoteElevatorData[event.Id]=event.ElevatorData

	case "hallRequests":
		SharedData.HallRequests = event.HallRequests
	}
	
	AssignLocalHallRequests(elevator, SharedData, externalConn)
	SetAllLights(elevator, SharedData)
	Start_if_idle(elevator)
}

func FSM_DetectLocalEvents(localEvents chan<- LocalEvent) {
	buttonEvents 		:= make(chan elevio.ButtonEvent)
	floorEvents 		:= make(chan int)
	obstructionEvents 	:= make(chan bool)
	timerEvents 		:= make(chan bool)
	stuckEvents 		:= make(chan bool)


	go elevio.PollButtons(buttonEvents)
	go elevio.PollFloorSensor(floorEvents)
	go elevio.PollObstructionSwitch(obstructionEvents)
	go TimerIsDone(timerEvents)

	go StuckTimerIsDone(stuckEvents)

	for {
		select {
		case button := <-buttonEvents:
			localEvents <- LocalEvent{EventType: "button", Button: button}
		case floor := <-floorEvents:
			localEvents <- LocalEvent{EventType: "floor", Floor: floor}
		case obstructed := <-obstructionEvents:
			localEvents <- LocalEvent{EventType: "obstructed", Obstructed: obstructed}

		case stuck := <-stuckEvents:
			localEvents <- LocalEvent{EventType: "stuck", Stuck: stuck}
			if stuck {
				fmt.Println("stuck event happend, stuck is true")
			} else {
				fmt.Println("stuck event happend, stuck is false")
			}

		case <-timerEvents:
			localEvents <- LocalEvent{EventType: "timer"}
		}
	}
}

func FSM_HandleObstruction(elevator *Elevator, obstructed bool){
	if obstructed {
		elevator.obstructed = true	
	} else {
		elevator.obstructed = false
	}

}
