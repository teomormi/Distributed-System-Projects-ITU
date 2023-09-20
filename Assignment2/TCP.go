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
	status := Closed
	var seq, ack int
	for {
		switch status {
		case Closed:
			{
				msg := <-channel
				if msg.Type == Syn {
					// reply with synack
					ack = msg.Seq + 1
					seq = rand.Intn(100)
					synack := Message{
						Type:    SynAck,
						Seq:     seq,
						Ack:     ack,
						Payload: "test",
					}
					channel <- synack
					status = Listen
					fmt.Println("Server status: Listen")
				}
				break
			}
		case Listen:
			{
				select {
				case msg := <-channel:
					if msg.Type == Ack && msg.Ack == seq+1 && msg.Seq == ack { // estabilished connection
						status = Established
						fmt.Println("Server status: Established")
					}
					break
				case <-time.After(2 * time.Second):
					status = Closed
					break
				}
				// wait for the ack or go back to closed status (waiting for syn)
				break
			}
		case Established:
			{
				// waiting for incoming data packets
			}
		default:
			break
		}
	}
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
			fmt.Println("Client: Status changed to Await")
			break
		case Await:
			select {
			case msg := <-channel:
				if msg.Type == SynAck && msg.Ack == seq+1 {
					seq++
					ack = msg.Seq + 1
					msg := Message{
						Type: Ack,
						Ack:  ack,
						Seq:  seq,
					}
					channel <- msg
					status = Established
					fmt.Println("Client: Status changed to Established")
				} else {
					status = Closed
					fmt.Println("Client: Status changed to Closed")
				}
				break
			case <-time.After(2 * time.Second):
				status = Closed
				fmt.Println("Client: Status changed to Closed")
				break
			}
			break
		case Established:
			fmt.Println("Client: Connections established")
			for {

			}
		}
	}
}
