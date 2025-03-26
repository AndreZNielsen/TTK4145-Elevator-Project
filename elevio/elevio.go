package elevio

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const _pollRate = 20 * time.Millisecond

var _initialized bool = false
var NumFloors int
var _mtx sync.Mutex
var _mtx_initialized bool = false

var _conn net.Conn
var DisconnectedElevSever = make(chan bool)
var ServerAdrr string

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down MotorDirection = -1
	MD_Stop MotorDirection = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2
)

type Behaviour string

const (
	Idle     Behaviour = "idle"
	Moving   Behaviour = "moving"
	DoorOpen Behaviour = "doorOpen"
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

func Init(addr string, numFloors int) {
	if _initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	ServerAdrr = addr
	NumFloors = numFloors
	if !_mtx_initialized {
		_mtx = sync.Mutex{}
		_mtx_initialized = true
	}
		var err error
	for {
		_conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Error connecting to local elavator:",err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
	_initialized = true
	go handleDisconnect()
	return
	}
}

func SetMotorDirection(dir MotorDirection) {
	write([4]byte{1, byte(dir), 0, 0})
}

func SetButtonLamp(button ButtonType, floor int, value bool) {
	write([4]byte{2, byte(button), byte(floor), toByte(value)})
}

func SetFloorIndicator(floor int) {
	write([4]byte{3, byte(floor), 0, 0})
}

func SetDoorOpenLamp(value bool) {
	write([4]byte{4, toByte(value), 0, 0})
}

func SetStopLamp(value bool) {
	write([4]byte{5, toByte(value), 0, 0})
}

func PollButtons(receiver chan<- ButtonEvent) {
	prev := make([][3]bool, NumFloors)
	for {
		time.Sleep(_pollRate)
		for f := 0; f < NumFloors; f++ {
			for b := ButtonType(0); b < 3; b++ {
				v := GetButton(b, f)
				if v != prev[f][b] && v != false {
					receiver <- ButtonEvent{f, ButtonType(b)}
					fmt.Println("Button pressed:", f, b)
				}
				prev[f][b] = v
			}
		}
	}
}

func PollFloorSensor(receiver chan<- int) {
	prev := -1
	for {
		time.Sleep(_pollRate)
		v := GetFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}

func PollStopButton(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := GetStop()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func PollObstructionSwitch(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := GetObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func GetButton(button ButtonType, floor int) bool {
	a := read([4]byte{6, byte(button), byte(floor), 0})
	return toBool(a[1])
}

func GetFloor() int {
	for {
	if _initialized {
		a := read([4]byte{7, 0, 0, 0})
		if a[1] != 0 {
			return int(a[2])
		} else {
			return -1
		}
	} else {
		time.Sleep(1 * time.Second)
		continue
	}
}
}

func GetStop() bool {
	for {
		if _initialized {
			a := read([4]byte{8, 0, 0, 0})
			return toBool(a[1])
		} else {
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

func GetObstruction() bool {
	for {
		if _initialized {
	a := read([4]byte{9, 0, 0, 0})
	return toBool(a[1])
	} else {
		time.Sleep(1 * time.Second)
		continue
	}
}
}
func read(in [4]byte) [4]byte {
	var out [4]byte
	for {
		if _initialized {
			_mtx.Lock()
			defer _mtx.Unlock()

			_, err := _conn.Write(in[:])

			if err != nil {
				//panic("Lost connection to Elevator Server")
				DisconnectedElevSever<-true
			
				return [4]byte{0, 0, 0, 0}
			}

			_, err = _conn.Read(out[:])
			if err != nil {
				//panic("Lost connection to Elevator Server")
				DisconnectedElevSever<-true
				return [4]byte{0, 0, 0, 0}

			}
		} else {
			time.Sleep(1 * time.Second)
			continue
		}
	

		return out
	}
}

func write(in [4]byte) {
	for {
		if _initialized {
			_mtx.Lock()
			defer _mtx.Unlock()

			_, err := _conn.Write(in[:])
			if err != nil {
				//panic("Lost connection to Elevator Server")
				DisconnectedElevSever<-true
			}
		} else {
			time.Sleep(1 * time.Second)
			continue
		}
		return
}
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 {
		b = true
	}
	return b
}


func handleDisconnect(){
	
	<-DisconnectedElevSever
	_initialized = false
	go Init(ServerAdrr, 4)
	
}		
