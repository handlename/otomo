package event

type Subscriber func(event Event) error
