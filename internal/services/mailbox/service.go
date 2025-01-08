package mailbox

import (
	"context"
	"fmt"

	"github.com/danilshap/domains-generator/pkg/utils"
)

// Interface for different mailbox providers
type Provider interface {
	DeleteMailbox(ctx context.Context, address string) error
	CreateMailboxWithPassword(ctx context.Context, address, domain, password string) error
	CreateMailboxBatchWithPassword(ctx context.Context, addresses []string, domain, password string) error
	UpdateMailboxPassword(ctx context.Context, address, password string) error
}

type Service struct {
	provider Provider
	// can add cache, logger, etc.
}

func NewService(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) CreateMailboxWithPassword(ctx context.Context, address, domain, password string) error {
	return s.provider.CreateMailboxWithPassword(ctx, address, domain, password)
}

func (s *Service) DeleteMailbox(ctx context.Context, address string) error {
	return s.provider.DeleteMailbox(ctx, address)
}

func (s *Service) CreateBulkMailboxes(ctx context.Context, prefix, domain, password string, count int) ([]string, error) {
	addresses := make([]string, 0, count)

	for i := 0; i < count; i++ {
		suffix := utils.RandomAlphanumeric(6)
		address := fmt.Sprintf("%s.%s@%s", prefix, suffix, domain)

		if err := s.provider.CreateMailboxWithPassword(ctx, address, domain, password); err != nil {
			return nil, fmt.Errorf("failed to create mailbox %s: %w", address, err)
		}

		addresses = append(addresses, address)
	}

	return addresses, nil

	/* Real implementation:
	addresses := make([]string, 0, count)

	// Create mailboxes in batches of 10 for better performance
	batchSize := 10
	for i := 0; i < count; i += batchSize {
		batch := make([]string, 0, batchSize)
		for j := 0; j < batchSize && i+j < count; j++ {
			suffix := utils.RandomAlphanumeric(6)
			address := fmt.Sprintf("%s.%s@%s", prefix, suffix, domain)
			batch = append(batch, address)
		}

		// Create batch of mailboxes
		if err := s.provider.CreateMailboxBatchWithPassword(ctx, batch, domain, password); err != nil {
			return nil, fmt.Errorf("failed to create mailbox batch: %w", err)
		}

		addresses = append(addresses, batch...)
	}

	return addresses, nil
	*/
}

func (s *Service) UpdatePassword(ctx context.Context, address, domain, newPassword string) error {
	return s.provider.UpdateMailboxPassword(ctx, address, newPassword)
}
