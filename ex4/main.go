package main

import (
	
	"fmt"
	"os/exec"	
	"os"
	"net"
	"time"
	//"strconv"
)




func main() {
	
	raddr, _ := net.ResolveUDPAddr("udp", ":20008")
	recieve, _ := net.ListenUDP("udp", raddr)

	defer recieve.Close()

	for{
		buffer := make([]byte, 1024)
		recieve.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _, err := recieve.ReadFromUDP(buffer[0:])

		fmt.Printf("Message recieved:",string(buffer[0:]))
		if err!=nil {
			break
		}
	}
	
	send, _ := net.DialUDP("udp", nil, raddr)
		
	defer send.Close()

	fmt.Print("Spawning backup\n")
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()

	for{
		
		i:=0
		i++
		if i == 5 {
			os.Exit(0)
		}

		message := []byte("test")

		_, err := send.WriteToUDP(message, raddr)
		if err != nil {
			fmt.Println("Error:", err)
		}

		time.Sleep(time.Duration(1) * time.Second)

		fmt.Printf("\n%d\n", i)
		

		
	}
	


	

}
