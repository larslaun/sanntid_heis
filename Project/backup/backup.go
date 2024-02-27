package backup

import (
	"Elev-project/Network-go-master/network/bcast"
	"Elev-project/driver-go-master/elevator"
	"net"
	"time"
)


func BackupLoop(backupRX *net.UDPConn){
	print("\nThis is slave\n")

	backupRX.SetReadDeadline(time.Now().Add(2 * time.Second))

	for{


	}
}
