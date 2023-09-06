package main

import (
	"fmt"
	"sync"
)

var arbiter sync.Mutex

func main() {
	var phylo int = 5 // number of phylo to create
	var prev chan int
	first := make(chan int)
	for i := 0; i < phylo; i++ {
		next := make(chan int)
		if i < (phylo - 1) { // not the last
			if i == 0 { // the first
				go phylos(first, next)
			} else {
				go phylos(prev, next)
			}
			prev = make(chan int)
			go fork(next, prev)
		} else {
			go phylos(prev, next)
			go fork(next, first)
		}
	}
	for {
		// just to wait for the print
	}
}

func fork(left chan int, right chan int) {
	isTaken := false
	fmt.Println("hey")

	arbiter.Lock()
	fmt.Print(left)
	fmt.Print(" fork here ")
	fmt.Println(right)
	arbiter.Unlock()
}

func phylos(left chan int, right chan int) {
	arbiter.Lock()
	fmt.Print(left)
	fmt.Print(" phylo here ")
	fmt.Println(right)
	arbiter.Unlock()

	state := "thinking"

	for {
		switch state {
		case "thinking":

			if <-ch1 && <-ch2 {
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
