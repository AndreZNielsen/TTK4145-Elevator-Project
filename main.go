package main

import (
	"fmt"
	elevalgo "root/elevator"

	"root/elevio"
)

func main() {
	fmt.Println("Started!")

	elevio.Init("localhost:15657", elevalgo.NUM_FLOORS)

	elevalgo.MakeFsm()

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	poll_timer := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevalgo.PollTimer(poll_timer)

	for {
		select {
		case button := <-drv_buttons:
			elevalgo.FsmOnRequestButtonPress(button.Floor, elevalgo.Button(button.Button))
		case floor := <-drv_floors:
			if !elevalgo.IsDoorObstructed() {
				elevalgo.FsmOnFloorArrival(floor)
			}
		case obstructed := <-drv_obstr:
			if obstructed {
				elevalgo.DoorObstructed()
			} else {
				elevalgo.DoorUnobstructed()
			}
		case <-poll_timer:
			if !elevalgo.IsDoorObstructed() {
				elevalgo.StopTimer()
				elevalgo.FsmOnDoorTimeout()
			} else {
				elevalgo.StartTimer() // Restart the timer if the door is obstructed
			}
		}
	}
}
