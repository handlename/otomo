package service

import "context"

type Messenger interface {
	Send(context.Context, string) error
}
