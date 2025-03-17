package util

import (
	"time"
	"fmt"
)


func Start_timer(elvatorAlive chan bool,elvatorDead chan bool) {
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-elvatorAlive:
			timer.Reset(10 * time.Second) // Reset the timer when the elvator is alive
		case <-timer.C:
			fmt.Println("elvator process not detected, restarting...")
			elvatorDead <- true
		}
	}
}