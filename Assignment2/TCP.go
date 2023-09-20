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
	"time"
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
	Await           = 3
)

type Message struct {
	Type    int // syn , synack , data, ack
	Seq     int //
	Ack     int
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
	var seq int
	var ack int

	for {
		switch status {
		case Closed:
			seq = rand.Intn(100)
			msg := Message{
				Type: Syn,
				Seq:  seq,
			}
			channel <- msg
			status = Await
			break
		case Await:
			select {
			case msg := <-channel:
				if msg.Type == SynAck && msg.Ack == seq+1 {
					seq++
					ack = msg.Ack
					msg := Message{}
				}
				break
			case <-time.After(2 * time.Second):
				status = Closed
				break
			}
		}
	}

	// send syn with random sequence

	// send syn

	for {
		select {
		case message := <-channel:
			fmt.Printf("Type %d Seq %d Payload %s", message.Type, message.Seq, message.Payload)
		default:
		}
	}

}
