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
	"strings"
	"time"
)

const ( // message type
	Syn    int = 0
	SynAck     = 1
	Ack        = 2
	Fin        = 3
	FinAck     = 4
	Data       = 5
)

const ( // connection status
	Closed      int = 0
	Listen          = 1
	Established     = 2
	Await           = 3
	CloseWait       = 4
	LastAck         = 5
)

type Message struct {
	Type    int
	Seq     int
	Ack     int
	Payload string // message content
}

var list []string // buffer of messages' payload to be send

func main() {
	main := make(chan bool) // channel for end the main after the closing of connection
	counter := 0
	network := make(chan Message, 2) // need to be buffered or some reply will be lost (2 send at the same time)
	fmt.Println("After the establish of the connection it is possibile to send message from server to client\n" +
		"by only type in the terminal the payload of the message and press enter, type exit to quit")
	// 2 go routine simulate server and client
	go client(network, main)
	go server(network, main)
	for {
		select {
		case <-main: // a philosopher ate at least three times
			counter++
			break
		default:
			break
		}
		if counter == 2 { // client and server closed
			break
		}
	}
}

func server(channel chan Message, main chan bool) {
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
					if msg.Type == Ack && msg.Ack == seq+1 && msg.Seq == ack { // estabilished connection
						status = Established
						// increase ack number for checking the first data packet
						ack++
						fmt.Println("Server status: Established")
					}
					break
				// closed if not receive ack after synack
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
				if msg.Ack <= seq+1 && msg.Seq <= ack { // '<' for missing ack packet already sent
					if msg.Type == Data {
						fmt.Printf("Server: Received data message with payload: %s", msg.Payload)
						// test some delay or losses (retransit timeout 5 second)
						time.Sleep((time.Duration(rand.Intn(8)) + 1) * time.Second)
						// reply with ack
						ack = msg.Seq + 1
						seq = msg.Ack
						fmt.Printf("Server: Reply with Ack to client (Seq=%d Ack=%d)\n", seq, ack)
						reply := Message{
							Type: Ack,
							Ack:  ack,
							Seq:  seq,
						}
						channel <- reply
					}
					if msg.Type == Fin {
						ack = msg.Seq + 1
						seq = msg.Ack
						fmt.Printf("Server: Reply with FinAck to client (Seq=%d Ack=%d)\n", seq, ack)
						finack := Message{
							Type: FinAck,
							Ack:  ack,
							Seq:  seq,
						}
						channel <- finack
						status = CloseWait
					}
				}
				break
			}
		case CloseWait:
			{
				fmt.Printf("Server: Send Fin to client (Seq=%d Ack=%d)\n", seq, ack)
				fin := Message{
					Type: Fin,
					Seq:  seq,
					Ack:  ack,
				}
				channel <- fin
				status = LastAck
				break
			}
		case LastAck:
			{
				msg := <-channel
				if msg.Type == FinAck && msg.Ack == seq+1 {
					fmt.Println("Server Status: Closed")
					main <- true
					status = Closed
				}
			}
		default:
			{
				break
			}
		}
	}
}

func client(channel chan Message, main chan bool) {
	status := Closed
	var seq, ack int
	var retransmit bool = true
	// var to store the last message sent (retransmit it if needed)
	last := Message{
		Type: -1,
	}

	fmt.Println("Client Status: Closed")
	for {
		switch status {
		case Closed:
			// send syn message
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
					ack = msg.Seq + 1
					seq = msg.Ack
					msg := Message{
						Type: Ack,
						Ack:  ack,
						Seq:  seq,
					}
					channel <- msg
					status = Established
					fmt.Println("Client Status: Established")
					// first data message
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
				// receive ack for a data message sent
				if msg.Type == Ack && msg.Ack == seq+1 && msg.Seq == ack {
					ack = msg.Seq + 1
					seq = msg.Ack
					if len(list) != 0 {
						// send next message
						last = sendMessage(seq, ack)
						if last.Type == Fin {
							retransmit = false
						}
						channel <- last
					} else { // no more transimmision
						last.Type = -1
					}
				}
				if msg.Type == Fin { // close connection, reply with finack
					ack = msg.Seq + 1
					seq = msg.Ack
					fmt.Printf("Client: Reply with FinAck to server (Seq=%d Ack=%d)\n", seq, ack)
					msg := Message{
						Type: FinAck,
						Ack:  ack,
						Seq:  seq,
					}
					channel <- msg
					fmt.Println("Client Status: Closed")
					main <- true
					status = Closed
				}
				break
			case <-time.After(5 * time.Second): // resend last packet or send the new one
				if last.Type == -1 { // no retransmittion
					if len(list) > 0 { // a new message in the list
						last = sendMessage(seq, ack)
						channel <- last
					} else {
						fmt.Println("Nothing to send")
					}
				} else {
					// resend the packet
					if retransmit {
						fmt.Printf("Client: REsend message to server (Seq=%d Ack=%d)\n", last.Seq, last.Ack)
						channel <- last
					}
				}
			}
		}
	}
}

func sendMessage(seq int, ack int) Message {
	var msg_type int
	if strings.Contains(list[0], "exit") {
		fmt.Printf("Client: Send Fin to server (Seq=%d Ack=%d)\n", seq, ack)
		msg_type = Fin
	} else {
		fmt.Printf("Client: Send Data to server (Seq=%d Ack=%d)\n", seq, ack)
		msg_type = Data
	}
	msg := Message{
		Type:    msg_type,
		Ack:     ack,
		Seq:     seq,
		Payload: list[0],
	}
	list = list[1:] // remove first element from the list
	return msg
}

func input() {
	var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ = reader.ReadString('\n')
		list = append(list, input) // insert at end
		if input == "exit" {
			break
		}
	}
}
