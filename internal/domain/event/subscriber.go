package event

import "context"

type Subscriber func(context.Context, Event) error
