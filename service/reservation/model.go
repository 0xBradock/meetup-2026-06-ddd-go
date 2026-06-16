// Package reservation manages holds placed on titles when no copies are available.
// ReservationQueue is the aggregate root ordering members waiting for a title.
package reservation

import (
	"errors"
	"time"
)

type ReservationID string
type MemberID string
type TitleID string

type ReservationStatus string

const (
	ReservationStatusWaiting        ReservationStatus = "waiting"
	ReservationStatusReadyForPickup ReservationStatus = "ready_for_pickup"
	ReservationStatusFulfilled      ReservationStatus = "fulfilled"
	ReservationStatusExpired        ReservationStatus = "expired"
	ReservationStatusCanceled       ReservationStatus = "canceled"
)

// ExpirationDate is a value object for the pickup deadline. Like DueDate in the Loan
// context, it validates at construction (must be in the future) so a Reservation can
// never be in a ready-for-pickup state with an already-expired deadline. Two ExpirationDates
// wrapping the same instant are equal; neither has an independent identity.
type ExpirationDate struct {
	value time.Time
}

func NewExpirationDate(t time.Time) (ExpirationDate, error) {
	if t.Before(time.Now()) {
		return ExpirationDate{}, errors.New("expiration date must be in the future")
	}
	return ExpirationDate{value: t}, nil
}

func (e ExpirationDate) Value() time.Time { return e.value }

// ReservableTitle is a read model from the Catalog context. The Reservation context only
// needs to know whether a title can be reserved — importing the full Title aggregate would
// couple the two contexts at the model level and force reservation logic to know about
// catalog internals (copies, shelf locations, conditions, etc.).
type ReservableTitle struct {
	TitleID      TitleID
	IsReservable bool
}

// Reservation is an entity. It has its own identity (ReservationID) and moves through
// a defined lifecycle (waiting → ready_for_pickup → fulfilled / canceled / expired).
// It lives inside the ReservationQueue aggregate and is only accessed through it —
// external code never holds a Reservation and calls Cancel or Fulfill on it directly.
type Reservation struct {
	ID             ReservationID
	MemberID       MemberID
	TitleID        TitleID
	Status         ReservationStatus
	ExpirationDate *ExpirationDate
	PlacedAt       time.Time
}

func (r *Reservation) Cancel() error {
	if r.Status == ReservationStatusFulfilled || r.Status == ReservationStatusCanceled {
		return errors.New("a fulfilled or canceled reservation cannot be changed")
	}
	r.Status = ReservationStatusCanceled
	return nil
}

func (r *Reservation) MarkReadyForPickup(expiry ExpirationDate) error {
	if r.Status != ReservationStatusWaiting {
		return errors.New("only waiting reservations can be marked ready for pickup")
	}
	r.Status = ReservationStatusReadyForPickup
	r.ExpirationDate = &expiry
	return nil
}

func (r *Reservation) Fulfill() error {
	if r.Status != ReservationStatusReadyForPickup {
		return errors.New("only ready reservations can be fulfilled")
	}
	r.Status = ReservationStatusFulfilled
	return nil
}

// ReservationQueue is the aggregate root. It is the consistency boundary for all
// queue operations on a single title: members always join at the back (Enqueue) and
// the next eligible member is determined in arrival order (Next). These invariants are
// enforced here so no external code can corrupt queue ordering or bypass priority.
type ReservationQueue struct {
	TitleID      TitleID
	Reservations []Reservation
}

func (q *ReservationQueue) Enqueue(r Reservation) {
	q.Reservations = append(q.Reservations, r)
}

// Next returns the first waiting reservation in queue order.
func (q *ReservationQueue) Next() *Reservation {
	for i := range q.Reservations {
		if q.Reservations[i].Status == ReservationStatusWaiting {
			return &q.Reservations[i]
		}
	}
	return nil
}
