package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/slidebolt/plugin-test-slow/pkg"
	runner "github.com/slidebolt/sdk-runner"
	"github.com/slidebolt/sdk-types"
)

func TestNewSlowPluginAdapter(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	if adapter == nil {
		t.Fatal("NewSlowPluginAdapter() returned nil")
	}

	if adapter.provider != provider {
		t.Error("Adapter provider mismatch")
	}
}

func TestSlowPluginAdapter_OnInitialize(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	config := runner.Config{}
	state := types.Storage{Meta: "test-meta", Data: json.RawMessage(`{}`)}

	manifest, newState := adapter.OnInitialize(config, state)

	if manifest.ID != "plugin-test-slow" {
		t.Errorf("Manifest.ID = %q, want %q", manifest.ID, "plugin-test-slow")
	}

	if manifest.Name != "Slow Loader" {
		t.Errorf("Manifest.Name = %q, want %q", manifest.Name, "Slow Loader")
	}

	if manifest.Version != "1.0.0" {
		t.Errorf("Manifest.Version = %q, want %q", manifest.Version, "1.0.0")
	}

	// State Meta should be preserved
	if newState.Meta != state.Meta {
		t.Errorf("OnInitialize should preserve state.Meta = %q, got %q", state.Meta, newState.Meta)
	}
}

func TestSlowPluginAdapter_OnHealthCheck(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	status, err := adapter.OnHealthCheck()
	if err != nil {
		t.Fatalf("OnHealthCheck() failed: %v", err)
	}

	if status != "perfect" {
		t.Errorf("OnHealthCheck() returned %q, want %q", status, "perfect")
	}
}

func TestSlowPluginAdapter_OnDevicesList(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	// Test with empty current list
	current := []types.Device{}
	result, err := adapter.OnDevicesList(current)
	if err != nil {
		t.Fatalf("OnDevicesList() with empty list failed: %v", err)
	}

	// Should have discovered device + any core entities
	if len(result) == 0 {
		t.Fatal("OnDevicesList() returned no devices")
	}

	// Check that our discovered device is in the list
	found := false
	for _, d := range result {
		if d.ID == "discovered-after-slow-wait" {
			found = true
			if d.SourceID != "slow-src" {
				t.Errorf("Discovered device SourceID = %q, want %q", d.SourceID, "slow-src")
			}
			if d.SourceName != "Slow Discovered Device" {
				t.Errorf("Discovered device SourceName = %q, want %q", d.SourceName, "Slow Discovered Device")
			}
			break
		}
	}

	if !found {
		t.Error("OnDevicesList() did not include discovered device")
	}

	// Test with existing device - should not duplicate
	result2, err := adapter.OnDevicesList(result)
	if err != nil {
		t.Fatalf("OnDevicesList() with existing devices failed: %v", err)
	}

	// Count occurrences of our discovered device
	count := 0
	for _, d := range result2 {
		if d.ID == "discovered-after-slow-wait" {
			count++
		}
	}

	if count > 1 {
		t.Errorf("OnDevicesList() duplicated the discovered device (found %d times)", count)
	}
}

func TestSlowPluginAdapter_OnCommand_Success(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	entity := types.Entity{
		ID:       "test-entity",
		DeviceID: "test-device",
		Domain:   "light",
		Data:     types.EntityData{},
	}

	req := types.Command{
		ID:         "cmd-1",
		PluginID:   "plugin-test-slow",
		DeviceID:   "test-device",
		EntityID:   "test-entity",
		EntityType: "light",
		Payload:    json.RawMessage(`{"state":"on"}`),
		CreatedAt:  time.Now(),
	}

	result, err := adapter.OnCommand(req, entity)
	if err != nil {
		t.Fatalf("OnCommand() failed: %v", err)
	}

	// Verify sync status is "synced"
	if result.Data.SyncStatus != "synced" {
		t.Errorf("OnCommand() SyncStatus = %q, want %q", result.Data.SyncStatus, "synced")
	}

	// Verify no error in reported state
	var reported map[string]interface{}
	if len(result.Data.Reported) > 0 {
		json.Unmarshal(result.Data.Reported, &reported)
		if _, exists := reported["error"]; exists {
			t.Error("OnCommand() should not have error in reported state for successful command")
		}
	}
}

func TestSlowPluginAdapter_OnEvent_Success(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	entity := types.Entity{
		ID:       "test-entity",
		DeviceID: "test-device",
		Domain:   "sensor",
		Data:     types.EntityData{},
	}

	evt := types.Event{
		ID:         "evt-1",
		PluginID:   "plugin-test-slow",
		DeviceID:   "test-device",
		EntityID:   "test-entity",
		EntityType: "sensor",
		Payload:    json.RawMessage(`{"value":42}`),
		CreatedAt:  time.Now(),
	}

	result, err := adapter.OnEvent(evt, entity)
	if err != nil {
		t.Fatalf("OnEvent() failed: %v", err)
	}

	// Verify sync status is "synced"
	if result.Data.SyncStatus != "synced" {
		t.Errorf("OnEvent() SyncStatus = %q, want %q", result.Data.SyncStatus, "synced")
	}
}

