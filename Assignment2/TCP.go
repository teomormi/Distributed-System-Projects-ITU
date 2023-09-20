/*
Authors:
 - Lucas Roy Guldbrandsen
 - Rafael Steffen Nguyen Jensen
 - Matteo Mormile

*/

package main

import (
	"fmt"
	"math/rand"
)

const (
	Syn    int = 0
	SynAck     = 1
	Ack        = 2
	Data       = 3
)

const (
	Closed      int = 0
	Listen          = 1
	Established     = 2
)

type Message struct {
	Type    int    // syn , synack , data, ack
	Seq     int    //
	Payload string // message content
}

func main() {

	network := make(chan Message)
	// 2 go routine simulate server and client
	go client(network)
	go server(network)
	for {
		// do nothing
	}
}

func server(channel chan Message) {
	//var msg Message
	//msg.Seq = 1
	//msg.Payload = "Hello"
	//msg.Type = 1
	status := Closed
	//channel <- msg
}

func client(channel chan Message) {
	status := Closed
	// send syn with random sequence
	var msg Message
	msg.Type = Syn
	msg.Seq = rand.Intn(100)
	// send syn
	channel <- msg
	for {
		select {
		case message := <-channel:
			fmt.Printf("Type %d Seq %d Payload %s", message.Type, message.Seq, message.Payload)
		default:
		}
	}

}
