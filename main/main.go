// PubSub implementation taking a channel as argument in Subscibe, w/o Close.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"fmt"
	"sync"
	"time"
)

type Pubsub struct {
	mu   sync.RWMutex // what is this?
	subs map[string][]chan string
}

func NewPubsub() *Pubsub {
	ps := &Pubsub{}
	// what is this?
	ps.subs = make(map[string][]chan string)
	// this is a map with key as strings(different topics) to a slice of
	// channels of type string
	return ps
}

func (ps *Pubsub) Subscribe(topic string, ch chan string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subs[topic] = append(ps.subs[topic], ch)
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.subs[topic] {
		ch <- msg
	}
}

func main() {
	ps := NewPubsub()

	// makesub creates a new channel, subscribes it to the given topic, and
	// returns it. The channel is buffered so that the subscription will not
	// block.
	makesub := func(topic string) chan string {
		ch := make(chan string, 1)
		ps.Subscribe(topic, ch)
		return ch
	}
	ch1 := makesub("tech")
	ch2 := makesub("travel")
	ch3 := makesub("travel")

	// listener is a goroutine that listens on the given channel and prints
	// everything it receives from it.
	listener := func(name string, ch chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}
	}

	go listener("1", ch1)
	go listener("2", ch2)
	go listener("3", ch3)

	// pub publishes some messages to topics.
	pub := func(topic string, msg string) {
		fmt.Printf("Publishing @%s: %s\n", topic, msg)
		ps.Publish(topic, msg)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)
	pub("tech", "tablets")
	pub("health", "vitamins")
	pub("tech", "robots")
	pub("travel", "beaches")
	pub("travel", "hiking")
	pub("tech", "drones")

	time.Sleep(50 * time.Millisecond)
}
