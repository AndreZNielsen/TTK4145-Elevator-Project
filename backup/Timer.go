package backup
import(
	"time"
	"fmt"
)


var timer *time.Timer

func Start_timer(backupDead chan bool) {
	timer = time.NewTimer(10 * time.Second)
	<-timer.C
	fmt.Println("backup process not detected, restarting...")
	backupDead <- true
}
	

func Reset_timer(){
	timer.Reset(10 * time.Second)
}

func Stop_timer(){
	timer.Stop()
}