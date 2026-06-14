package loan

import (
	"context"
	"log/slog"
	"time"
)

// Service is the application service for the Loan context.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- BorrowBook ---

type BorrowBookInput struct {
	MemberID MemberID
	CopyID   CopyID
}

type BorrowBookOutput struct {
	LoanID  LoanID
	DueDate time.Time
}

func (s *Service) BorrowBook(ctx context.Context, in BorrowBookInput) (BorrowBookOutput, error) {
	// 1. repo.GetBorrower(ctx, in.MemberID)
	// 2. repo.GetBorrowableCopy(ctx, in.CopyID)
	// 3. borrower.CanBorrow() — enforces eligibility and loan limit
	// 4. if !borrowableCopy.IsBorrowable — return error
	// 5. Build Loan with DueDate = today + 21 days
	// 6. repo.SaveLoan(ctx, loan)
	// 7. Publish LoanCreated and CopyBorrowed events
	return BorrowBookOutput{}, nil
}

// --- ReturnCopy ---

type ReturnCopyInput struct {
	LoanID LoanID
}

type ReturnCopyOutput struct {
	LoanID LoanID
}

func (s *Service) ReturnCopy(ctx context.Context, in ReturnCopyInput) (ReturnCopyOutput, error) {
	// 1. repo.GetLoan(ctx, in.LoanID)
	// 2. loan.Return() — enforces active status
	// 3. repo.SaveLoan(ctx, loan)
	// 4. Publish CopyReturned event
	return ReturnCopyOutput{}, nil
}

// --- ExtendLoan ---

type ExtendLoanInput struct {
	LoanID     LoanID
	ExtendDays int
}

type ExtendLoanOutput struct {
	LoanID     LoanID
	NewDueDate time.Time
}

func (s *Service) ExtendLoan(ctx context.Context, in ExtendLoanInput) (ExtendLoanOutput, error) {
	// 1. repo.GetLoan(ctx, in.LoanID)
	// 2. Compute new due date from current DueDate + ExtendDays
	// 3. NewDueDate(newDate) — value object validates future date
	// 4. loan.Extend(newDueDate) — enforces active status
	// 5. repo.SaveLoan(ctx, loan)
	// 6. Publish LoanExtended event
	return ExtendLoanOutput{}, nil
}

// --- GetLoan ---

type GetLoanOutput struct {
	Loan Loan
}

func (s *Service) GetLoan(ctx context.Context, id LoanID) (GetLoanOutput, error) {
	// 1. repo.GetLoan(ctx, id)
	return GetLoanOutput{}, nil
}
