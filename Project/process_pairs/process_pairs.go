package process_pairs

import(
	"fmt"
)




//Putting process pair code from main here to clean up
//Decide if it should be included later



	//Processing pairs
	fmt.print("This is slave\n")
	timer1 := time.NewTimer(2 * time.Second)
	
	backupLoop:
		for {
			select {
			case elev = <-elevStateRx:
				fmt.Print("\n\nElev msg recieved:\n")
				elevator.Elevator_print(elev)
				fmt.Print("\n\n")
				timer1.Reset(2 * time.Second)
			case <-timer1.C:
				break backupLoop
			}
		}
	fmt.Print("Spawning backup\n")
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	print("This is now master\n")


