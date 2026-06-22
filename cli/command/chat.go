package command

type ChatCmd struct{}

func (c *ChatCmd) Run(ctx *Context) error {
	return ctx.App.RunChat(ctx.Ctx)
}
