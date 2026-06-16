// Package loan manages borrowing, returning, and extending copies.
// The Loan aggregate is the consistency boundary for all borrowing rules.
package loan

import (
	"errors"
	"time"
)

type LoanID string
type MemberID string
type CopyID string

type LoanStatus string

const (
	LoanStatusActive   LoanStatus = "active"
	LoanStatusReturned LoanStatus = "returned"
	LoanStatusOverdue  LoanStatus = "overdue"
	LoanStatusLost     LoanStatus = "lost"
	LoanStatusClosed   LoanStatus = "closed"
)

// DueDate is a value object: it wraps time.Time with a domain-specific rule (must be
// in the future). Two DueDates at the same instant are equal; neither needs an identity.
// Validation at construction means the Loan aggregate can never hold an invalid due date.
type DueDate struct {
	value time.Time
}

func NewDueDate(t time.Time) (DueDate, error) {
	if t.IsZero() {
		return DueDate{}, errors.New("due date is required")
	}
	if t.Before(time.Now()) {
		return DueDate{}, errors.New("due date must be in the future")
	}
	return DueDate{value: t}, nil
}

func (d DueDate) Value() time.Time { return d.value }

// Extension is a value object. It records an immutable fact — when a loan was extended
// and what the new due date became. It has no identity of its own; two extensions at the
// same instant with the same new due date are interchangeable. Stored inside Loan, never
// addressed individually from outside the aggregate.
type Extension struct {
	ExtendedAt time.Time
	NewDueDate DueDate
}

// Borrower is a read model: the minimum projection of a Member the Loan context needs.
// Importing the full Member aggregate would couple two bounded contexts at the model level,
// leaking membership rules into borrowing logic. Only eligibility data crosses the boundary.
type Borrower struct {
	MemberID    MemberID
	Name        string
	IsEligible  bool
	ActiveLoans int
}

func (b Borrower) CanBorrow() error {
	if !b.IsEligible {
		return errors.New("member is not eligible to borrow")
	}
	if b.ActiveLoans >= 5 {
		return errors.New("member has reached the maximum number of active loans")
	}
	return nil
}

// BorrowableCopy is a read model from the Catalog context. The Loan context only needs
// to know whether a copy can be borrowed — bringing in the full Copy or Title aggregate
// would violate the bounded context boundary and expose catalog internals to loan logic.
type BorrowableCopy struct {
	CopyID       CopyID
	IsBorrowable bool
}

// Loan is the aggregate root — the consistency boundary for all borrowing rules.
// All state changes (returning, extending, marking overdue) go through Loan methods so
// invariants (e.g. only active loans can be returned or extended) are always enforced.
// External code never mutates Loan fields directly.
type Loan struct {
	ID         LoanID
	BorrowerID MemberID
	CopyID     CopyID
	Status     LoanStatus
	DueDate    DueDate
	Extensions []Extension
	CreatedAt  time.Time
}

func (l *Loan) Return() error {
	if l.Status != LoanStatusActive {
		return errors.New("only active loans can be returned")
	}
	l.Status = LoanStatusReturned
	return nil
}

func (l *Loan) Extend(newDueDate DueDate) error {
	if l.Status != LoanStatusActive {
		return errors.New("only active loans can be extended")
	}
	l.Extensions = append(l.Extensions, Extension{
		ExtendedAt: time.Now(),
		NewDueDate: newDueDate,
	})
	l.DueDate = newDueDate
	return nil
}

func (l *Loan) MarkOverdue() error {
	if l.Status != LoanStatusActive {
		return errors.New("only active loans can be marked overdue")
	}
	l.Status = LoanStatusOverdue
	return nil
}
