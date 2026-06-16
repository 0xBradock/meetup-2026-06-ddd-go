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

// Author is a value object. In catalog terms two authors with the same name are the
// same author — no separate identity is required. Wrapping the name in a struct prevents
// primitive obsession and makes the domain language explicit in method signatures.
type Author struct {
	Name string
}

// Publisher is a value object. Equality is by name alone; no identity beyond content.
type Publisher struct {
	Name string
}

// Category is a value object. It classifies a title by subject; two categories with
// the same name represent the same classification, so no identity is needed.
type Category struct {
	Name string
}

// ShelfLocation is a value object describing where a copy sits in the library.
// Two copies at the same Section and Row share the same location concept; the location
// has no lifecycle or identity of its own.
type ShelfLocation struct {
	Section string
	Row     string
}

// CopyCondition is a value object describing the physical state of a copy at a point
// in time. It is purely descriptive — condition is always set on a Copy, never addressed
// independently, and two conditions with the same level and description are identical.
type CopyCondition struct {
	Level       string
	Description string
}

// Copy is an entity. It represents one physical book in the library, has its own identity
// (CopyID), and transitions through availability states (available → on_loan → available).
// It lives inside the Title aggregate — external code accesses it through Title, which
// ensures the catalog's copy inventory is always consistent.
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

// Title is the aggregate root — the consistency boundary for the catalog entry and its
// physical copies. Copies are added through Title.AddCopy and availability is queried
// through Title.AvailableCopies, so the catalog can never report a copy available that
// does not belong to this title, or allow direct copy mutation that bypasses title rules.
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
