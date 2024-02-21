package primary

import(
	"fmt"
	"time"
	"os/exec"
	"ex4/udpnetwork"
)


func Primary() {
	fmt.Print("Spawning backup\n")
	exec.Command("gnome-terminal", "--", "go", "run", "backup/backup.go").Run()

	
	for{
	udpnetwork.WriteToServerUDP()
	time.Sleep(time.Duration(1) * time.Second)
	}
}



