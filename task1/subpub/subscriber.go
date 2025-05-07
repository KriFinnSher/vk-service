package subpub

const messageBuffer = 16

type subscriber struct {
	handler MessageHandler
	message chan interface{}
}

func newSubscriber(handler MessageHandler) *subscriber {
	return &subscriber{
		handler: handler,
		message: make(chan interface{}, messageBuffer),
	}
}
