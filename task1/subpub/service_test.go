package subpub

import (
	"context"
	"sync"
	"testing"
)

const (
	n       = 100
	subject = "weather"
)

func newService() *Service {
	return &Service{
		mu:      sync.RWMutex{},
		storage: make(map[string][]*subscriber),
		Quit:    make(chan struct{}),
	}
}

func TestService_100Subscribe(t *testing.T) {
	service := newService()
	cb := func(msg interface{}) {}

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_, err := service.Subscribe(subject, cb)
			if err != nil {
				t.Fail()
			}
		}()
	}
	wg.Wait()

	wait, got := n, len(service.storage[subject])
	if wait != got {
		t.Errorf("expected %d subscribers, got %d", wait, got)
	}
}

func TestService_Publish(t *testing.T) {
	service := newService()
	counter := 0
	var mu sync.Mutex
	var received sync.WaitGroup
	received.Add(n)
	cb := func(msg interface{}) {
		mu.Lock()
		counter++
		mu.Unlock()
		received.Done()
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_, err := service.Subscribe(subject, cb)
			if err != nil {
				t.Fail()
			}
		}()
	}
	wg.Wait()

	err := service.Publish(subject, "it's awful today")
	if err != nil {
		t.Fail()
	}
	received.Wait()

	wait, got := n, counter
	if wait != got {
		t.Errorf("expected %d messages, got %d", wait, got)
	}
}

func TestService_Close_ClosesAllSubscriberChannels(t *testing.T) {
	service := newService()
	for i := 0; i < n; i++ {
		sub := &subscriber{message: make(chan interface{}, 1)}
		service.storage[subject] = append(service.storage[subject], sub)
	}

	ctx := context.Background()
	if err := service.Close(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, sub := range service.storage[subject] {
		select {
		case _, ok := <-sub.message:
			if ok {
				t.Errorf("subscriber #%d: expected closed channel, but it's still open", i)
			}
		default:
			t.Errorf("subscriber #%d: expected closed channel, but no read was possible", i)
		}
	}
}
