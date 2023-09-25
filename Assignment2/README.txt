a)

b)
Our implementation uses threads. 
It is not realistic to use threads because they can not properly simulate the problems 
that can occur when communicating across a network. 

c)
Each message has two values ack and seq. 
With these two values it is possible to figure out which message was sent before another.
Then the messages can be reordered based on the values. 

d)
After sending a message the client waits for a certain amount of time.
If after the time has run out and the client has not recieved a message from the server acknowledging that it has recieved the message,
then the client will sent the message again. 

e)
The 3-way handshake is important to ensure that there exist a connection between the client and the server.
If the 3-way handshake is not performed then the client could be sending messages without knowing if the server will recieve them.
The sequences that gets established at the handshake also makes it easier to detect if messages gets lost or reordered.
