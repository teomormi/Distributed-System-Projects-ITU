package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)
	ch4 := make(chan string)
	ch5 := make(chan string)
	ch6 := make(chan string)
	ch7 := make(chan string)
	ch8 := make(chan string)
	ch9 := make(chan string)
	ch10 := make(chan string)

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

func fork(ch1, ch2 chan string){
	isTaken := false
	fmt.Println("hey")

}

func philosopher(ch1, ch2 chan string) {
	state := "thinking"

	switch state {
	case "thinking":
	
	state = "eating"
	break

	case "eating":
	
	state = "thinking"
	break
	}
}
