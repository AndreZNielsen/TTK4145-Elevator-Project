package util

import (
	"time"
	"fmt"
)
var timer *time.Timer

func Start_timer(elvatorDead chan bool) {
	timer = time.NewTimer(10 * time.Second)
	<-timer.C
	fmt.Println("elvator process not detected, restarting...")
	elvatorDead <- true
}
	

func Reset_timer(){
	timer.Reset(10 * time.Second) // Reset the timer when the elvator is alive
}