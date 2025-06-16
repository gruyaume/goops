package goops

import (
	"fmt"

	"github.com/canonical/pebble/client"
)

var defaultPebbleGetter PebbleGetter = realPebbleGetter{}

type PebbleGetter interface {
	GetPebble(container string) PebbleClient
}

func GetPebble(container string) PebbleClient {
	return defaultPebbleGetter.GetPebble(container)
}

type PebbleClient interface {
	Stop(opts *client.ServiceOptions) (changeID string, err error)
	WaitChange(changeID string, options *client.WaitChangeOptions) (*client.Change, error)
	Exec(opts *client.ExecOptions) (PebbleExecProcess, error)
	SysInfo() (*client.SysInfo, error)
	Push(opts *client.PushOptions) error
	Pull(opts *client.PullOptions) error
	AddLayer(opts *client.AddLayerOptions) error
	Restart(opts *client.ServiceOptions) (changeID string, err error)
	Start(opts *client.ServiceOptions) (changeID string, err error)
}

type PebbleExecProcess interface {
	Wait() error
	SendResize(width int, height int) error
	SendSignal(signal string) error
}

type realPebbleGetter struct{}

func (g realPebbleGetter) GetPebble(container string) PebbleClient {
	pebble, err := client.New(&client.Config{
		Socket: fmt.Sprintf("/charm/containers/%s/pebble.socket", container),
	})
	if err != nil {
		panic(err) // shouldn't happen
	}

	return &realPebbleStub{pebble}
}

type realPebbleStub struct {
	*client.Client
}

func (p *realPebbleStub) Exec(opts *client.ExecOptions) (PebbleExecProcess, error) {
	return p.Client.Exec(opts)
}