func TestSlowPluginAdapter_OnEntitiesList(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	entities := []types.Entity{
		{
			ID:        "entity-1",
			DeviceID:  "device-1",
			Domain:    "light",
			LocalName: "test-light",
			Data:      types.EntityData{},
		},
	}

	result, err := adapter.OnEntitiesList("device-1", entities)
	if err != nil {
		t.Fatalf("OnEntitiesList() failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("OnEntitiesList() returned no entities")
	}
}

func TestSlowPluginAdapter_updateEntityAvailability(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	// Test availability entity
	entity := types.Entity{
		ID:        "avail-1",
		DeviceID:  "device-1",
		Domain:    "binary_sensor",
		LocalName: "availability",
		Data:      types.EntityData{},
	}

	result := adapter.updateEntityAvailability(entity)

	if result.Data.SyncStatus != "synced" {
		t.Errorf("Availability entity SyncStatus = %q, want %q", result.Data.SyncStatus, "synced")
	}

	// Verify reported state contains availability
	var reported map[string]interface{}
	if len(result.Data.Reported) > 0 {
		json.Unmarshal(result.Data.Reported, &reported)
		state, exists := reported["state"]
		if !exists {
			t.Error("Availability entity reported state should contain 'state' field")
		}
		if state != true {
			t.Errorf("Availability entity state = %v, want true", state)
		}
	} else {
		t.Error("Availability entity should have reported data")
	}

	// Test non-availability entity (should not be modified)
	normalEntity := types.Entity{
		ID:        "light-1",
		DeviceID:  "device-1",
		Domain:    "light",
		LocalName: "living-room",
		Data:      types.EntityData{},
	}

	result2 := adapter.updateEntityAvailability(normalEntity)
	if len(result2.Data.Reported) > 0 {
		t.Error("Non-availability entity should not have reported data modified")
	}
}

func TestSlowPluginAdapter_updateEntitySyncStatus(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	entity := types.Entity{
		ID:       "test-entity",
		DeviceID: "test-device",
		Domain:   "light",
		Data:     types.EntityData{},
	}

	// Test setting to "synced" without error
	result := adapter.updateEntitySyncStatus(entity, types.SyncStatusSynced, "")

	if result.Data.SyncStatus != "synced" {
		t.Errorf("updateEntitySyncStatus() SyncStatus = %q, want %q", result.Data.SyncStatus, "synced")
	}

	// Test setting to "failed" with error message
	result2 := adapter.updateEntitySyncStatus(entity, types.SyncStatusFailed, "device is offline")

	if result2.Data.SyncStatus != "failed" {
		t.Errorf("updateEntitySyncStatus() SyncStatus = %q, want %q", result2.Data.SyncStatus, "failed")
	}

	// Verify error in reported state
	var reported map[string]interface{}
	if len(result2.Data.Reported) > 0 {
		json.Unmarshal(result2.Data.Reported, &reported)
		errMsg, exists := reported["error"]
		if !exists {
			t.Error("Failed entity should have error in reported state")
		}
		if errMsg != "device is offline" {
			t.Errorf("Error message = %v, want %q", errMsg, "device is offline")
		}
	} else {
		t.Error("Failed entity should have reported data with error")
	}

	// Test "pending" status
	result3 := adapter.updateEntitySyncStatus(entity, types.SyncStatusPending, "")
	if result3.Data.SyncStatus != "pending" {
		t.Errorf("updateEntitySyncStatus() SyncStatus = %q, want %q", result3.Data.SyncStatus, "pending")
	}
}

func TestSlowPluginAdapter_PassThroughHandlers(t *testing.T) {
	provider := pkg.NewSlowProvider()
	adapter := NewSlowPluginAdapter(provider)

	// Test OnReady (should not panic)
	adapter.OnReady()

	// Test OnShutdown (should not panic)
	adapter.OnShutdown()

	// Test OnStorageUpdate
	state := types.Storage{Meta: "test"}
	newState, err := adapter.OnStorageUpdate(state)
	if err != nil {
		t.Errorf("OnStorageUpdate() error = %v", err)
	}
	if newState.Meta != state.Meta {
		t.Errorf("OnStorageUpdate should preserve state.Meta = %q, got %q", state.Meta, newState.Meta)
	}

	// Test OnDeviceCreate
	device := types.Device{ID: "test-device"}
	result, err := adapter.OnDeviceCreate(device)
	if err != nil {
		t.Errorf("OnDeviceCreate() error = %v", err)
	}
	if result.ID != device.ID {
		t.Error("OnDeviceCreate should preserve device")
	}

	// Test OnDeviceUpdate
	result, err = adapter.OnDeviceUpdate(device)
	if err != nil {
		t.Errorf("OnDeviceUpdate() error = %v", err)
	}
	if result.ID != device.ID {
		t.Error("OnDeviceUpdate should preserve device")
	}

	// Test OnDeviceDelete
	err = adapter.OnDeviceDelete("test-device")
	if err != nil {
		t.Errorf("OnDeviceDelete() error = %v", err)
	}

	// Test OnDeviceSearch
	searchResult, err := adapter.OnDeviceSearch(types.SearchQuery{}, []types.Device{})
	if err != nil {
		t.Errorf("OnDeviceSearch() error = %v", err)
	}
	if len(searchResult) != 0 {
		t.Error("OnDeviceSearch should return unchanged result")
	}

	// Test OnEntityCreate
	entity := types.Entity{ID: "test-entity"}
	entityResult, err := adapter.OnEntityCreate(entity)
	if err != nil {
		t.Errorf("OnEntityCreate() error = %v", err)
	}
	if entityResult.ID != entity.ID {
		t.Error("OnEntityCreate should preserve entity ID")
	}

	// Test OnEntityUpdate
	entityResult, err = adapter.OnEntityUpdate(entity)
	if err != nil {
		t.Errorf("OnEntityUpdate() error = %v", err)
	}
	if entityResult.ID != entity.ID {
		t.Error("OnEntityUpdate should preserve entity ID")
	}

	// Test OnEntityDelete
	err = adapter.OnEntityDelete("device-1", "entity-1")
	if err != nil {
		t.Errorf("OnEntityDelete() error = %v", err)
	}
}
