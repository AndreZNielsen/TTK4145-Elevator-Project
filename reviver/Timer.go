package reviver
import(
	"time"
	"fmt"
)

var timer *time.Timer

func Start_timer(reviverDead chan bool) {
	timer = time.NewTimer(10 * time.Second)
	<-timer.C
	fmt.Println("reviver process not detected, restarting...")
	reviverDead <- true
}
	
func Reset_timer(){
	timer.Reset(10 * time.Second)
}

func Stop_timer(){
	timer.Stop()
}
