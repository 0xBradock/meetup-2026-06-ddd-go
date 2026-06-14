package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"
)

// addRoutes is the function responsible to map all the routes and its handlers.
func addRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	logger *slog.Logger,
	hh healthHandler,
	ch catalogHandler,
	lh loanHandler,
	oh overdueHandler,
	rh reservationHandler,
	mh membershipHandler,
) {
	// Health
	mux.Handle("/health", hh.getHealth(ctx, logger))
	mux.Handle("/ping", hh.getPing(ctx, logger))

	// Catalog — title and copy management
	mux.Handle("POST /catalog/titles", ch.registerTitle(ctx, logger))
	mux.Handle("GET /catalog/titles/{id}", ch.getTitle(ctx, logger))
	mux.Handle("POST /catalog/titles/{id}/copies", ch.addCopy(ctx, logger))

	// Loan — borrowing, returning, and extending copies
	mux.Handle("POST /loan/borrow", lh.borrowBook(ctx, logger))
	mux.Handle("POST /loan/{id}/return", lh.returnCopy(ctx, logger))
	mux.Handle("POST /loan/{id}/extend", lh.extendLoan(ctx, logger))
	mux.Handle("GET /loan/{id}", lh.getLoan(ctx, logger))

	// Overdue — overdue detection and fine management
	mux.Handle("GET /overdue/cases", oh.listOpenCases(ctx, logger))
	mux.Handle("POST /overdue/cases/{id}/fine", oh.applyFine(ctx, logger))
	mux.Handle("POST /overdue/cases/{id}/waive", oh.waiveFine(ctx, logger))
	mux.Handle("POST /overdue/cases/{id}/resolve", oh.resolveCase(ctx, logger))

	// Reservation — holds on titles when no copies are available
	mux.Handle("POST /reservation", rh.placeReservation(ctx, logger))
	mux.Handle("DELETE /reservation/{id}", rh.cancelReservation(ctx, logger))
	mux.Handle("GET /reservation/queue/{titleID}", rh.getQueue(ctx, logger))

	// Membership — members and borrowing eligibility
	mux.Handle("POST /membership/members", mh.registerMember(ctx, logger))
	mux.Handle("GET /membership/members/{id}", mh.getMember(ctx, logger))
	mux.Handle("POST /membership/members/{id}/block", mh.blockMember(ctx, logger))
	mux.Handle("POST /membership/members/{id}/unblock", mh.unblockMember(ctx, logger))

	// Profiling
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
