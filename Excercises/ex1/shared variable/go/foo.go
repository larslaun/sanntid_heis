// Use `go run foo.go` to run your program

package main

import (
    . "fmt"
    "runtime"
    "time"
)

func incrementing(ch chan bool, finish chan bool) {
    //TODO: increment i 1000000 times
    for a := 0; a < 1000000; a++{
        ch<-true
    } 
    finish<-true

}

func decrementing(ch chan bool, finish chan bool) {
    //TODO: decrement i 1000000 times
    for a := 0; a < 1000000; a++{
       ch<-true
    }
    finish<-true
}

func server(ch1 chan bool, ch2 chan bool,finish chan bool, ch_read chan int) {
    var i = 0
    var k = 0
    for {
        select {
        case <-ch1:
            i++
        case <-ch2:
            i--
        case <-finish:
            k++
            if(k == 2){
                ch_read<-i
            }
            
        }
    }
    


}
func main() {
    // What does GOMAXPROCS do? What happens if you set it to 1?
    runtime.GOMAXPROCS(2)    
	
    // TODO: Spawn both functions as goroutines
    ch1:=make(chan bool)
    ch2:=make(chan bool)
    ch_read:=make(chan int)
    finish:=make(chan bool)
    
    go incrementing(ch1, finish)
    go decrementing(ch2, finish)
    
    go server(ch1, ch2,finish, ch_read)
    
	
    // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
    // We will do it properly with channels soon. For now: Sleep.
    time.Sleep(500*time.Millisecond)
    Println("The magic number is:", <-ch_read)
}
