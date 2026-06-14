package catalog

import (
	"context"
	"log/slog"
)

// Service is the application service for the Catalog context.
// It orchestrates use cases; business rules live in the domain types.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- RegisterTitle ---

type RegisterTitleInput struct {
	Name      string
	ISBN      string
	Author    string
	Publisher string
	Category  string
}

type RegisterTitleOutput struct {
	TitleID TitleID
}

func (s *Service) RegisterTitle(ctx context.Context, in RegisterTitleInput) (RegisterTitleOutput, error) {
	// 1. Validate required fields (name, ISBN)
	// 2. Build Title aggregate
	// 3. repo.SaveTitle(ctx, title)
	// 4. Publish TitleRegistered event
	return RegisterTitleOutput{}, nil
}

// --- AddCopy ---

type AddCopyInput struct {
	TitleID      TitleID
	Condition    string
	ShelfSection string
	ShelfRow     string
}

type AddCopyOutput struct {
	CopyID CopyID
}

func (s *Service) AddCopy(ctx context.Context, in AddCopyInput) (AddCopyOutput, error) {
	// 1. repo.GetTitle(ctx, in.TitleID)
	// 2. Build Copy value with CopyStatusAvailable
	// 3. title.AddCopy(copy)
	// 4. repo.SaveTitle(ctx, title)
	// 5. Publish CopyAddedToCatalog event
	return AddCopyOutput{}, nil
}

// --- GetTitle ---

type GetTitleOutput struct {
	Title Title
}

func (s *Service) GetTitle(ctx context.Context, id TitleID) (GetTitleOutput, error) {
	// 1. repo.GetTitle(ctx, id)
	return GetTitleOutput{}, nil
}
