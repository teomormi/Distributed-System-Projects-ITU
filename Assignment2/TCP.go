/*
Authors:
 - Lucas Roy Guldbrandsen
 - Rafael Steffen Nguyen Jensen
 - Matteo Mormile

*/

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	Syn    int = 0
	SynAck     = 1
	Ack        = 2
	Fin        = 3
	FinAck     = 4
	Data       = 5
)

const (
	Closed      int = 0
	Listen          = 1
	Established     = 2
	Await           = 3
	CloseWait       = 4
	LastAck         = 5
)

type Message struct {
	Type    int // syn , synack , data, ack, fin, finack
	Seq     int //
	Ack     int
	Payload string // message content
}

var list []string

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
				// wait for the ack or go back to closed status (waiting for syn)
				case msg := <-channel:
					if msg.Type == Ack && msg.Ack == seq+1 && msg.Seq == ack { // estabilished connection
						status = Established
						// save last seq and ack to check the following messages
						seq = msg.Seq
						ack = msg.Ack
						fmt.Println("Server status: Established")
					}
					break
					/*case <-time.After(2 * time.Second):
					status = Closed
					break*/
				}
				break
			}
		case Established:
			{
				// waiting for incoming data packets
				msg := <-channel
				if msg.Type == Data && msg.Ack == ack && msg.Seq == seq+1 { // check seq and ack number
					fmt.Printf("Message payload: %s\n", msg.Payload)
					// test some delay or losses (retransit timeout 5 second - line 221)
					time.Sleep(time.Duration(rand.Intn(4)) * time.Millisecond)
					// reply with ack
					seq = msg.Ack
					ack = msg.Seq + 1
					fmt.Printf("Server: Send ack to client Seq=%d Ack=%d\n", seq, ack)
					reply := Message{
						Type: Ack,
						Ack:  ack,
						Seq:  seq,
					}
					// save last seq and ack to check the following messages
					channel <- reply
				}

				if msg.Type == Fin {
					ack = msg.Seq + 1
					finack := Message{
						Type: FinAck,
						Ack:  ack,
					}
					channel <- finack
					status = CloseWait
				}
				break
			}
		case CloseWait:
			{
				seq = rand.Intn(100)
				fin := Message{
					Type: Fin,
					Seq:  seq,
				}
				channel <- fin
				status = LastAck
				break
			}
		case LastAck:
			{
				msg := <-channel
				if msg.Type == FinAck {
					if msg.Ack == seq+1 {
						status = Closed
					}
				}
			}
		default:
			{
				break
			}

		}

	}
}

func client(channel chan Message) {
	status := Closed
	var seq int
	var ack int
	last := Message{
		Type: -1,
	}
	fmt.Println("Client Status: Closed")

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
			fmt.Println("Client Status: Await")
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
					fmt.Println("Client Status: Established")
					// go routine for input
					go input()
				} else {
					status = Closed
					fmt.Println("Client Status: Closed")
				}
				break
			case <-time.After(2 * time.Second):
				status = Closed
				fmt.Println("Client Status: Closed")
				break
			}
			break
		case Established:
			select {
			case msg := <-channel:
				if msg.Type == Ack && msg.Ack == seq+1 && msg.Seq == ack {
					fmt.Println("Client: Received Ack")
					if len(list) != 0 {
						ack = msg.Seq
						seq = msg.Ack + 1
						fmt.Printf("Client: Send message to server Seq=%d Ack=%d\n", seq, ack)
						next := Message{
							Type:    Data,
							Ack:     ack,
							Seq:     seq,
							Payload: list[0],
						}
						list = list[1:]
						last = next
						channel <- next
					} else { // no more transimmision
						last.Type = -1
					}
				}
			case <-time.After(5 * time.Second): // resend or send the first packet
				if last.Type == -1 { // no retransmittion
					if len(list) > 0 { // at least one packet in the list
						seq++
						fmt.Printf("Client: Send message to server Seq=%d Ack=%d\n", seq, ack)
						next := Message{
							Type:    Data,
							Ack:     ack,
							Seq:     seq,
							Payload: list[0],
						}
						list = list[1:] // remove first element from the list
						last = next
						channel <- next
					}
				} else {
					// resend the packet
					fmt.Println("Client: REsend message to server")
					channel <- last
				}
			}
		}
	}
}

func input() {
	var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ = reader.ReadString('\n')

		if input == "exit" {
			break
		}

		list = append(list, input) // insert at end
	}
}
