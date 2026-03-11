package synch_contract

type PubSub interface {
	Publish(channel string, message []byte) error
	Subscribe(channel string) (Subscription, error)
}

type Subscription interface {
	Channel() <-chan []byte
	Unsubscribe() error
}
