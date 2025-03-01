package utility

import (
    "encoding/gob"
    "fmt"
    sharedData "root/SharedData"
    "net"
    "sync"
    "time"
)

var conn_lift1 net.Conn
var sendMu sync.Mutex

func init() {
    gob.Register(ElevatorUpdate{})
}

func Start_tcp_call(port string, ip string) {
    var err error
    conn_lift1, err = net.Dial("tcp", ip+":"+port)
    if err != nil {
        fmt.Println("Error connecting to pc:", ip, err)
        time.Sleep(5 * time.Second)
        Start_tcp_call(port, ip)
    }
}

func Send_Elevator_data(data Elevator_data) {
    sendMu.Lock()         // Locking before sending
    defer sendMu.Unlock() // Ensure to unlock after sending
    time.Sleep(5 * time.Millisecond)
    encoder := gob.NewEncoder(conn_lift1)
    err := encoder.Encode("elevator_data") // Type ID
    if err != nil {
        fmt.Println("Encoding error:", err)
        return
    }
    time.Sleep(5 * time.Millisecond)
    err = encoder.Encode(data)
    if err != nil {
        fmt.Println("Error encoding data:", err)
        return
    }
}

func Send_ElevatorUpdate(update [3]int, data Elevator_data) {
    sendMu.Lock()         // Locking before sending
    defer sendMu.Unlock() // Ensure to unlock after sending
    time.Sleep(5 * time.Millisecond)

    elevatorUpdate := ElevatorUpdate{
        Update:        update,
        ElevatorState: data,
    }

    encoder := gob.NewEncoder(conn_lift1)
    err := encoder.Encode(elevatorUpdate)
    if err != nil {
        fmt.Println("Error encoding ElevatorUpdate:", err)
        return
    }
}

func Transmitt_update_and_update_localHallRequests(update [3]int, data Elevator_data) {
    sharedData.UpdatesharedHallRequests(update)
    Send_ElevatorUpdate(update, data)
}

type ElevatorUpdate struct {
    Update        [3]int
    ElevatorState Elevator_data
}