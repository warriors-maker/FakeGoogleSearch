package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Search func(query string)

type Result struct {
	s string
}

// Return a closure function
func fakeSearch(kind string, ch chan string) Search {
	return func(query string) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		ch <- fmt.Sprintf("%s result for %q\n", kind, query)
	}
}

// // Sequential Search
// func SeqGoogle(query string) []string {
// 	var s []string
// 	s = append(s, Web(query).s)
// 	s = append(s, Image(query).s)
// 	s = append(s, Video(query).s)
// 	return s
// }

// Parallel Search
func ParallelGoogle(query string) []string {
	var s []string
	ch := make(chan string)
	var (
		Web   = fakeSearch("Web", ch)
		Image = fakeSearch("Image", ch)
		Video = fakeSearch("Video", ch)
	)
	go Web(query)
	go Image(query)
	go Video(query)
	// Wait with no time limit on it
	for i := 0; i < 3; i++ {
		s = append(s, <-ch)
	}
	return s
}

// Parallel Search but with TimeLimits
func PTimeLimitGoogle(query string) []string {
	var s []string
	ch := make(chan string)
	var (
		Web   = fakeSearch("Web", ch)
		Image = fakeSearch("Image", ch)
		Video = fakeSearch("Video", ch)
	)

	time_ch := time.After(50 * time.Millisecond)

	go Web(query)
	go Image(query)
	go Video(query)

	for {
		select {
		case v := <-ch:
			s = append(s, v)
		case <-time_ch:
			fmt.Println("Sorry, time is up!")
			return s
		}
	}
}

// TODO:
//Need to do the replication part
// Individual server needs to do go Web, go Video, go Image
// A channel to communicate between each server with the main
// if main gets a result return that because it is the most efficient one
func ParallelGoogleReplica(query string, replicaID int) chan []string {
	quit := make(chan []string)
	go func() {
		var s []string
		ch := make(chan string)
		var (
			Web   = fakeSearch("Web", ch)
			Image = fakeSearch("Image", ch)
			Video = fakeSearch("Video", ch)
		)
		go Web(query)
		go Image(query)
		go Video(query)
		for i := 0; i < 3; i++ {
			s = append(s, <-ch)
		}
		fmt.Printf("Finished in replica %d\n", replicaID)
		quit <- s
	}()
	return quit
}

// fanIn pattern
func clientListening(replica1, replica2, replica3 chan []string) {
	// Only listen for the first value that comes out
	select {
	case value1 := <-replica1:
		fmt.Println(value1)
	case value2 := <-replica2:
		fmt.Println(value2)
	case value3 := <-replica3:
		fmt.Println(value3)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	query := "golang"
	clientListening(
		ParallelGoogleReplica(query, 1),
		ParallelGoogleReplica(query, 2),
		ParallelGoogleReplica(query, 3),
	)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
