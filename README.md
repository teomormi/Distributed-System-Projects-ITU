# Distributed System Hand-ins
The repository contains the 6 mandatory hand-ins for the fall course of Distributed Systems at ITU. \
The course covers a number of topics, including: 
  - Communication over the network leads to dropped and reordered packets which requires robust request/reply and multicast protocols. 
  - Since programs communicate over an open and potentially faulty network, it is difficult for them to synchronise on time, leaders, and data. Protocols for time synchronisation, logical time, consensus, and data distribution address such issues. 
  - In order to support changes and evolution of it-systems, loosely coupled architectures such as Service Oriented Architecture (e.g., Microservices) are employed, as well as process-oriented architectures for combining services in process flows in an adaptable way. 


# Authors
- Lucas Roy Guldbrandsen
- Rafael Steffen Nguyen Jensen
- Matteo Mormile

## First mandatory Hand-in
The Dining Philosophers is a well-known problem in Computer Science that concerns concurrency. At a dining round table, there are five philosophers who are supposed to have dinner. Philosophers are kind of special and while they have dinner, they either *eat* their food or *think* about something. In order to be able to eat, they must get hold of two forks (the food is very special and cannot be handled with one fork). Unfortunately, there are only five forks at the table, each of them uniquely placed between two philosophers (the table is round, there is exactly one fork between any two philosophers -- each philosopher can only reach the two forks that are nearby). As a consequence, it is never the case that all philosophers can eat at the same time (max two at a time).  Eating is not limited by food quantity or stomach space (which are both assumed to be infinite). This problem is interesting because, depending on how they decide to pick the forks, the philosopher may reach a deadlock.
The goal of this project is to implement the dining philosophers problem in Go, with the following requirements:

- Each fork must have its own thread (goroutine) 
- Each philosopher must have its own thread (goroutine) 
- Philosophers and forks must communicate with each other *only* by  using channels 
- The system must be designed in a way that does not lead to a deadlock  (and each philosopher must eat at least 3 times).  Comment in the code why the system does not deadlock. 
- A sequentialisation of the system (executing only one philosopher at a time) is not acceptable. I.e., philosophers must be able to request a fork at any time. 
- Philosophers must display (print on screen) any state change (eating or thinking) during their execution.

## Second mandatory Hand-in
Implement the TCP/IP protocol in Go using threads. 
Attach to your submission, a *README* file answering the following questions:

a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
b) Does your implementation use threads or processes? Why is it not realistic to use threads?
c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
d) In case messages can be delayed or lost, how does your implementation handle message loss?
e) Why is the 3-way handshake important?

## Third mandatory Hand-in
You have to implement Chitty-Chat a distributed system, that is providing a chatting service, and keeps track of logical time using Lamport Timestamps.
We call clients of the Chitty-Chat service Participants. 

### System Requirements
    R1: Chitty-Chat is a distributed service, that enables its clients to chat. The service is using gRPC for communication. You have to design the API, including gRPC methods and data types. 
    R2: Clients in Chitty-Chat can Publish a valid chat message at any time they wish.  A valid message is a string of UTF-8 encoded text with a maximum length of 128 characters. A client publishes a message by making a gRPC call to Chitty-Chat.
    R3: The Chitty-Chat service has to broadcast every published message, together with the current logical timestamp, to all participants in the system, by using gRPC. It is an implementation decision left to the students, whether a Vector Clock or a Lamport timestamp is sent.
    R4: When a client receives a broadcasted message, it has to write the message and the current logical timestamp to the log
    R5: Chat clients can join at any time. 
    R6: A "Participant X  joined Chitty-Chat at Lamport time L" message is broadcast to all Participants when client X joins, including the new Participant.
    R7: Chat clients can drop out at any time. 
    R8: A "Participant X left Chitty-Chat at Lamport time L" message is broadcast to all remaining Participants when Participant X leaves.

 ### Technical Requirements:
    Use gRPC for all messages passing between nodes
    Use Golang to implement the service and clients
    Every client has to be deployed as a separate process
    Log all service calls (Publish, Broadcast, ...) using the log package
    Demonstrate that the system can be started with at least 3 client nodes 
    Demonstrate that a client node can join the system
    Demonstrate that a client node can leave the system
    Optional: All elements of the Chitty-Chat service are deployed as Docker containers

## Fourth mandatory Hand-in
You have to implement distributed mutual exclusion between nodes in your distributed system. 
Your system has to consist of a set of peer nodes, and you are not allowed to base your implementation on a central server solution.
You can decide to base your implementation on one of the algorithms, that were discussed in lecture 7.

### System Requirements:
```
R1: Implement a system with a set of peer nodes, and a Critical Section, that represents a sensitive system operation. Any node can at any time decide it wants access to the Critical Section. Critical section in this exercise is emulated, for example by a print statement, or writing to a shared file on the network.
R2: Safety: Only one node at the same time is allowed to enter the Critical Section 
R3: Liveliness: Every node that requests access to the Critical Section, will get access to the Critical Section (at some point in time)
```

### Technical Requirements:
    Use Golang to implement the service's nodes
    In you source code repo, provide a README.md, that explains how to start your system
    Use gRPC for message passing between nodes
    Your nodes need to find each other. This is called service discovery. You could consider  one of the following options for implementing service discovery:
        Supply a file with IP addresses/ports of other nodes
        Enter IP address/ports through the command line
        use an existing package or service
    Demonstrate that the system can be started with at least 3 nodes
    Demonstrate using your system's logs,  a sequence of messages in the system, that leads to a node getting access to the Critical Section. You should provide a discussion of your algorithm, using examples from your logs.

## Fifth mandatory Hand-in
You must implement a **distributed auction system** using replication: a distributed component which handles auctions, and provides operations for bidding and querying the state of an auction. The component must faithfully implement the semantics of the system described below, and must at least be resilient to one (1) crash failure.
The goal of this mandatory activity is that you learn (by doing) how to use replication to design a service that is resilent to crashes. In particular, it is important that you can recognise what the key issues that may arise are and understand how to deal with them.

## Sixth mandatory Hand-in
This last mandatory activity is about Raft and it requires each group to find an implementation of Raft (in Go) and write a 2-page report about it. 

### What to do? 
Find an online implementation of Raft
It must be in Go: if you find something in another language, you need to translate it yourself or implement it from scratch write a max 2-page report on the implementation you have found/written focussing on:
- Go-specific features (Like Goroutines, channels and so on
- communication method
- properties guaranteed by the protocol implementation (is it faithful to what we saw at the lecture?)
