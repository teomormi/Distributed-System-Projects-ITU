/*
Authors:
 - Lucas Roy Guldbrandsen
 - Rafael Steffen Nguyen Jensen
 - Matteo Mormile

Our code that prevents a dealock works like so:
If a phylo is in the thinking state it will first make a request to get the left fork
If the the left fork is available it wil then make a request to the right fork
If the right fork is not available then it will return the left fork ensuring that a case where every phylo only has one fork will not occur
After a random time, if he previously failed to get the forks he will try again
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	main := make(chan bool) // channel for the phylos to communicate with main if they have eaten at least 3 times
	counter := 0            // counter for how many phylos that have eaten at least 3 times
	var phylo int = 5       // number of phylo to create
	var prev chan bool
	first := make(chan bool)
	for i := 0; i < phylo; i++ { // cycle to create and assign the correct channel to the fork and phylo routines
		next := make(chan bool)
		if i < (phylo - 1) { // not the last phylo
			if i == 0 { // the first one
				go phylos(i, first, next, main)
			} else {
				go phylos(i, prev, next, main)
			}
			prev = make(chan bool)
			go fork(next, prev)
		} else { // the last phylo
			go phylos(i, prev, next, main)
			go fork(next, first)
		}
	}
	for {
		select {
		case <-main: // a philosopher ate at least three times
			counter++
			break
		default:
			break
		}
		if counter == phylo { // all phylos has eaten at least 3 times
			break
		}
	}
}

func fork(left chan bool, right chan bool) { // code executed by fork routines
	// left and right channels are the next philosophers
	isTaken := false // the fork is available

	for {
		select {
		case message := <-left:
			if message {
				if isTaken { // if the fork is taken it replies with false
					left <- false
				} else { // if receives true request and fork is not taken responds with true
					left <- true
					isTaken = true
				}
			} else { // if receives false, it makes available the fork
				isTaken = false
			}
		case message := <-right: // same as before for the right side
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
	counter := 0 // number of times the philosopher has eaten

	for {
		switch state {
		case "thinking": // If he is thinking, after a random time he starts to eat
			left <- true
			if <-left { // if he takes the left fork
				right <- true
				if <-right { // if he takes also the right one
					state = "eating" // he starts to eat
					counter++
					fmt.Printf("%d I'm eating for the %d time \n", id, counter)
				} else { // release the left fork if he can't take the left one
					left <- false
				}
			}
			break
		case "eating":
			// release all the forks
			left <- false
			right <- false
			state = "thinking"
			fmt.Printf("%d I'm thinking\n", id)
			if counter == 3 { // send message to the main
				main <- true
			}
			break
		}
		n := rand.Intn(1000)
		time.Sleep(time.Duration(n) * time.Millisecond)
	}
}
