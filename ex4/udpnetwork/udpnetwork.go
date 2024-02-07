package udpnetwork

import (
	"fmt"
	"net"
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

	msg, err := net.ListenUDP("udp", raddr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buffer := make([]byte, 1024)

	msg.Read(buffer[0:])

	fmt.Print("From server: ", string(buffer[0:]))
	
	msg.Close()
}
