package event

type Handler func(event Event) error

type Publisher interface {
	Subscribe(kind Kind, hanler Handler)
	Publish(event Event) error
}
