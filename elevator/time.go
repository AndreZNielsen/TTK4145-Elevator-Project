package elevator

import (
	"fmt"
	"time"

)

var (
	pollRate  = 20 * time.Millisecond
	timeOut   = 3 * time.Second
	timeOfStart time.Time
	timerActive    bool
	stuckTimeOfStart time.Time
	stuckTimerActive bool
	stuckTimeOut = 10 * time.Second
)

func StartTimer() {
	timeOfStart = time.Now()
	timerActive = true
}

func StopTimer() {
    timerActive = false
}

func TimerIsDone(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		timedOut := timerActive && time.Since(timeOfStart) > timeOut
		if timedOut && timedOut != prev {
			receiver <- true
		}
		prev = timedOut
	}
}

func TimedOut() bool {
	return timerActive && time.Since(timeOfStart) > timeOut
}

func StartStuckTimer(){
	fmt.Println("stucktimer started")
	stuckTimeOfStart = time.Now()
	stuckTimerActive = true
}

func StopStuckTimer() {
	fmt.Println("stucktimer stopped")
    stuckTimerActive = false
}

func StuckTimedOut() bool {
	return stuckTimerActive && time.Since(stuckTimeOfStart) > stuckTimeOut
}

func StuckTimerIsDone(stuckEvents chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		timedOut := stuckTimerActive && time.Since(timeOfStart) > timeOut
		if timedOut && timedOut != prev {
			stuckEvents <- true
		}
		prev = timedOut
	}
}

