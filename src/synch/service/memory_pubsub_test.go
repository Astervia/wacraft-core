package synch_service

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryPubSub_PublishSubscribe(t *testing.T) {
	ps := NewMemoryPubSub()

	sub, err := ps.Subscribe("ch1")
	if err != nil {
		t.Fatalf("Subscribe() error: %v", err)
	}
	defer sub.Unsubscribe()

	msg := []byte("hello")
	if err := ps.Publish("ch1", msg); err != nil {
		t.Fatalf("Publish() error: %v", err)
	}

	select {
	case got := <-sub.Channel():
		if string(got) != "hello" {
			t.Errorf("received %q, want %q", got, "hello")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestMemoryPubSub_MultipleSubscribers(t *testing.T) {
	ps := NewMemoryPubSub()

	subs := make([]interface {
		Channel() <-chan []byte
		Unsubscribe() error
	}, 3)

	for i := range subs {
		sub, err := ps.Subscribe("ch1")
		if err != nil {
			t.Fatalf("Subscribe() error: %v", err)
		}
		defer sub.Unsubscribe()
		subs[i] = sub
	}

	ps.Publish("ch1", []byte("broadcast"))

	for i, sub := range subs {
		select {
		case got := <-sub.Channel():
			if string(got) != "broadcast" {
				t.Errorf("subscriber %d: got %q, want %q", i, got, "broadcast")
			}
		case <-time.After(time.Second):
			t.Fatalf("subscriber %d: timeout", i)
		}
	}
}

func TestMemoryPubSub_Unsubscribe(t *testing.T) {
	ps := NewMemoryPubSub()

	sub, _ := ps.Subscribe("ch1")
	sub.Unsubscribe()

	// Channel should be closed
	_, ok := <-sub.Channel()
	if ok {
		t.Fatal("channel should be closed after Unsubscribe")
	}

	// Publish should not panic
	if err := ps.Publish("ch1", []byte("after-unsub")); err != nil {
		t.Fatalf("Publish() after Unsubscribe error: %v", err)
	}
}

func TestMemoryPubSub_DoubleUnsubscribe(t *testing.T) {
	ps := NewMemoryPubSub()

	sub, _ := ps.Subscribe("ch1")
	sub.Unsubscribe()

	// Second unsubscribe should not panic
	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("second Unsubscribe() error: %v", err)
	}
}

func TestMemoryPubSub_IsolatedChannels(t *testing.T) {
	ps := NewMemoryPubSub()

	sub1, _ := ps.Subscribe("ch1")
	defer sub1.Unsubscribe()

	sub2, _ := ps.Subscribe("ch2")
	defer sub2.Unsubscribe()

	ps.Publish("ch2", []byte("only-ch2"))

	select {
	case <-sub1.Channel():
		t.Fatal("ch1 subscriber should not receive ch2 messages")
	case <-time.After(50 * time.Millisecond):
		// Expected
	}

	select {
	case got := <-sub2.Channel():
		if string(got) != "only-ch2" {
			t.Errorf("got %q, want %q", got, "only-ch2")
		}
	case <-time.After(time.Second):
		t.Fatal("ch2 subscriber should have received message")
	}
}

func TestMemoryPubSub_BufferFull(t *testing.T) {
	ps := NewMemoryPubSub()

	sub, _ := ps.Subscribe("ch1")
	defer sub.Unsubscribe()

	// Fill the buffer (capacity is 100)
	for i := 0; i < 100; i++ {
		ps.Publish("ch1", []byte("fill"))
	}

	// This should not panic or block — message is dropped
	if err := ps.Publish("ch1", []byte("overflow")); err != nil {
		t.Fatalf("Publish() on full buffer error: %v", err)
	}
}

func TestMemoryPubSub_ConcurrentPublishSubscribe(t *testing.T) {
	ps := NewMemoryPubSub()
	var wg sync.WaitGroup

	// Concurrent subscribers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub, _ := ps.Subscribe("ch1")
			time.Sleep(10 * time.Millisecond)
			sub.Unsubscribe()
		}()
	}

	// Concurrent publishers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ps.Publish("ch1", []byte("msg"))
		}()
	}

	wg.Wait()
}
