package util

import (
	"time"
	"fmt"
)


func Start_timer(parentAlive chan bool,parentDead chan bool) {
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-parentAlive:
			timer.Reset(10 * time.Second) // Reset the timer when the parent is alive
		case <-timer.C:
			fmt.Println("Parent process not detected, restarting...")
			parentDead <- true
		}
	}
}