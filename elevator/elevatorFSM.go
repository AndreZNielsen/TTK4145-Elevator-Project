package elevator

import (
	//"container/list"
	"fmt"
	"root/elevio"
	"root/sharedData"
	"root/config"
	"root/transmitter"

)

type LocalEvent struct {
	EventType  string
	Button     elevio.ButtonEvent
	Floor      int
	Obstructed bool
}

func FSM_MakeElevator(elevator *Elevator, elevator_ip string, Num_floors int) {
	elevio.Init(elevator_ip, Num_floors)
	*elevator = MakeUninitializedelevator()
	FSM_InitBetweenFloors(elevator)
}

func FSM_InitBetweenFloors(elevator *Elevator) { // Create Move-down function
	elevio.SetMotorDirection(elevio.MD_Down)
	elevator.direction = Dir_down
	elevator.behaviour = Behaviour_moving
}

func FSM_HandleButtonPress(elevator *Elevator, btn_floor int, btn_type Button, SharedData *sharedData.SharedData) []config.Update {
	
	updates := []config.Update{}

	if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
		DoorOpen(elevator) 
		return updates
	}

	if btn_type == Btn_hallcab {
		elevator.requests[btn_floor][btn_type] = true
	}
	
	update := config.Update{
		Floor:		btn_floor, 
		ButtonType: int(btn_type), 
		Value: 		true,
	}
	return append(updates, update)
}

func FSM_HandleFloorArrival(elevator *Elevator, newFloor int, SharedData *sharedData.SharedData) []config.Update {
	
	updates := []config.Update{}

	elevator.floor = newFloor
	elevio.SetFloorIndicator(elevator.floor)

	if elevator.behaviour != Behaviour_moving {
		return updates
	}

	if elevator.ShouldStop() {
		elevio.SetMotorDirection(elevio.MD_Stop)
		updates = elevator.RequestsClearAtCurrentFloor(SharedData)
		DoorOpen(elevator)
	}
	
	return updates
}

func FSM_startNextRequest(elevator *Elevator, SharedData *sharedData.SharedData, externalConn *sharedData.ExternalConn) {
	DoorClose(elevator)
	nextBehaviourPair := elevator.SelectNextDirection()
	elevator.direction = nextBehaviourPair.dir
	elevator.behaviour = nextBehaviourPair.behaviour

	switch elevator.behaviour {
	case Behaviour_door_open: //hvis neste tilstand er "door_open", skal døra åpnes
		DoorOpen(elevator)
		updates := elevator.RequestsClearAtCurrentFloor(SharedData)
		if len(updates) == 0 {
			return
		}

		fmt.Println("lenght of updates:", len(updates))
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, SharedData, updates[i])
			transmitter.Send_update(updates[i], externalConn)

		}

	case Behaviour_moving, Behaviour_idle:
		elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
	}
	//Send_Elevator_data(GetElevatorData(elevator), externalConn)
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
		AssignLocalHallRequests(elevator, SharedData)
		Start_if_idle(elevator)
		// startMotor() // doesnt exist yet, but this function should be created. Or something similar
		

	case "floor":
		updates := FSM_HandleFloorArrival(elevator, event.Floor, SharedData)

		if len(updates) == 0 {
			return
		}
		
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, SharedData, updates[i])
			transmitter.Send_update(updates[i], externalConn)
		}
		
		//StopMotor() these two could be here, but the current solution might be good too
		//DoorOpen(elevator)
		AssignLocalHallRequests(elevator, SharedData)
		SetAllLights(elevator, SharedData)

	case "obstructed":
		FSM_HandleObstruction(elevator, event.Obstructed)
		//send_elevator_data for å sende obstruction

	case "timer":
		if IsDoorObstructed(elevator) {
			DoorOpen(elevator) // Door is kept open if it is obstructed
						
		} else {
			FSM_startNextRequest(elevator, SharedData, externalConn) 
			SetAllLights(elevator, SharedData)
		}
	}
}

func FSM_HandleRemoteEvent(elevator *Elevator, SharedData *sharedData.SharedData, event config.Update) { // Ideally this should say RemoteEvent, instead of [3]int, maybe fix this later
	UpdatesharedHallRequests(elevator, SharedData, event)
	AssignLocalHallRequests(elevator, SharedData)
	SetAllLights(elevator, SharedData)
	Start_if_idle(elevator) // should be called here instead of in ChangeLocalHallRequests

	// Once this change is made I am very happy with this function
}

func FSM_DetectLocalEvents(localEvents chan<- LocalEvent) {
	buttonEvents 		:= make(chan elevio.ButtonEvent)
	floorEvents 		:= make(chan int)
	obstructionEvents 	:= make(chan bool)
	timerEvents 		:= make(chan bool)

	go elevio.PollButtons(buttonEvents)
	go elevio.PollFloorSensor(floorEvents)
	go elevio.PollObstructionSwitch(obstructionEvents)
	go TimerIsDone(timerEvents)

	for {
		select {
		case button := <-buttonEvents:
			localEvents <- LocalEvent{EventType: "button", Button: button}
		case floor := <-floorEvents:
			localEvents <- LocalEvent{EventType: "floor", Floor: floor}
		case obstructed := <-obstructionEvents:
			localEvents <- LocalEvent{EventType: "obstructed", Obstructed: obstructed}
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