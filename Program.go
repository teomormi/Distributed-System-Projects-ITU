package main

import (
	"fmt"
)



func main() {
	counter1 := 0
	counter2 := 0
	counter3 := 0
	counter4 := 0
	counter5 := 0

	ch1 := make(chan bool)
	ch2 := make(chan bool)
	ch3 := make(chan bool)
	ch4 := make(chan bool)
	ch5 := make(chan bool)
	ch6 := make(chan bool)
	ch7 := make(chan bool)
	ch8 := make(chan bool)
	ch9 := make(chan bool)
	ch10 := make(chan bool)

	go fork(ch1, ch3)
	go fork(ch3, ch5)
	go fork(ch5, ch7)
	go fork(ch7, ch9)
	go fork(ch9, ch1)

	go philosopher(ch10, ch2)
	go philosopher(ch2, ch4)
	go philosopher(ch4, ch6)
	go philosopher(ch6, ch8)
	go philosopher(ch8, ch10)
}

func fork(ch1, ch2 chan bool){
	isTaken := false
	fmt.Println("hey")

}

func philosopher(ch1, ch2 chan bool) {
	state := "thinking"

	for {
		switch state {
		case "thinking":
		
			if (<-ch1 && <-ch2) {
				ch1 <- true
				ch2 <- true
				state = "eating"
			}
			break
		case "eating":
		
			state = "thinking"
			break
		}
	}
}
