package main

import (
	"fmt"
	elevalgo "root/elevator"
	"root/elevio"
	"root/reciver"
	"root/transmitter"
	"root/network"
)

var elevator_1_ip = "localhost"

/*
hvordan kjøre:
start to simulatorer med port 12345 og 12346 (./SimElevatorServer --port ______ i simulator mappen)
kjør go run -ldflags="-X root/assigner.elevator_id=A" main.go
og så go run -ldflags="-X root/assigner.elevator_id=B" main2.go
på samme maskin
*/


func main() {
	fmt.Println("Started!")


	
	/*
	go utility.Start_tcp_call2("8081", elevator_2_ip) // for the third elevator
	utility.Start_tcp_listen2("8081")
	*/
	elevio.Init("localhost:12345", elevalgo.NUM_FLOORS)

	elevalgo.MakeFsm()

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	poll_timer := make(chan bool)
	alive_timer := make(chan bool)
	update_recived := make(chan [3]int)
	network.Start_network(update_recived)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevalgo.PollTimer(poll_timer)
	go reciver.AliveTimer(alive_timer)
	go transmitter.Send_alive()

	elevalgo.Start_if_idle()
	transmitter.Send_Elevator_data(elevalgo.GetElevatordata())
	for {
		select {
		case button := <-drv_buttons:
			elevalgo.FsmOnRequestButtonPress(button.Floor, elevalgo.Button(button.Button))
			elevalgo.SetAllLights()
			
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
		case update := <-update_recived:
			elevalgo.UpdatesharedHallRequests(update)
			elevalgo.ChangeLocalHallRequests()

			elevalgo.SetAllLights()

		case <-alive_timer:

		}
		
	}
}

