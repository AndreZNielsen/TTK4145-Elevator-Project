package elevator

import "time"

var (
	pollRate  = 20 * time.Millisecond
	timeOut   = 3 * time.Second
	timeOfStart time.Time
	On    bool
)

func StartTimer() {
	timeOfStart = time.Now()
	On = true
}

func WaitTime() {
	On = false
}

func PollTimer(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(pollRate)
		timedOut := On && time.Since(timeOfStart) > timeOut
		if timedOut && timedOut != prev {
			receiver <- true
		}
		prev = timedOut
	}
}

func TimedOut() bool {
	return On && time.Since(timeOfStart) > timeOut
}
