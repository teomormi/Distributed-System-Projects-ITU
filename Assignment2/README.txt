a)
- bufio: to read the message payload from the command line
- fmt: we use it for print some messages on the command line
- math/rand: to generate the initial values for sequence and ack number
- os: indicate the standard input to the NewReader function
- strings: to isolate the string with the "exit" command
- time: trigger the retransmit of the message after some amount of time

b)
Our implementation uses threads. 
It is not realistic to use threads because they can not properly simulate the problems 
that can occur when communicating across a network. 

c)
Each message has two values ack and seq. 
With these two values it is possible to figure out which message was sent before another.
We drop a message in case of missing a previous one and we re-ack the already acked messages. 

d)
After sending a message the client waits for a certain amount of time.
If after the time has run out and the client has not recieved a message from the server acknowledging that it has recieved the message,
then the client will sent the message again. 

e)
The 3-way handshake is important to ensure that there exist a connection between the client and the server.
If the 3-way handshake is not performed then the client could be sending messages without knowing if the server will recieve them.
The sequences that gets established at the handshake also makes it easier to detect if messages gets lost or reordered.
