package synch_service

import (
	"sync"

	synch_contract "github.com/Astervia/wacraft-core/src/synch/contract"
)

type MemoryPubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]*memorySubscription
}

type memorySubscription struct {
	ch      chan []byte
	channel string
	pubsub  *MemoryPubSub
	closed  bool
	mu      sync.Mutex
}

func NewMemoryPubSub() *MemoryPubSub {
	return &MemoryPubSub{
		subscribers: make(map[string][]*memorySubscription),
	}
}

func (p *MemoryPubSub) Publish(channel string, message []byte) error {
	p.mu.RLock()
	original := p.subscribers[channel]
	subs := make([]*memorySubscription, len(original))
	copy(subs, original)
	p.mu.RUnlock()

	for _, sub := range subs {
		sub.mu.Lock()
		if !sub.closed {
			select {
			case sub.ch <- message:
			default:
				// Drop message if subscriber buffer is full
			}
		}
		sub.mu.Unlock()
	}

	return nil
}

func (p *MemoryPubSub) Subscribe(channel string) (synch_contract.Subscription, error) {
	sub := &memorySubscription{
		ch:      make(chan []byte, 100),
		channel: channel,
		pubsub:  p,
	}

	p.mu.Lock()
	p.subscribers[channel] = append(p.subscribers[channel], sub)
	p.mu.Unlock()

	return sub, nil
}

func (p *MemoryPubSub) removeSubscription(sub *memorySubscription) {
	p.mu.Lock()
	defer p.mu.Unlock()

	subs := p.subscribers[sub.channel]
	for i, s := range subs {
		if s == sub {
			p.subscribers[sub.channel] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	if len(p.subscribers[sub.channel]) == 0 {
		delete(p.subscribers, sub.channel)
	}
}

func (s *memorySubscription) Channel() <-chan []byte {
	return s.ch
}

func (s *memorySubscription) Unsubscribe() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true
	s.pubsub.removeSubscription(s)
	close(s.ch)
	return nil
}
