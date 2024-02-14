package main

import ( 
	"net"
	"fmt"
	"time"
)

func CheckError(err error){
	if err != nil{
		fmt.Println("Error: %d", err)
	}
}

func main(){
	raddr, err := net.ResolveUDPAddr("udp", ":20009")
	CheckError(err)
	recieve, err := net.ListenUDP("udp", raddr)
	CheckError(err)

	defer recieve.Close()
}