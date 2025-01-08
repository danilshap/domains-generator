package domain

import (
	"context"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) CreateDomain(ctx context.Context, name string) error {
	// Simulate domain creation
	return nil
}

func (p *MockProvider) VerifyDomain(ctx context.Context, name string) (bool, error) {
	// Simulate domain verification
	return true, nil
}

func (p *MockProvider) DeleteDomain(ctx context.Context, name string) error {
	// Simulate domain deletion
	return nil
}

/*
// For future implementation of real providers:

type GoDaddyProvider struct {
    client *godaddy.Client
}

type NamecheapProvider struct {
    client *namecheap.Client
}

type CloudflareProvider struct {
    client *cloudflare.Client
}
*/
