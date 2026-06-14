package reservation

import "context"

// Repository is the persistence port for the Reservation context.
type Repository interface {
	GetReservation(ctx context.Context, id ReservationID) (Reservation, error)
	GetQueue(ctx context.Context, titleID TitleID) (ReservationQueue, error)
	SaveReservation(ctx context.Context, r Reservation) error
	SaveQueue(ctx context.Context, q ReservationQueue) error
	ListMemberReservations(ctx context.Context, memberID MemberID) ([]Reservation, error)
}
