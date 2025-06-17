package goopstest

import (
	"fmt"

	"github.com/canonical/pebble/client"
	"github.com/gruyaume/goops"
)

type fakePebbleGetter struct {
	CanConnect bool
}

type FakePebbleClient struct {
	CanConnect bool
}

func (f *FakePebbleClient) AddLayer(*client.AddLayerOptions) error {
	return nil
}

func (f *FakePebbleClient) Exec(*client.ExecOptions) (goops.PebbleExecProcess, error) {
	return nil, nil
}

func (f *FakePebbleClient) Pull(*client.PullOptions) error {
	return nil
}

func (f *FakePebbleClient) Push(*client.PushOptions) error {
	return nil
}

func (f *FakePebbleClient) Restart(*client.ServiceOptions) (string, error) {
	return "", nil
}

func (f *FakePebbleClient) Start(*client.ServiceOptions) (string, error) {
	return "", nil
}

func (f *FakePebbleClient) Stop(*client.ServiceOptions) (string, error) {
	return "", nil
}

func (f *FakePebbleClient) SysInfo() (*client.SysInfo, error) {
	if !f.CanConnect {
		return nil, fmt.Errorf("cannot connect to Pebble")
	}

	return nil, nil
}

func (f *FakePebbleClient) WaitChange(string, *client.WaitChangeOptions) (*client.Change, error) {
	return nil, nil
}

func (f *fakePebbleGetter) Pebble(string) goops.PebbleClient {
	return &FakePebbleClient{
		CanConnect: f.CanConnect,
	}
}

type Container struct {
	Name            string
	CanConnect      bool
	Layers          map[string][]byte
	ServiceStatuses map[string]client.ServiceStatus
	Mounts          map[string]Mount
	Execs           []Exec
	Notices         []client.Notice
	CheckInfos      []client.CheckInfo
}
