package subpub

import (
	"context"
	"sync"
)

// Service stands for Publisher-Subscriber service
type Service struct {
	mu      sync.RWMutex
	storage map[string][]*subscriber
	Quit    chan struct{}
}

func (s *Service) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	sub := newSubscriber(cb)

	s.mu.Lock()
	s.storage[subject] = append(s.storage[subject], sub)
	s.mu.Unlock()

	go func() {
		for {
			select {
			case msg, ok := <-sub.message:
				if !ok {
					return
				}
				sub.handler(msg)
			}
		}
	}()

	subc := newSubscription(subject, s, sub)
	return subc, nil
}

func (s *Service) Publish(subject string, msg interface{}) error {
	s.mu.RLock()
	subjectSubs := s.storage[subject]
	s.mu.RUnlock()

	for _, sub := range subjectSubs {
		select {
		case sub.message <- msg:
		default:
		}
	}

	return nil
}

func (s *Service) Close(ctx context.Context) error {
	close(s.Quit)

	done := make(chan struct{})
	go func() {
		s.mu.Lock()
		for _, subs := range s.storage {
			for _, sub := range subs {
				close(sub.message)
			}
		}
		s.mu.Unlock()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
