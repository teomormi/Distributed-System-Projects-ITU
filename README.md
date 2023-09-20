# Distributed System Project

# Authors
- Lucas Roy Guldbrandsen
- Rafael Steffen Nguyen Jensen
- Matteo Mormile

## First mandatory Hand-in
The goal of this project is to implement the dining philosophers problem in Go, with the following requirements:

- Each fork must have its own thread (goroutine)

- Each philosopher must have its own thread (goroutine)

- Philosophers and forks must communicate with each other *only* by  using channels

- the system must be designed in a way that does not lead to a deadlock  (and each philosopher must eat at least 3 times).  Comment in the code why the system does not deadlock.

- A sequentialisation of the system (executing only one philosopher at a time) is not acceptable. I.e., philosophers must be able to request a fork at any time.

- Philosophers must display (print on screen) any state change (eating or thinking) during their execution.

## Second mandatory Hand-in
