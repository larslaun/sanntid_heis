package main

import (
	
	/* "ex4/backup"
	"ex4/primary" */
	"ex4/udpnetwork"
	"time"
	"fmt"
)




func main() {
	

	
	
	 if no message recieved {
		fmt.Print("Spawning backup\n")
		exec.Command("gnome-terminal", "--", "go", "run", "backup/backup.go").Run()

		
		for{
		
			udpnetwork.WriteToServerUDP()
			time.Sleep(time.Duration(1) * time.Second)

			fmt.Printf("\n%d\n", i)
			i++
			if i == 3 {
				os.Exit(0)
			}
		}
	} else{
		i++
	}



	/* go primary.Primary()	
	go backup.Backup()


	select{} */

}
