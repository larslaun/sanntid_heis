package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

func main() {

	raddr, _ := net.ResolveUDPAddr("udp", ":20011")
	recieve, _ := net.ListenUDP("udp", raddr)
	print(recieve)
	defer recieve.Close()
	print("This is slave\n")

	msg := 0

	for {

		buffer := make([]byte, 1024)
		recieve.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := recieve.ReadFromUDP(buffer[0:])

		if err != nil {
			break
		}

		msg, _ = strconv.Atoi(string(buffer[:n]))

		
		if err != nil {
			fmt.Println("Error:", err)
		}


		fmt.Printf("MSG VAL: %d\n", msg)
		fmt.Printf("Message recieved: %v\n", string(buffer[0:]))
	}
	recieve.Close()

	fmt.Printf("MSG: %d", msg)
	i := msg

	fmt.Print("Spawning backup\n")
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	time.Sleep(1 * time.Second)
	send, _ := net.DialUDP("udp", nil, raddr)
	defer send.Close()
	print("This is master\n")

	for {

		
		message := []byte(strconv.Itoa(i))

		_, err := send.Write(message)
		if err != nil {
			fmt.Println("Error:", err)
		}

		time.Sleep(time.Duration(1) * time.Second)

		fmt.Printf("\n%d\n", i)
		i++

	}

}
