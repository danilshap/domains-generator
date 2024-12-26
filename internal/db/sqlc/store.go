package db

import (
	"context"
	"database/sql"
	"log"
)

type Store interface {
	CreateDomain(ctx context.Context, arg CreateDomainParams) (Domain, error)
	GetDomainByID(ctx context.Context, id int32) (Domain, error)
	GetAllDomains(ctx context.Context, arg GetAllDomainsParams) ([]GetAllDomainsRow, error)
	GetDomainByName(ctx context.Context, name string) (GetDomainByNameRow, error)
	GetDomainsCount(ctx context.Context) (int64, error)
	UpdateDomain(ctx context.Context, arg UpdateDomainParams) error
	SetDomainStatus(ctx context.Context, arg SetDomainStatusParams) error
	DeleteDomain(ctx context.Context, id int32) error
	UpdateDomainAndMailboxesStatus(ctx context.Context, domainID int32, status int32) error

	CreateMailbox(ctx context.Context, arg CreateMailboxParams) (Mailbox, error)
	GetMailboxByID(ctx context.Context, id int32) (Mailbox, error)
	GetAllMailboxes(ctx context.Context, arg GetAllMailboxesParams) ([]GetAllMailboxesRow, error)
	GetMailboxesByDomain(ctx context.Context, address string) ([]GetMailboxesByDomainRow, error)
	GetMailboxesByDomainID(ctx context.Context, arg GetMailboxesByDomainIDParams) ([]Mailbox, error)
	GetMailboxesCountByDomainID(ctx context.Context, domainID int32) (int64, error)
	GetMailboxesCount(ctx context.Context) (int64, error)
	UpdateMailbox(ctx context.Context, arg UpdateMailboxParams) error
	SetMailboxStatus(ctx context.Context, arg SetMailboxStatusParams) error
	DeleteMailbox(ctx context.Context, id int32) error
	GetMailboxesWithFilters(ctx context.Context, arg GetMailboxesWithFiltersParams) ([]GetMailboxesWithFiltersRow, error)
	GetMailboxesStats(ctx context.Context) (GetMailboxesStatsRow, error)
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(fn func(*Queries) error) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

// GetDB returns the underlying database connection
func (store *SQLStore) GetDB() *sql.DB {
	return store.db
}

func (store *SQLStore) UpdateDomainAndMailboxesStatus(ctx context.Context, domainID int32, status int32) error {
	err := store.execTx(func(q *Queries) error {
		err := q.SetDomainStatus(ctx, SetDomainStatusParams{
			ID:     domainID,
			Status: status,
		})
		if err != nil {
			return err
		}

		if status == 2 {
			err = q.UpdateMailboxesStatusByDomainID(ctx, UpdateMailboxesStatusByDomainIDParams{
				Status:   status,
				DomainID: domainID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err == nil && status == 2 {
		go func() {
			notifyCtx := context.Background()
			mailboxes, err := store.GetMailboxesByDomainID(notifyCtx, GetMailboxesByDomainIDParams{
				DomainID: domainID,
			})
			if err != nil {
				log.Printf("Error fetching mailboxes for notifications: %v", err)
				return
			}

			for _, mailbox := range mailboxes {
				log.Printf("Notifying user about mailbox status change: %s", mailbox.Address)
			}
		}()
	}

	return err
}
