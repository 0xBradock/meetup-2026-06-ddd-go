package overdue

import "context"

// Repository is the persistence port for the Overdue context.
type Repository interface {
	GetOverdueCase(ctx context.Context, id OverdueCaseID) (OverdueCase, error)
	GetFinePolicy(ctx context.Context) (FinePolicy, error)
	ListOpenCases(ctx context.Context) ([]OverdueCase, error)
	SaveOverdueCase(ctx context.Context, oc OverdueCase) error
}
