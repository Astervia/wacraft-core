package synch_redis

import (
	"testing"
	"time"
)

func TestRedisPubSub_CrossInstance(t *testing.T) {
	client := testRedisClient(t)

	// Two PubSub instances simulating two app instances
	psA := NewRedisPubSub(client)
	psB := NewRedisPubSub(client)

	sub, err := psA.Subscribe("events")
	if err != nil {
		t.Fatalf("Subscribe() error: %v", err)
	}
	defer sub.Unsubscribe()

	// Small delay to ensure subscription is active
	time.Sleep(50 * time.Millisecond)

	// Publish from instance B
	if err := psB.Publish("events", []byte("cross-instance")); err != nil {
		t.Fatalf("Publish() error: %v", err)
	}

	select {
	case msg := <-sub.Channel():
		if string(msg) != "cross-instance" {
			t.Errorf("received %q, want %q", msg, "cross-instance")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: subscriber on instance A should receive message from instance B")
	}
}

func TestRedisPubSub_MultipleChannels(t *testing.T) {
	client := testRedisClient(t)
	ps := NewRedisPubSub(client)

	sub1, err := ps.Subscribe("ch1")
	if err != nil {
		t.Fatalf("Subscribe(ch1) error: %v", err)
	}
	defer sub1.Unsubscribe()

	sub2, err := ps.Subscribe("ch2")
	if err != nil {
		t.Fatalf("Subscribe(ch2) error: %v", err)
	}
	defer sub2.Unsubscribe()

	time.Sleep(50 * time.Millisecond)

	ps.Publish("ch1", []byte("msg-ch1"))

	select {
	case msg := <-sub1.Channel():
		if string(msg) != "msg-ch1" {
			t.Errorf("ch1: got %q, want %q", msg, "msg-ch1")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ch1 subscriber should receive message")
	}

	select {
	case <-sub2.Channel():
		t.Fatal("ch2 subscriber should NOT receive ch1 message")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}

func TestRedisPubSub_Unsubscribe(t *testing.T) {
	client := testRedisClient(t)
	ps := NewRedisPubSub(client)

	sub, err := ps.Subscribe("unsub-test")
	if err != nil {
		t.Fatalf("Subscribe() error: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := sub.Unsubscribe(); err != nil {
		t.Fatalf("Unsubscribe() error: %v", err)
	}

	// Publish after unsubscribe — should not panic or error
	if err := ps.Publish("unsub-test", []byte("after")); err != nil {
		t.Fatalf("Publish() after Unsubscribe error: %v", err)
	}
}

func TestRedisPubSub_MultipleSubscribersSameChannel(t *testing.T) {
	client := testRedisClient(t)

	psA := NewRedisPubSub(client)
	psB := NewRedisPubSub(client)

	subA, _ := psA.Subscribe("shared")
	defer subA.Unsubscribe()

	subB, _ := psB.Subscribe("shared")
	defer subB.Unsubscribe()

	time.Sleep(50 * time.Millisecond)

	psA.Publish("shared", []byte("broadcast"))

	for name, sub := range map[string]interface{ Channel() <-chan []byte }{
		"A": subA,
		"B": subB,
	} {
		select {
		case msg := <-sub.Channel():
			if string(msg) != "broadcast" {
				t.Errorf("subscriber %s: got %q, want %q", name, msg, "broadcast")
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("subscriber %s: timeout", name)
		}
	}
}
