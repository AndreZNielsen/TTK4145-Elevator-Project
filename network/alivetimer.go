package network

import(
	"time"
	"root/config"
)

var (
	timer = make(map[string]*time.Timer)
	ini = make(map[string]bool)
)

func InitAliveTimer(){
	for _, id := range config.RemoteIDs{
		ini[id]=false
	}
}


func StartAliveTimer(elvatorDead chan string , id string) {
	timer[id] = time.NewTimer(10 * time.Second)
	ini[id] = true
	<-timer[id].C
	elvatorDead <- id
}
	

func ResetAliveTimer(id string){
	if(ini[id]){
		timer[id].Reset(10 * time.Second) // Reset the timer when the elvator is alive
	}

}

func StopAliveTimer(id string){
	timer[id].Stop()
	ini[id] = false
}
