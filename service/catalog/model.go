// Package catalog manages the library catalog: titles and their physical copies.
// It is the authoritative source for what exists and whether a copy can be borrowed.
package catalog

import (
	"errors"
	"time"
)

type TitleID string
type CopyID string
type ISBN string

type CopyStatus string

const (
	CopyStatusAvailable   CopyStatus = "available"
	CopyStatusOnLoan      CopyStatus = "on_loan"
	CopyStatusReserved    CopyStatus = "reserved"
	CopyStatusLost        CopyStatus = "lost"
	CopyStatusDamaged     CopyStatus = "damaged"
	CopyStatusUnavailable CopyStatus = "unavailable"
)

type Author struct {
	Name string
}

type Publisher struct {
	Name string
}

type Category struct {
	Name string
}

type ShelfLocation struct {
	Section string
	Row     string
}

type CopyCondition struct {
	Level       string
	Description string
}

// Copy is a physical instance of a title.
type Copy struct {
	ID            CopyID
	TitleID       TitleID
	Status        CopyStatus
	Condition     CopyCondition
	ShelfLocation ShelfLocation
}

func (c Copy) IsAvailable() bool {
	return c.Status == CopyStatusAvailable
}

func (c *Copy) MarkOnLoan() error {
	if c.Status != CopyStatusAvailable {
		return errors.New("copy is not available for loan")
	}
	c.Status = CopyStatusOnLoan
	return nil
}

func (c *Copy) MarkAvailable() error {
	if c.Status == CopyStatusLost {
		return errors.New("lost copies cannot be made available")
	}
	c.Status = CopyStatusAvailable
	return nil
}

// Title is the aggregate root. It owns its physical copies.
type Title struct {
	ID           TitleID
	Name         string
	ISBN         ISBN
	Author       Author
	Publisher    Publisher
	Category     Category
	Copies       []Copy
	RegisteredAt time.Time
}

func (t *Title) AddCopy(copy Copy) {
	t.Copies = append(t.Copies, copy)
}

func (t *Title) AvailableCopies() []Copy {
	var available []Copy
	for _, c := range t.Copies {
		if c.IsAvailable() {
			available = append(available, c)
		}
	}
	return available
}

// BorrowableCopy is the read model exposed to the Loan context.
// The Loan context never receives a full Title or Copy.
type BorrowableCopy struct {
	CopyID       CopyID
	TitleID      TitleID
	IsBorrowable bool
}
