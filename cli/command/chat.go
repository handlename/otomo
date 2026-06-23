package command

type ChatCmd struct {
	MCP     bool `help:"Run MCP server in background." name:"mcp"`
	MCPPort int  `help:"Port for MCP server to listen on." name:"mcp-port" default:"8000"`
}

func (c *ChatCmd) Run(ctx *Context) error {
	return ctx.App.RunChat(ctx.Ctx, c.MCP, c.MCPPort)
}
