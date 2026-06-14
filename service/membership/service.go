package membership

import (
	"context"
	"log/slog"
)

// Service is the application service for the Membership context.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- RegisterMember ---

type RegisterMemberInput struct {
	Name  string
	Email string
	Phone string
}

type RegisterMemberOutput struct {
	MemberID MemberID
}

func (s *Service) RegisterMember(ctx context.Context, in RegisterMemberInput) (RegisterMemberOutput, error) {
	// 1. Validate name and email are non-empty
	// 2. Build Member aggregate with MembershipStatusActive and BorrowingStatusAllowed
	// 3. repo.SaveMember(ctx, member)
	// 4. Publish MemberRegistered event
	return RegisterMemberOutput{}, nil
}

// --- GetMember ---

type GetMemberOutput struct {
	Member Member
}

func (s *Service) GetMember(ctx context.Context, id MemberID) (GetMemberOutput, error) {
	// 1. repo.GetMember(ctx, id)
	return GetMemberOutput{}, nil
}

// --- BlockMember (reacts to FineApplied event) ---

type BlockMemberInput struct {
	MemberID MemberID
}

func (s *Service) BlockMember(ctx context.Context, in BlockMemberInput) error {
	// 1. repo.GetMember(ctx, in.MemberID)
	// 2. member.Block() — enforces active membership
	// 3. repo.SaveMember(ctx, member)
	// 4. Publish MemberBlocked event
	return nil
}

// --- UnblockMember (reacts to FinePaid or FineWaived events) ---

type UnblockMemberInput struct {
	MemberID MemberID
}

func (s *Service) UnblockMember(ctx context.Context, in UnblockMemberInput) error {
	// 1. repo.GetMember(ctx, in.MemberID)
	// 2. member.Unblock() — enforces blocked status
	// 3. repo.SaveMember(ctx, member)
	// 4. Publish MemberUnblocked event
	return nil
}

// --- SuspendMember ---

type SuspendMemberInput struct {
	MemberID MemberID
}

func (s *Service) SuspendMember(ctx context.Context, in SuspendMemberInput) error {
	// 1. repo.GetMember(ctx, in.MemberID)
	// 2. member.Suspend() — enforces active membership
	// 3. repo.SaveMember(ctx, member)
	// 4. Publish MembershipSuspended event
	return nil
}
