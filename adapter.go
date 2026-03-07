// adapter.go - SDK adapter implementation
package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/slidebolt/plugin-test-slow/pkg"
	runner "github.com/slidebolt/sdk-runner"
	"github.com/slidebolt/sdk-types"
)

// SlowPluginAdapter adapts the core provider to the SDK runner interface
type SlowPluginAdapter struct {
	provider *pkg.SlowProvider
}

// NewSlowPluginAdapter creates a new adapter with the given provider
func NewSlowPluginAdapter(provider *pkg.SlowProvider) *SlowPluginAdapter {
	return &SlowPluginAdapter{provider: provider}
}

// OnInitialize initializes the plugin
func (a *SlowPluginAdapter) OnInitialize(config runner.Config, state types.Storage) (types.Manifest, types.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.provider.Initialize(ctx); err != nil {
		// Log initialization error but return manifest anyway
		// The error will be surfaced through health checks
	}

	return types.Manifest{
		ID:      "plugin-test-slow",
		Name:    "Slow Loader",
		Version: "1.0.0",
		Schemas: types.CoreDomains(),
	}, state
}

// OnReady is called when the plugin is ready
func (a *SlowPluginAdapter) OnReady() {}

// WaitReady waits for the plugin to be ready
func (a *SlowPluginAdapter) WaitReady(ctx context.Context) error {
	return a.provider.WaitReady(ctx)
}

// OnShutdown is called during shutdown
func (a *SlowPluginAdapter) OnShutdown() {}

// OnHealthCheck performs a health check
func (a *SlowPluginAdapter) OnHealthCheck() (string, error) {
	return a.provider.GetHealthStatus()
}

// OnStorageUpdate handles storage updates
func (a *SlowPluginAdapter) OnStorageUpdate(current types.Storage) (types.Storage, error) {
	return current, nil
}

// OnDeviceCreate handles device creation
func (a *SlowPluginAdapter) OnDeviceCreate(dev types.Device) (types.Device, error) {
	return dev, nil
}

// OnDeviceUpdate handles device updates
func (a *SlowPluginAdapter) OnDeviceUpdate(dev types.Device) (types.Device, error) {
	return dev, nil
}

// OnDeviceDelete handles device deletion
func (a *SlowPluginAdapter) OnDeviceDelete(id string) error {
	return nil
}

// OnDevicesList handles device listing and discovery
func (a *SlowPluginAdapter) OnDevicesList(current []types.Device) ([]types.Device, error) {
	// Convert current devices to provider format
	var providerDevices []pkg.DiscoveryResult
	for _, d := range current {
		providerDevices = append(providerDevices, pkg.DiscoveryResult{
			ID:         d.ID,
			SourceID:   d.SourceID,
			SourceName: d.SourceName,
		})
	}

	// Perform discovery through the provider
	ctx := context.Background()
	discovered, err := a.provider.DiscoverDevices(ctx, providerDevices)
	if err != nil {
		return nil, err
	}

	// Convert back to SDK types
	var result []types.Device
	for _, d := range discovered {
		result = append(result, types.Device{
			ID:         d.ID,
			SourceID:   d.SourceID,
			SourceName: d.SourceName,
		})
	}

	return runner.EnsureCoreDevice("plugin-test-slow", result), nil
}

// OnDeviceSearch handles device search
func (a *SlowPluginAdapter) OnDeviceSearch(q types.SearchQuery, res []types.Device) ([]types.Device, error) {
	return res, nil
}

// OnEntityCreate handles entity creation
func (a *SlowPluginAdapter) OnEntityCreate(e types.Entity) (types.Entity, error) {
	return a.updateEntityAvailability(e), nil
}

// OnEntityUpdate handles entity updates
func (a *SlowPluginAdapter) OnEntityUpdate(e types.Entity) (types.Entity, error) {
	return a.updateEntityAvailability(e), nil
}

// OnEntityDelete handles entity deletion
func (a *SlowPluginAdapter) OnEntityDelete(d, e string) error {
	return nil
}

// OnEntitiesList handles entity listing
func (a *SlowPluginAdapter) OnEntitiesList(d string, c []types.Entity) ([]types.Entity, error) {
	// Ensure all entities have proper availability and sync status
	var result []types.Entity
	for _, e := range c {
		result = append(result, a.updateEntityAvailability(e))
	}
	return runner.EnsureCoreEntities("plugin-test-slow", d, result), nil
}

// OnCommand handles commands with standardized error reporting
func (a *SlowPluginAdapter) OnCommand(req types.Command, entity types.Entity) (types.Entity, error) {
	// Check device availability before processing command
	if !a.provider.IsDeviceAvailable(entity.DeviceID) {
		// Update entity with failed sync status and error
		entity = a.updateEntitySyncStatus(entity, types.SyncStatusFailed, "device is offline")
		return entity, pkg.ErrOffline
	}

	// Command processed successfully
	entity = a.updateEntitySyncStatus(entity, types.SyncStatusSynced, "")
	return entity, nil
}

// OnEvent handles events with standardized error reporting
func (a *SlowPluginAdapter) OnEvent(evt types.Event, entity types.Entity) (types.Entity, error) {
	// Check device availability before processing event
	if !a.provider.IsDeviceAvailable(entity.DeviceID) {
		// Update entity with failed sync status and error
		entity = a.updateEntitySyncStatus(entity, types.SyncStatusFailed, "device is offline")
		return entity, pkg.ErrOffline
	}

	// Event processed successfully
	entity = a.updateEntitySyncStatus(entity, types.SyncStatusSynced, "")
	return entity, nil
}

// updateEntityAvailability updates the availability entity for a device
func (a *SlowPluginAdapter) updateEntityAvailability(entity types.Entity) types.Entity {
	// If this is an availability entity, update its state
	if entity.Domain == "binary_sensor" && entity.LocalName == "availability" {
		available := a.provider.IsDeviceAvailable(entity.DeviceID)

		// Create availability state
		stateData := map[string]interface{}{
			"state": available,
		}

		reportedBytes, _ := json.Marshal(stateData)
		entity.Data.Reported = reportedBytes
		entity.Data.SyncStatus = types.SyncStatusSynced
		entity.Data.UpdatedAt = time.Now()
	}
	return entity
}

// updateEntitySyncStatus updates the sync status and error information
func (a *SlowPluginAdapter) updateEntitySyncStatus(entity types.Entity, status types.SyncStatus, errorMsg string) types.Entity {
	entity.Data.SyncStatus = status
	entity.Data.UpdatedAt = time.Now()

	if errorMsg != "" {
		// Add error to the reported state
		var reported map[string]interface{}
		if len(entity.Data.Reported) > 0 {
			json.Unmarshal(entity.Data.Reported, &reported)
		}
		if reported == nil {
			reported = make(map[string]interface{})
		}
		reported["error"] = errorMsg

		reportedBytes, _ := json.Marshal(reported)
		entity.Data.Reported = reportedBytes
	}

	return entity
}
