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

// DueDate is a value object. Validates on construction.
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

// Extension records a single loan extension.
type Extension struct {
	ExtendedAt time.Time
	NewDueDate DueDate
}

// Borrower is the Member seen from the borrowing perspective.
// The Loan context never holds a full Member aggregate.
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

// BorrowableCopy is the Copy seen from the borrowing perspective.
type BorrowableCopy struct {
	CopyID       CopyID
	IsBorrowable bool
}

// Loan is the aggregate root.
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
