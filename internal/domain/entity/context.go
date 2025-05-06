package entity

import "context"

type Context interface {
	// AddRefresher adds function to refresh items.
	AddRefresher(ContextRefresher)

	// Refresh rereshes its context items.
	Refresh(context.Context) error

	// Items returns slice of ContextItem.
	// If Refresh() is not called, Items may return empty slice.
	Items() []ContextItem
}

type ContextRefresher func(context.Context) error

type ContextItem interface {
	// Text returns itself as string.
	Text() string
}

// var _ Context = (*Context2)(nil)

// // Context2 is contextual knowledge linked to a Session used by Otomo when interpreting Instructions.
// type Context2 struct {
// }

// // AddRefresher implements Context.
// func (c *Context2) AddRefresher(ContextRefresher) {
// 	panic("unimplemented")
// }

// // Refresh implements Context.
// func (c *Context2) Refresh(context.Context) error {
// 	panic("unimplemented")
// }

// // Items implements Context.
// func (c *Context2) Items() []ContextItem {
// 	panic("unimplemented")
// }
