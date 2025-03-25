package reciver

import(
	"time"
)
var (
	pollRate  = 20 * time.Millisecond
	timeOut   = 3 * time.Second
	timeOfStart time.Time
	timerActive    bool
)

func StartTimer() {
	timeOfStart = time.Now()
	timerActive = true
}

func StopTimer() {
    timerActive = false
}

func TimedOut() bool {
	return timerActive && time.Since(timeOfStart) > timeOut
}
