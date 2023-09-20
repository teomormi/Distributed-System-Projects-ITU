/*
Authors:
 - Lucas Roy Guldbrandsen
 - Rafael Steffen Nguyen Jensen
 - Matteo Mormile

*/

package main

import (
	"fmt"
)

type Message struct {
	Type    int
	Seq     int
	Payload string
}

func main() {

	network := make(chan Message)
	go client(network)
	go server(network)
	for {
	}
	//fmt.Printf("Hello world!")
	// 2 go routine simulate server and client
}

func server(channel chan Message) {
	var msg Message
	msg.Seq = 1
	msg.Payload = "Hello"
	msg.Type = 1
	channel <- msg
}

func client(channel chan Message) {
	answer := <-channel
	fmt.Printf("Type %d Seq %d Payload %s", answer.Type, answer.Seq, answer.Payload)
}
