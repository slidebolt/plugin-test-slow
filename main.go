package main

import (
	"context"
	"log"
	"time"

	runner "github.com/slidebolt/sdk-runner"
	"github.com/slidebolt/sdk-types"
)

type SlowPlugin struct{}

func (p *SlowPlugin) OnInitialize(config runner.Config, state types.Storage) (types.Manifest, types.Storage) {
	time.Sleep(1 * time.Second)
	return types.Manifest{ID: "plugin-test-slow", Name: "Slow Loader", Version: "1.0.0", Schemas: types.CoreDomains()}, state
}

func (p *SlowPlugin) OnReady() {}
func (p *SlowPlugin) WaitReady(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *SlowPlugin) OnShutdown()                    {}
func (p *SlowPlugin) OnHealthCheck() (string, error) { return "perfect", nil }
func (p *SlowPlugin) OnStorageUpdate(current types.Storage) (types.Storage, error) {
	return current, nil
}

func (p *SlowPlugin) OnDeviceCreate(dev types.Device) (types.Device, error) {
	return dev, nil
}
func (p *SlowPlugin) OnDeviceUpdate(dev types.Device) (types.Device, error) { return dev, nil }
func (p *SlowPlugin) OnDeviceDelete(id string) error                        { return nil }
func (p *SlowPlugin) OnDevicesList(current []types.Device) ([]types.Device, error) {
	// Include all existing devices plus the hard-coded discovered one.
	for _, d := range current {
		if d.ID == "discovered-after-slow-wait" {
			return runner.EnsureCoreDevice("plugin-test-slow", current), nil
		}
	}
	current = append(current, types.Device{
		ID:         "discovered-after-slow-wait",
		SourceID:   "slow-src",
		SourceName: "Slow Discovered Device",
	})
	return runner.EnsureCoreDevice("plugin-test-slow", current), nil
}
func (p *SlowPlugin) OnDeviceSearch(q types.SearchQuery, res []types.Device) ([]types.Device, error) {
	return res, nil
}

func (p *SlowPlugin) OnEntityCreate(e types.Entity) (types.Entity, error) { return e, nil }
func (p *SlowPlugin) OnEntityUpdate(e types.Entity) (types.Entity, error) { return e, nil }
func (p *SlowPlugin) OnEntityDelete(d, e string) error                    { return nil }
func (p *SlowPlugin) OnEntitiesList(d string, c []types.Entity) ([]types.Entity, error) {
	return runner.EnsureCoreEntities("plugin-test-slow", d, c), nil
}

func (p *SlowPlugin) OnCommand(req types.Command, entity types.Entity) (types.Entity, error) {
	return entity, nil
}

func (p *SlowPlugin) OnEvent(evt types.Event, entity types.Entity) (types.Entity, error) {
	return entity, nil
}

func main() {
	r, err := runner.NewRunner(&SlowPlugin{})
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
