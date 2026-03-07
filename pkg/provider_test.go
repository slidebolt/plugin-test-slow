package pkg

import (
	"context"
	"testing"
	"time"
)

func TestNewSlowProvider(t *testing.T) {
	provider := NewSlowProvider()
	if provider == nil {
		t.Fatal("NewSlowProvider() returned nil")
	}
	if provider.initialized {
		t.Error("New provider should not be initialized")
	}
}

func TestSlowProvider_Initialize(t *testing.T) {
	provider := NewSlowProvider()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := provider.Initialize(ctx)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if !provider.initialized {
		t.Error("Provider should be initialized after Initialize()")
	}

	// Verify it takes at least 1 second (the sleep duration)
	start := time.Now()
	provider2 := NewSlowProvider()
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	provider2.Initialize(ctx2)
	elapsed := time.Since(start)

	if elapsed < 1*time.Second {
		t.Errorf("Initialize() should take at least 1 second, took %v", elapsed)
	}
}

func TestSlowProvider_Initialize_ContextCancel(t *testing.T) {
	provider := NewSlowProvider()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := provider.Initialize(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Initialize() with cancelled context should return DeadlineExceeded, got: %v", err)
	}
}

func TestSlowProvider_WaitReady(t *testing.T) {
	provider := NewSlowProvider()

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	start := time.Now()
	err := provider.WaitReady(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("WaitReady() failed: %v", err)
	}

	// Should take at least 5 seconds
	if elapsed < 5*time.Second {
		t.Errorf("WaitReady() should take at least 5 seconds, took %v", elapsed)
	}
}

func TestSlowProvider_WaitReady_ContextCancel(t *testing.T) {
	provider := NewSlowProvider()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := provider.WaitReady(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("WaitReady() with cancelled context should return DeadlineExceeded, got: %v", err)
	}
}

func TestSlowProvider_GetHealthStatus(t *testing.T) {
	provider := NewSlowProvider()

	status, err := provider.GetHealthStatus()
	if err != nil {
		t.Fatalf("GetHealthStatus() failed: %v", err)
	}

	if status != "perfect" {
		t.Errorf("GetHealthStatus() returned %q, want %q", status, "perfect")
	}
}

func TestSlowProvider_IsDeviceAvailable(t *testing.T) {
	provider := NewSlowProvider()

	// All devices should be available in test mode
	if !provider.IsDeviceAvailable("test-device") {
		t.Error("IsDeviceAvailable() should return true for test devices")
	}

	if !provider.IsDeviceAvailable("") {
		t.Error("IsDeviceAvailable() should return true even for empty device ID")
	}
}

func TestSlowProvider_DiscoverDevices(t *testing.T) {
	provider := NewSlowProvider()

	ctx := context.Background()

	// Test with empty current list
	current := []DiscoveryResult{}
	discovered, err := provider.DiscoverDevices(ctx, current)
	if err != nil {
		t.Fatalf("DiscoverDevices() with empty list failed: %v", err)
	}

	if len(discovered) != 1 {
		t.Fatalf("DiscoverDevices() returned %d devices, want 1", len(discovered))
	}

	if discovered[0].ID != "discovered-after-slow-wait" {
		t.Errorf("DiscoverDevices() returned device with ID %q, want %q", discovered[0].ID, "discovered-after-slow-wait")
	}

	if discovered[0].SourceID != "slow-src" {
		t.Errorf("DiscoverDevices() returned device with SourceID %q, want %q", discovered[0].SourceID, "slow-src")
	}

	if discovered[0].SourceName != "Slow Discovered Device" {
		t.Errorf("DiscoverDevices() returned device with SourceName %q, want %q", discovered[0].SourceName, "Slow Discovered Device")
	}

	// Test when device already exists - should not duplicate
	discovered2, err := provider.DiscoverDevices(ctx, discovered)
	if err != nil {
		t.Fatalf("DiscoverDevices() with existing device failed: %v", err)
	}

	if len(discovered2) != 1 {
		t.Errorf("DiscoverDevices() with existing device returned %d devices, want 1 (no duplicates)", len(discovered2))
	}
}

func TestProviderError(t *testing.T) {
	err := ErrOffline
	if err.Error() != "device is offline" {
		t.Errorf("ErrOffline.Error() = %q, want %q", err.Error(), "device is offline")
	}

	err = ErrUnauthorized
	if err.Error() != "unauthorized" {
		t.Errorf("ErrUnauthorized.Error() = %q, want %q", err.Error(), "unauthorized")
	}

	err = ErrTimeout
	if err.Error() != "operation timed out" {
		t.Errorf("ErrTimeout.Error() = %q, want %q", err.Error(), "operation timed out")
	}

	err = ErrNotFound
	if err.Error() != "device not found" {
		t.Errorf("ErrNotFound.Error() = %q, want %q", err.Error(), "device not found")
	}
}

func TestSlowProvider_IsInitialized(t *testing.T) {
	provider := NewSlowProvider()

	if provider.IsInitialized() {
		t.Error("IsInitialized() should return false for new provider")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	provider.Initialize(ctx)

	if !provider.IsInitialized() {
		t.Error("IsInitialized() should return true after Initialize()")
	}
}
