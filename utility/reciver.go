package utility

import (
    "encoding/gob"
    "fmt"
    "net"
    "time"
    "root/SharedData"
)

var lis_lift1 net.Conn

type Elevator_data struct {
    Behavior    string
    Floor       int
    Direction   string
    CabRequests []bool
    HallRequests [][2]bool
}


func init() {
    gob.Register(ElevatorUpdate{})
}

func Start_tcp_listen(port string) {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Error starting listen:", err)
    }
    lis_lift1, err = ln.Accept()
    if err != nil {
        fmt.Println("Error accepting connection:", err)
    }
}

func Listen_recive(receiver chan<- bool) {
    for {
        Decode(receiver)
    }
}

func Decode(receiver chan<- bool) {
    decoder := gob.NewDecoder(lis_lift1)

    var elevatorUpdate ElevatorUpdate
    err := decoder.Decode(&elevatorUpdate)
    if err != nil {
        fmt.Println("Error decoding ElevatorUpdate:", err)
        time.Sleep(1 * time.Second)
        return
    }

    fmt.Println("Received ElevatorUpdate:", elevatorUpdate)
    sharedData.UpdatesharedHallRequests(elevatorUpdate.Update)
    // Process elevator state as needed
    receiver <- true
}