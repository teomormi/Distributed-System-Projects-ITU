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

	network := make(chan Message, 2) // need to be buffered or some reply will be lost
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
				if msg.Type == Syn { // ok
					// reply with synack
					ack = msg.Seq + 1
					seq = rand.Intn(100)
					synack := Message{
						Type: SynAck,
						Seq:  seq,
						Ack:  ack,
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
					if msg.Type == Ack && msg.Ack == seq && msg.Seq == ack+1 { // estabilished connection
						status = Established
						// increase ack for check the first data packet
						ack++
						fmt.Println("Server status: Established")
					}
					break
				case <-time.After(5 * time.Second):
					status = Closed
					break
				}
				break
			}
		case Established:
			{
				// waiting for incoming data packets
				msg := <-channel
				if msg.Type == Data && msg.Ack == seq && msg.Seq == ack+1 { // check seq and ack number
					fmt.Printf("Server: Message payload: %s", msg.Payload)
					// test some delay or losses (retransit timeout 5 second - line 221)
					time.Sleep((time.Duration(rand.Intn(6)) + 2) * time.Second)
					// reply with ack
					ack = msg.Seq + 1
					seq = msg.Ack
					fmt.Printf("Server: Reply with ack to client (Seq=%d Ack=%d)\n", seq, ack)
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
					seq = msg.Ack
					finack := Message{
						Type: FinAck,
						Ack:  ack,
						Seq:  seq,
					}
					channel <- finack
					status = CloseWait
				}
				break
			}
			// check, probabluy wrong
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
	var seq, ack int
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
					ack = msg.Seq
					seq = msg.Ack + 1
					msg := Message{
						Type: Ack,
						Ack:  ack,
						Seq:  seq,
					}
					channel <- msg
					status = Established
					fmt.Println("Client Status: Established")
					time.Sleep(2 * time.Second) // rilegge il suo messaggio
					// first data package
					seq++
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
					ack = msg.Seq
					seq = msg.Ack + 1
					if len(list) != 0 {
						fmt.Printf("Client: Send next message to server (Seq=%d Ack=%d)\n", seq, ack)
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
				break
			case <-time.After(5 * time.Second): // resend or send the first packet
				if last.Type == -1 { // no retransmittion
					if len(list) > 0 { // at least one packet in the list
						fmt.Printf("Client: Send message to server (Seq=%d Ack=%d)\n", seq, ack)
						next := Message{
							Type:    Data,
							Ack:     ack,
							Seq:     seq,
							Payload: list[0],
						}
						list = list[1:] // remove first element from the list
						last = next
						channel <- next
					} else {
						fmt.Println("Nothing to send")
					}
				} else {
					// resend the packet
					fmt.Printf("Client: REsend message to server (Seq=%d Ack=%d)\n", last.Seq, last.Ack)
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
