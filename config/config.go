package config

//shared config file for the elevators

import (
	"time"
)

const (
	NumFloors       = 4
	NumElevators    = 3
	NumButtons      = 3
	PeersPortNumber = 12345
	BcastPortNumber = 54321
	Buffer          = 1024

	DisconnectTime   = 1 * time.Second
	DoorOpenDuration = 3 * time.Second
	WatchdogTime     = 4 * time.Second
	HeartbeatTime    = 15 * time.Millisecond
)
