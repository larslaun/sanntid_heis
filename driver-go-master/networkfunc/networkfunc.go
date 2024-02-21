package networkfunc

import (
	"fmt"
	"net"
	"Driver-go/elevator"
	"encoding/json"
	//"time"
)



//Initialize connection between master-slave
func InitConn(ip string, port string)(conn *net.UDPConn){
	var address = ip + ":" + port

	raddr, err := net.ResolveUDPAddr("udp", address)

	conn, err = net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return conn
}


func WriteStateToUDP(conn *net.UDPConn, elev elevator.Elevator){
	
	elevJSON , err := json.Marshal(elev)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn.Write(elevJSON)
}


func ReadFromUDP(conn *net.UDPConn)(msg elevator.Elevator){
	msgJSON := make([]byte, 1024)
	_ ,_ , err := conn.ReadFromUDP(msgJSON[0:])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	
	err = json.Unmarshal(msgJSON, &msg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	elevator.Elevator_print(msg)

	return msg
}