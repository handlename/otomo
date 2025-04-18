package event

type Publisher interface {
	Subscribe(Kind, Subscriber)
	Publish(Event) error
}
