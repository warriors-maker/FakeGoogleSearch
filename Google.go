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

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	result := PTimeLimitGoogle("golang")
	elapsed := time.Since(start)
	fmt.Println(result)
	fmt.Println(elapsed)
}
