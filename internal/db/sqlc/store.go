package db

import (
	"context"
	"database/sql"
	"log"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(fn func(*Queries) error) error {
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
func (store *Store) GetDB() *sql.DB {
	return store.db
}

func (store *Store) UpdateDomainAndMailboxesStatus(ctx context.Context, domainID int32, status int32) error {
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
