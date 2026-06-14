package catalog

import "context"

// Repository is the persistence port for the Catalog context.
// Implementations live in the infrastructure layer; never here.
type Repository interface {
	GetTitle(ctx context.Context, id TitleID) (Title, error)
	ListTitles(ctx context.Context) ([]Title, error)
	SaveTitle(ctx context.Context, title Title) error
	GetBorrowableCopy(ctx context.Context, id CopyID) (BorrowableCopy, error)
	UpdateCopyStatus(ctx context.Context, id CopyID, status CopyStatus) error
}
