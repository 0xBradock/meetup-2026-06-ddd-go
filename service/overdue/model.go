// Package overdue covers overdue detection and fine management.
// The OverdueCase aggregate owns the full lifecycle from detection through resolution.
package overdue

import (
	"errors"
	"time"
)

type OverdueCaseID string
type FineID string
type LoanID string
type MemberID string

type OverdueCaseStatus string

const (
	OverdueCaseStatusOpen         OverdueCaseStatus = "open"
	OverdueCaseStatusNotified     OverdueCaseStatus = "notified"
	OverdueCaseStatusFineApplied  OverdueCaseStatus = "fine_applied"
	OverdueCaseStatusResolved     OverdueCaseStatus = "resolved"
	OverdueCaseStatusEscalated    OverdueCaseStatus = "escalated"
)

type FineStatus string

const (
	FineStatusPending  FineStatus = "pending"
	FineStatusPaid     FineStatus = "paid"
	FineStatusWaived   FineStatus = "waived"
	FineStatusCanceled FineStatus = "canceled"
)

type FineReason string

const (
	FineReasonOverdue  FineReason = "overdue"
	FineReasonLostCopy FineReason = "lost_copy"
)

// FineAmount is a value object representing a monetary charge.
type FineAmount struct {
	Cents    int64
	Currency string
}

// Fine records the monetary penalty for an overdue case.
type Fine struct {
	ID     FineID
	Amount FineAmount
	Reason FineReason
	Status FineStatus
}

// FinePolicy encodes the rule for calculating a fine.
type FinePolicy struct {
	DailyRateCents int64
	MaxFineCents   int64
}

func (p FinePolicy) Calculate(daysOverdue int) FineAmount {
	cents := p.DailyRateCents * int64(daysOverdue)
	if cents > p.MaxFineCents {
		cents = p.MaxFineCents
	}
	return FineAmount{Cents: cents, Currency: "EUR"}
}

// OverdueCase is the aggregate root.
// It represents a loan that has passed its due date without return.
type OverdueCase struct {
	ID       OverdueCaseID
	LoanID   LoanID
	MemberID MemberID
	Fine     *Fine
	Status   OverdueCaseStatus
	OpenedAt time.Time
}

func (oc *OverdueCase) ApplyFine(fine Fine) error {
	if oc.Status != OverdueCaseStatusOpen && oc.Status != OverdueCaseStatusNotified {
		return errors.New("fine can only be applied to an open or notified case")
	}
	oc.Fine = &fine
	oc.Status = OverdueCaseStatusFineApplied
	return nil
}

func (oc *OverdueCase) WaiveFine() error {
	if oc.Fine == nil || oc.Fine.Status != FineStatusPending {
		return errors.New("no pending fine to waive")
	}
	oc.Fine.Status = FineStatusWaived
	oc.Status = OverdueCaseStatusResolved
	return nil
}

func (oc *OverdueCase) Resolve() error {
	if oc.Fine == nil {
		return errors.New("cannot resolve a case without a fine")
	}
	if oc.Fine.Status != FineStatusPaid && oc.Fine.Status != FineStatusWaived {
		return errors.New("fine must be paid or waived before resolving")
	}
	oc.Status = OverdueCaseStatusResolved
	return nil
}
