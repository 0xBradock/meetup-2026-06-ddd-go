package loan

import "context"

// Repository is the persistence port for the Loan context.
type Repository interface {
	GetLoan(ctx context.Context, id LoanID) (Loan, error)
	GetBorrower(ctx context.Context, id MemberID) (Borrower, error)
	GetBorrowableCopy(ctx context.Context, id CopyID) (BorrowableCopy, error)
	SaveLoan(ctx context.Context, loan Loan) error
	ListActiveLoansByMember(ctx context.Context, memberID MemberID) ([]Loan, error)
}
