package reservation

import (
	"context"
	"log/slog"
)

// Service is the application service for the Reservation context.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// --- PlaceReservation ---

type PlaceReservationInput struct {
	MemberID MemberID
	TitleID  TitleID
}

type PlaceReservationOutput struct {
	ReservationID ReservationID
}

func (s *Service) PlaceReservation(ctx context.Context, in PlaceReservationInput) (PlaceReservationOutput, error) {
	// 1. Load ReservableTitle — ACL call to Catalog to verify no copies are available
	// 2. if reservableTitle.IsReservable == false — return error (copies still available)
	// 3. repo.GetQueue(ctx, in.TitleID)
	// 4. Build Reservation aggregate with status Waiting
	// 5. queue.Enqueue(reservation)
	// 6. repo.SaveReservation(ctx, reservation)
	// 7. repo.SaveQueue(ctx, queue)
	// 8. Publish ReservationPlaced event
	return PlaceReservationOutput{}, nil
}

// --- CancelReservation ---

type CancelReservationInput struct {
	ReservationID ReservationID
}

func (s *Service) CancelReservation(ctx context.Context, in CancelReservationInput) error {
	// 1. repo.GetReservation(ctx, in.ReservationID)
	// 2. reservation.Cancel() — enforces terminal state check
	// 3. repo.SaveReservation(ctx, reservation)
	// 4. Publish ReservationCanceled event
	return nil
}

// --- NotifyReadyForPickup (reacts to CopyReturned or CopyMarkedAvailable events) ---

type NotifyReadyInput struct {
	TitleID       TitleID
	ExpiresInDays int
}

func (s *Service) NotifyReadyForPickup(ctx context.Context, in NotifyReadyInput) error {
	// 1. repo.GetQueue(ctx, in.TitleID)
	// 2. queue.Next() — get first waiting reservation
	// 3. Build ExpirationDate value object (today + ExpiresInDays)
	// 4. reservation.MarkReadyForPickup(expiry) — enforces waiting status
	// 5. repo.SaveReservation(ctx, reservation)
	// 6. Publish ReservationReadyForPickup event
	return nil
}

// --- GetQueue ---

type GetQueueOutput struct {
	Queue ReservationQueue
}

func (s *Service) GetQueue(ctx context.Context, titleID TitleID) (GetQueueOutput, error) {
	// 1. repo.GetQueue(ctx, titleID)
	return GetQueueOutput{}, nil
}
