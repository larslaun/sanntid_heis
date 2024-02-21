package udpnetwork

import (
	"fmt"
	"net"
	"time"
)

func WriteToServerUDP() {
	raddr, err := net.ResolveUDPAddr("udp", ":20007")

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	message := []byte("Test")
	conn.Write(message)

	conn.Close()
}

func ReadfromServerUDP() {
	raddr, err := net.ResolveUDPAddr("udp", ":20007")

	conn, err := net.ListenUDP("udp", raddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	_, _, err = conn.ReadFromUDP(buffer[0:])

	if err != nil {
		panic(err)
		
	}

	

	
	fmt.Print("From server: ", string(buffer[0:]))

	conn.Close()
}
