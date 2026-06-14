package overdue

import (
	"context"
	"log/slog"
)

// Service is the application service for the Overdue context.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- OpenCase (reacts to LoanBecameOverdue event) ---

type OpenCaseInput struct {
	LoanID   LoanID
	MemberID MemberID
}

type OpenCaseOutput struct {
	CaseID OverdueCaseID
}

func (s *Service) OpenCase(ctx context.Context, in OpenCaseInput) (OpenCaseOutput, error) {
	// 1. Validate LoanID and MemberID are non-empty
	// 2. Build OverdueCase aggregate with status Open
	// 3. repo.SaveOverdueCase(ctx, overdueCase)
	// 4. Publish LoanBecameOverdue event
	return OpenCaseOutput{}, nil
}

// --- ApplyFine ---

type ApplyFineInput struct {
	CaseID      OverdueCaseID
	DaysOverdue int
}

type ApplyFineOutput struct {
	FineID FineID
	Amount FineAmount
}

func (s *Service) ApplyFine(ctx context.Context, in ApplyFineInput) (ApplyFineOutput, error) {
	// 1. repo.GetOverdueCase(ctx, in.CaseID)
	// 2. repo.GetFinePolicy(ctx)
	// 3. policy.Calculate(in.DaysOverdue) — pure value computation
	// 4. Build Fine value
	// 5. overdueCase.ApplyFine(fine) — enforces state transition
	// 6. repo.SaveOverdueCase(ctx, overdueCase)
	// 7. Publish FineApplied event
	return ApplyFineOutput{}, nil
}

// --- WaiveFine ---

type WaiveFineInput struct {
	CaseID OverdueCaseID
	Reason string
}

func (s *Service) WaiveFine(ctx context.Context, in WaiveFineInput) error {
	// 1. repo.GetOverdueCase(ctx, in.CaseID)
	// 2. overdueCase.WaiveFine() — enforces pending fine exists
	// 3. repo.SaveOverdueCase(ctx, overdueCase)
	// 4. Publish FineWaived event
	return nil
}

// --- ResolveCase ---

type ResolveCaseInput struct {
	CaseID OverdueCaseID
}

func (s *Service) ResolveCase(ctx context.Context, in ResolveCaseInput) error {
	// 1. repo.GetOverdueCase(ctx, in.CaseID)
	// 2. overdueCase.Resolve() — enforces fine is paid or waived
	// 3. repo.SaveOverdueCase(ctx, overdueCase)
	// 4. Publish OverdueCaseResolved event
	return nil
}

// --- ListOpenCases ---

type ListOpenCasesOutput struct {
	Cases []OverdueCase
}

func (s *Service) ListOpenCases(ctx context.Context) (ListOpenCasesOutput, error) {
	// 1. repo.ListOpenCases(ctx)
	return ListOpenCasesOutput{}, nil
}
