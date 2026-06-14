// Package membership manages members and their borrowing eligibility.
// The Member aggregate is the authoritative source for whether a patron can borrow.
package membership

import (
	"errors"
	"time"
)

type MemberID string

type MembershipStatus string

const (
	MembershipStatusActive    MembershipStatus = "active"
	MembershipStatusSuspended MembershipStatus = "suspended"
	MembershipStatusExpired   MembershipStatus = "expired"
	MembershipStatusClosed    MembershipStatus = "closed"
)

type BorrowingStatus string

const (
	BorrowingStatusAllowed BorrowingStatus = "allowed"
	BorrowingStatusBlocked BorrowingStatus = "blocked"
	BorrowingStatusLimited BorrowingStatus = "limited"
)

// ContactInformation is a value object.
type ContactInformation struct {
	Email string
	Phone string
}

// BorrowingPrivilege records the current borrowing rights for a member.
type BorrowingPrivilege struct {
	Status         BorrowingStatus
	MaxActiveLoans int
}

// Member is the aggregate root for the Membership context.
type Member struct {
	ID                 MemberID
	Name               string
	MembershipStatus   MembershipStatus
	BorrowingPrivilege BorrowingPrivilege
	ContactInformation ContactInformation
	RegisteredAt       time.Time
}

func (m *Member) IsEligibleToBorrow() bool {
	return m.MembershipStatus == MembershipStatusActive &&
		m.BorrowingPrivilege.Status == BorrowingStatusAllowed
}

func (m *Member) Block() error {
	if m.MembershipStatus != MembershipStatusActive {
		return errors.New("only active members can be blocked")
	}
	m.BorrowingPrivilege.Status = BorrowingStatusBlocked
	return nil
}

func (m *Member) Unblock() error {
	if m.BorrowingPrivilege.Status != BorrowingStatusBlocked {
		return errors.New("member is not blocked")
	}
	m.BorrowingPrivilege.Status = BorrowingStatusAllowed
	return nil
}

func (m *Member) Suspend() error {
	if m.MembershipStatus != MembershipStatusActive {
		return errors.New("only active members can be suspended")
	}
	m.MembershipStatus = MembershipStatusSuspended
	return nil
}
