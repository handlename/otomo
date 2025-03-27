package entity

import "context"

type Brain interface {
	Think(context.Context, Context, *Instruction) (*Answer, error)
}
