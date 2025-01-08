package mailbox

import (
	"context"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) DeleteMailbox(ctx context.Context, address string) error {
	// Simulate mailbox deletion
	return nil
}

func (p *MockProvider) CreateMailboxWithPassword(ctx context.Context, address, domain, password string) error {
	// Simulate mailbox creation with password
	return nil
}

func (p *MockProvider) CreateMailboxBatchWithPassword(ctx context.Context, addresses []string, domain, password string) error {
	// Simulate batch creation with password
	return nil
}

func (p *MockProvider) UpdateMailboxPassword(ctx context.Context, address, password string) error {
	// Simulate password update
	return nil
}

/*
// For future implementation of real providers:

type GoogleWorkspaceProvider struct {
	client *googleapi.Client
}

type Office365Provider struct {
	client *msgraph.Client
}

type ExchangeProvider struct {
	client *exchange.Client
}
*/
