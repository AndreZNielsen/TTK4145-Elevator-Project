package network

import(
	"time"
)
var (
	timer = make(map[string]*time.Timer)
)


func StartAliveTimer(elvatorDead chan string , id string) {
	timer[id] = time.NewTimer(10 * time.Second)
	<-timer[id].C
	elvatorDead <- id
}
	

func ResetAliveTimer(id string){
	timer[id].Reset(10 * time.Second) // Reset the timer when the elvator is alive
}
