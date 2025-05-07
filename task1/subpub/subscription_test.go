package subpub

import "testing"

func TestSubscription_Unsubscribe(t *testing.T) {
	service := newService()
	sub := &subscriber{message: make(chan interface{}, 1)}
	service.storage[subject] = []*subscriber{sub}
	subc := &subscription{
		subject: subject,
		srv:     service,
		sub:     sub,
	}

	if len(service.storage[subject]) != 1 {
		t.Fatalf("expected 1 subscriber before unsubscribe, got %d", len(service.storage[subject]))
	}

	subc.Unsubscribe()
	if len(service.storage[subject]) != 0 {
		t.Errorf("expected 0 subscribers after unsubscribe, got %d", len(service.storage[subject]))
	}
}
