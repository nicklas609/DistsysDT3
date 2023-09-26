// Why it doesn't deadlock:
// Each philosopher checks both forks on either side.
// If one of the forks are unavailable, the philosopher will wait 1 second before checking again
// If both of the forks are availble the philosopher picks up the forks, eats for 100ms, puts the forks down and then waits for 2 seconds before wanting to eat again.
// This offset means that there will be no overlap between the philosophers trying to eat, and that no philosopher can hoard the forks.
package main

import (
	"fmt"
	"time"
)

var ch1 = make(chan bool)
var ch2 = make(chan bool)
var ch3 = make(chan bool)
var ch4 = make(chan bool)
var ch5 = make(chan bool)

var running = true
var finished = 0

func fork(ch chan bool, index int) {
	ch <- false
	for running {
		var response = <-ch
		ch <- response
	}
}

func phil(fork1 chan bool, fork2 chan bool, index int) {
	var count = 0

	for count < 3 {
		var rightFork = <-fork1
		var leftFork = <-fork2

		if !rightFork && !leftFork {
			fork1 <- true
			fork2 <- true

			fmt.Println("Phil #", index, "has picked up his right fork")
			fmt.Println("Phil #", index, "has picked up his left fork")

			time.Sleep(100 * time.Millisecond)
			if <-fork1 {
				fork1 <- false
			}
			if <-fork2 {
				fork2 <- false
			}
			count++
			fmt.Println("Phil #", index, " has eaten! (", count, " times )")
			time.Sleep(time.Second * 2)
		} else {
			fork1 <- rightFork
			fork2 <- leftFork
			fmt.Println("Phil #", index, "is thinking... hrmm")
			time.Sleep(time.Second * 1)
		}
	}
	fmt.Println("Phil #", index, "has finished eating :)\n========")
	finished++
}

func main() {
	go fork(ch1, 1)
	go fork(ch2, 2)
	go fork(ch3, 3)
	go fork(ch4, 4)
	go fork(ch5, 5)

	go phil(ch5, ch1, 1)
	go phil(ch1, ch2, 2)
	go phil(ch2, ch3, 3)
	go phil(ch3, ch4, 4)
	go phil(ch4, ch5, 5)

	for finished < 5 {

	}
}
