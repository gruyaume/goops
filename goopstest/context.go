package goopstest

type Context struct {
	Charm func() error
}

func (c *Context) Run(event string, state State) *State {
	return &State{
		UnitStatus: "active",
	}
}

type State struct {
	UnitStatus string
}
