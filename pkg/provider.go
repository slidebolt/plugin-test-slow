// pkg/provider.go - Core provider logic (decoupled from SDK)
package pkg

import (
	"context"
	"time"
)

// Error types for standardized error reporting
type ProviderError string

func (e ProviderError) Error() string { return string(e) }

const (
	ErrOffline      ProviderError = "device is offline"
	ErrUnauthorized ProviderError = "unauthorized"
	ErrTimeout      ProviderError = "operation timed out"
	ErrNotFound     ProviderError = "device not found"
)

// SlowProvider implements the core slow plugin logic
type SlowProvider struct {
	initialized bool
}

// NewSlowProvider creates a new slow provider instance
func NewSlowProvider() *SlowProvider {
	return &SlowProvider{}
}

// Initialize performs slow initialization
func (p *SlowProvider) Initialize(ctx context.Context) error {
	select {
	case <-time.After(1 * time.Second):
		p.initialized = true
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitReady waits for the provider to be ready
func (p *SlowProvider) WaitReady(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// IsInitialized returns whether the provider is initialized
func (p *SlowProvider) IsInitialized() bool {
	return p.initialized
}

// GetHealthStatus returns the health status
func (p *SlowProvider) GetHealthStatus() (string, error) {
	return "perfect", nil
}

// DiscoveryResult represents a discovered device
type DiscoveryResult struct {
	ID         string
	SourceID   string
	SourceName string
}

// DiscoverDevices performs device discovery
func (p *SlowProvider) DiscoverDevices(ctx context.Context, current []DiscoveryResult) ([]DiscoveryResult, error) {
	// Check for existing discovered device
	for _, d := range current {
		if d.ID == "discovered-after-slow-wait" {
			return current, nil
		}
	}

	// Add the hard-coded discovered device
	current = append(current, DiscoveryResult{
		ID:         "discovered-after-slow-wait",
		SourceID:   "slow-src",
		SourceName: "Slow Discovered Device",
	})

	return current, nil
}

// IsDeviceAvailable checks if a device is available
func (p *SlowProvider) IsDeviceAvailable(deviceID string) bool {
	// For test purposes, all devices are available
	return true
}
