package backup

import(
	"fmt"
	"time"
	"os"
	"ex4/udpnetwork"
)


func Backup(){

	fmt.Print("\nhello\n")
	var i = 0
	
	for{
		fmt.Printf("\n%d\n", i)
		i++
		time.Sleep(time.Duration(1) * time.Second)
		
		udpnetwork.ReadfromServerUDP()

		if i == 3 {
			os.Exit(0)
		}
	}

}