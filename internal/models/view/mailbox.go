package view

import (
	"time"

	db "github.com/danilshap/domains-generator/internal/db/sqlc"
)

type MailboxView struct {
	ID         int32
	Address    string
	Password   string
	Status     int32
	CreatedAt  time.Time
	DomainID   int32
	DomainName string
}

type MailboxListData struct {
	Mailboxes   []MailboxView
	CurrentPage int32
	TotalPages  int32
	PageSize    int32
	DomainID    int32
}

// ToMailboxViewFromDomain конвертирует из GetMailboxesByDomainIDRow
func ToMailboxViewFromDomain(m db.GetMailboxesByDomainIDRow) MailboxView {
	return MailboxView{
		ID:         m.ID,
		Address:    m.Address,
		Password:   m.Password,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt.Time,
		DomainID:   m.DomainID,
		DomainName: m.DomainName,
	}
}

// ToMailboxViewFromFilter конвертирует из GetMailboxesWithFiltersRow
func ToMailboxViewFromFilter(m db.GetMailboxesWithFiltersRow) MailboxView {
	return MailboxView{
		ID:         m.ID,
		Address:    m.Address,
		Password:   m.Password,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt.Time,
		DomainID:   m.DomainID,
		DomainName: m.DomainName.String,
	}
}
