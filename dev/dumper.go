package main

import (
	"time"
	"log"
)

var version = 1
func dump() {
	counter := 0
	//now := 'a'
	
	for {
		
		select {
		case <-time.After(time.Second):
			log.Printf("%d - %d...", version, counter)
			counter += 1
			continue
		}
	}
}

func main() {
	go dump()
	stop := make(chan int, 0)
	<-stop
}