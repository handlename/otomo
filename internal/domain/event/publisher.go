package event

import "context"

type Publisher interface {
	Subscribe(Kind, Subscriber)
	Publish(context.Context, Event) error
}
