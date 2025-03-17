package main

import (
	"time"
	"fmt"
)


func Start_timer(parentAlive chan bool,parentDead chan bool) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-parentAlive:
			timer.Reset(5 * time.Second) // Reset the timer when the parent is alive
		case <-timer.C:
			fmt.Println("Parent process not detected, restarting...")
			parentDead <- true
		}
	}
}