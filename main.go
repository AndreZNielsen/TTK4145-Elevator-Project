package main

import (
	"fmt"
	elevalgo "root/elevator"

	"root/elevio"
	"root/utility"
)

func main() {
	fmt.Println("Started!")

	conn := utility.Start_tcp_call("8080","10.100.23.23")
	defer conn.Close()
	lis := utility.Start_tcp_listen("8081")
	defer lis.Close()


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
	go utility.Listen_recive(lis)
	
	

	for {
		select {
		case button := <-drv_buttons:
			elevalgo.FsmOnRequestButtonPress(button.Floor, elevalgo.Button(button.Button))
			e := elevalgo.GetElevatorRequests()
			utility.Send_tcp(conn,e)
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
				elevalgo.StartTimer() 
			}
		}
	}
}
