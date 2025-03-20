package elevator

import (
	//"container/list"
	"fmt"
	"root/elevio"
	"root/sharedData"
)

type LocalEvent struct {
	EventType  string
	Button     elevio.ButtonEvent
	Floor      int
	Obstructed bool
}

// type RemoteEvent struct{ // This might be a good idea to implement
//     Floor int
//     Button int
//     Update int
// }

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

func FSM_HandleButtonPress(elevator *Elevator, btn_floor int, btn_type Button, externalData *sharedData.ExternalData) [3]int {
	
	noUpdate := [3]int{0, 0, 0} // Empty array, indicating no update


	if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
		fmt.Println("Request cleared immediately")
		DoorOpen(elevator) 
		return noUpdate
	}

	if btn_type == Btn_hallcab {
		elevator.requests[btn_floor][btn_type] = true
	}
	
	UpdateAndTransmittLocalRequests(elevator, btn_floor, btn_type, 1, externalData)
	return [3]int{btn_floor, int(btn_type), 1}
}

// func FSM_GetUpdate(elevator *Elevator, btn_floor int, btn_type Button) {
// 	if elevator.RequestsShouldClearImmediately(btn_floor, btn_type) {
// 		DoorOpen(elevator) // Just create an openDoor function. It is not clear what this does at the moment
// 		return
// 	}
// 	localUpdate := [3]int{btn_floor, int(btn_type), update}
// }


func UpdateAndTransmittLocalRequests(elevator *Elevator, btn_floor int, btn_type Button, update int, externalData *sharedData.ExternalData) {
	localUpdate := [3]int{btn_floor, int(btn_type), update}
	// go Transmitt_update_and_update_localHallRequests(elevator, localUpdate, externalData)
	Transmitt_update_and_update_localHallRequests(elevator, localUpdate, externalData)
}

// ControlMovement-function, start_if_idle is not good enough. We can make one that is also capable of stopping

func FSM_HandleFloorArrival(elevator *Elevator, newFloor int, externalData *sharedData.ExternalData) [][3]int {
	
	//noUpdate := [3]int{0, 0, 1} // One in last element, indicates no orders should be cleared
	var updates [][3]int

	elevator.floor = newFloor
	elevio.SetFloorIndicator(elevator.floor)

	if elevator.behaviour != Behaviour_moving {
		return updates
	}

	if elevator.ShouldStop() {
		elevio.SetMotorDirection(elevio.MD_Stop)
		updates = elevator.RequestsClearAtCurrentFloor(externalData)
		//SetAllLights(elevator)
		DoorOpen(elevator)
	}
	
	return updates

	// if elevator.behaviour == Behaviour_moving {
	// 	if elevator.ShouldStop() {
	// 		elevio.SetMotorDirection(elevio.MD_Stop)
	// 		elevator.RequestsClearAtCurrentFloor(externalData)
	// 		//SetAllLights(elevator)
	// 		DoorOpen(elevator)
	// 	}
	// }
	//go Send_Elevator_data(GetElevatorData(elevator), externalData)
}

func FSM_startNextRequest(elevator *Elevator, externalData *sharedData.ExternalData) {
	DoorClose(elevator)
	nextBehaviourPair := elevator.SelectNextDirection()
	elevator.direction = nextBehaviourPair.dir
	elevator.behaviour = nextBehaviourPair.behaviour

	switch elevator.behaviour {
	case Behaviour_door_open: //hvis neste tilstand er "door_open", skal døra åpnes
		DoorOpen(elevator)
		updates := elevator.RequestsClearAtCurrentFloor(externalData)
		if len(updates) == 0 {
			return
		}
		
		fmt.Println("lenght of updates:", len(updates))
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, externalData, updates[i])     // call this in main instead, as it requires externalData
			//transmitter.Send_update(update, connectedConn)

		}


		SetAllLights(elevator, externalData)
	case Behaviour_moving, Behaviour_idle:
		elevio.SetMotorDirection(elevio.MotorDirection(elevator.direction))
	}
}

// go Send_Elevator_data(GetElevatorData(elevator), externalData)
// }

// Here I think it would be reasonable to create functions for each case
// like FSM_HandleButtonEvent, FSM_HandleFloorEvent, FSM_HandleObstructedEvent, FSM_HandleTimerEvent
// The cases are very different from one another, so I think this makes sense.

// FSM_HandleButtonEvent should do what FSM_RequestButtonPress does now, but in an organized way
func FSM_HandleLocalEvent(elevator *Elevator, event LocalEvent, externalData *sharedData.ExternalData) {
	switch event.EventType {
	case "button":
		update := FSM_HandleButtonPress(elevator, event.Button.Floor, Button(event.Button.Button), externalData)
		
		if update == [3]int{0, 0, 0} { // replace this with if !update.isValid 
			return
		}

		
		UpdatesharedHallRequests(elevator, externalData, update)     // call this in main instead, as it requires externalData
		//transmitter.Send_update(update_val, connectedConn)
		
		AssignLocalHallRequests(elevator, externalData)      //This one is called anyway should be called elsewhere
		Start_if_idle(elevator)
		// startMotor() // doesnt exist yet, but this function should be created. Or something similar
		SetAllLights(elevator, externalData)

	case "floor":
		updates := FSM_HandleFloorArrival(elevator, event.Floor, externalData)

		if len(updates) == 0 {
			return
		}
		
		fmt.Println("lenght of updates:", len(updates))
		for i:=0; i<len(updates); i++  {
			UpdatesharedHallRequests(elevator, externalData, updates[i])     // call this in main instead, as it requires externalData
			//transmitter.Send_update(update, connectedConn)

		}
		
		//StopMotor() these two could be here, but the current solution might be good too
		//DoorOpen(elevator)
		AssignLocalHallRequests(elevator, externalData)
		SetAllLights(elevator, externalData)

	case "obstructed":
		FSM_HandleObstruction(elevator, event.Obstructed)

	case "timer":
		if IsDoorObstructed(elevator) {
			DoorOpen(elevator) // Door is kept open if it is obstructed
						
		} else {
			FSM_startNextRequest(elevator, externalData) 
		}
	}
}

func FSM_HandleRemoteEvent(elevator *Elevator, externalData *sharedData.ExternalData, event [3]int) { // Ideally this should say RemoteEvent, instead of [3]int, maybe fix this later
	UpdatesharedHallRequests(elevator, externalData, event)
	AssignLocalHallRequests(elevator, externalData)
	SetAllLights(elevator, externalData)
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