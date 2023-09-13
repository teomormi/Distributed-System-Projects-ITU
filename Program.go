package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var phylo int = 5 // number of phylo to create
	var prev chan bool
	first := make(chan bool)
	for i := 0; i < phylo; i++ {
		next := make(chan bool)
		if i < (phylo - 1) { // not the last
			if i == 0 { // the first
				go phylos(i, first, next)
			} else {
				go phylos(i, prev, next)
			}
			prev = make(chan bool)
			go fork(next, prev)
		} else {
			go phylos(i, prev, next)
			go fork(next, first)
		}
	}
	for {
		// just to wait for the print
	}
}

func fork(left chan bool, right chan bool) {
	isTaken := false

	for {
		select {
		case message := <-left:
			if message {
				if isTaken {
					left <- false
				} else {
					left <- true
					isTaken = true
				}
			} else {
				isTaken = false
			}
		case message := <-right:
			if message {
				if isTaken {
					right <- false
				} else {
					right <- true
					isTaken = true
				}
			} else {
				isTaken = false
			}
		default:
		}
	}
}

func phylos(id int, left chan bool, right chan bool) {
	state := "thinking"
	counter := 0

	for {
		switch state {
		case "thinking":
			left <- true
			if <-left {
				right <- true
				if <-right {
					state = "eating"
					counter++
					fmt.Printf("%d I'm eating for the %d time \n", id, counter)
				} else {
					left <- false
				}
			}
			break
		case "eating":
			left <- false
			right <- false
			state = "thinking"
			fmt.Printf("%d I'm thinking\n", id)
			break
		}
		n := rand.Intn(1000)
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
}
