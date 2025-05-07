package subpub

type subscription struct {
	subject string
	srv     *Service
	sub     *subscriber
}

func (s *subscription) Unsubscribe() {
	s.srv.mu.Lock()
	defer s.srv.mu.Unlock()

	subjectSubs := s.srv.storage[s.subject]
	for i, sub := range subjectSubs {
		if sub == s.sub {
			s.srv.storage[s.subject] = append(subjectSubs[:i], subjectSubs[i+1:]...)
			return
		}
	}
}

func newSubscription(subject string, service *Service, subscriber *subscriber) *subscription {
	return &subscription{
		subject: subject,
		srv:     service,
		sub:     subscriber,
	}
}
