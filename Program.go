/*
Authors:
 - Lucas Roy Guldbrandsen
 - Rafael Steffen Nguyen Jensen
 - Matteo Mormile

Our code that prevents a dealock works like so:
If a phylo is in the thinking state it will first make a request to get the left fork
If the the left fork is available it wil then make a request to the right fork
If the right fork is not available then it will return the left fork ensuring that a case where every phylo only has one fork will not occur
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	main := make(chan bool) //channel for the phylos to communicate with main if they have eaten at least 3 times
	counter := 0            // counter for how many phylos that have eaten at least 3 times
	var phylo int = 5       // number of phylo to create
	var prev chan bool
	first := make(chan bool)
	for i := 0; i < phylo; i++ {
		next := make(chan bool)
		if i < (phylo - 1) { // not the last
			if i == 0 { // the first
				go phylos(i, first, next, main)
			} else {
				go phylos(i, prev, next, main)
			}
			prev = make(chan bool)
			go fork(next, prev)
		} else {
			go phylos(i, prev, next, main)
			go fork(next, first)
		}
	}
	for {
		select {
		case <-main:
			counter++
			break
		default:
			break
		}
		if counter == phylo {
			break
		}
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

func phylos(id int, left chan bool, right chan bool, main chan bool) {
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
			if counter == 3 {
				main <- true
			}
			break
		}
		n := rand.Intn(1000)
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
}
