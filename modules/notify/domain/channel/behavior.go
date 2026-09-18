package channel

func (c *Channel) Kind() string  { return c.state.Kind }
func (c *Channel) Enabled() bool { return c.state.Enabled }
func (c *Channel) Name() string  { return c.state.Name }
