package main

func main() {
	var phylo int = 5 // number of phylo to create
	var prev chan bool
	first := make(chan bool)
	for i := 0; i < phylo; i++ {
		next := make(chan bool)
		if i < (phylo - 1) { // not the last
			if i == 0 { // the first
				go phylos(first, next)
			} else {
				go phylos(prev, next)
			}
			prev = make(chan bool)
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

func phylos(left chan bool, right chan bool) {
	state := "thinking"

	for {
		switch state {
		case "thinking":
			if <-left && <-right {
				left <- true
				right <- true
				state = "eating"
			}
			break
		case "eating":
			left <- false
			right <- false
			state = "thinking"
			break
		}
	}
}
