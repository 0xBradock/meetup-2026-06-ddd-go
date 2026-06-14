package membership

import "context"

// Repository is the persistence port for the Membership context.
type Repository interface {
	GetMember(ctx context.Context, id MemberID) (Member, error)
	SaveMember(ctx context.Context, member Member) error
	ListMembers(ctx context.Context) ([]Member, error)
}
