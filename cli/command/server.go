package command

type Server struct{}

func (s *Server) Run(c *Context) error {
	return c.App.Run(c.Ctx)
}
