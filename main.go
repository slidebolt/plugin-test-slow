package main

import (
	"log"
	"time"

	runner "github.com/slidebolt/sdk-runner"
	"github.com/slidebolt/sdk-types"
)

type SlowPlugin struct{}

func (p *SlowPlugin) OnInitialize(config runner.Config, state types.Storage) (types.Manifest, types.Storage) {
	time.Sleep(1 * time.Second)
	return types.Manifest{ID: "plugin-test-slow", Name: "Slow Loader", Version: "1.0.0"}, state
}

func (p *SlowPlugin) OnReady()                       {}
func (p *SlowPlugin) OnHealthCheck() (string, error) { return "perfect", nil }
func (p *SlowPlugin) OnStorageUpdate(current types.Storage) (types.Storage, error) {
	return current, nil
}

func (p *SlowPlugin) OnDeviceCreate(dev types.Device) (types.Device, error) {
	dev.Config = types.Storage{Meta: "slow-metadata"}
	return dev, nil
}
func (p *SlowPlugin) OnDeviceUpdate(dev types.Device) (types.Device, error) { return dev, nil }
func (p *SlowPlugin) OnDeviceDelete(id string) error                        { return nil }
func (p *SlowPlugin) OnDevicesList(current []types.Device) ([]types.Device, error) {
	return current, nil
}
func (p *SlowPlugin) OnDeviceSearch(q types.SearchQuery, res []types.Device) ([]types.Device, error) {
	return res, nil
}

func (p *SlowPlugin) OnEntityCreate(e types.Entity) (types.Entity, error) { return e, nil }
func (p *SlowPlugin) OnEntityUpdate(e types.Entity) (types.Entity, error) { return e, nil }
func (p *SlowPlugin) OnEntityDelete(d, e string) error                    { return nil }
func (p *SlowPlugin) OnEntitiesList(d string, c []types.Entity) ([]types.Entity, error) {
	return c, nil
}

func (p *SlowPlugin) OnCommand(cmd types.Command, entity types.Entity) (types.Entity, error) {
	return entity, nil
}
func (p *SlowPlugin) OnEvent(evt types.Event, entity types.Entity) (types.Entity, error) {
	return entity, nil
}

func main() {
	if err := runner.NewRunner(&SlowPlugin{}).Run(); err != nil {
		log.Fatal(err)
	}
}
