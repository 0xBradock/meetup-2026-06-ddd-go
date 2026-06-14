# httpserver/

HTTP transport layer. This package translates between the HTTP wire format and domain calls — nothing more. **No business logic lives here.**

The patterns in this package are inspired by [How I write HTTP services in Go after 13 years](https://grafana.com/blog/how-i-write-http-services-in-go-after-13-years/) by Mat Ryer.

## Files

| File | Role |
|---|---|
| `server.go` | HTTP server construction and configuration |
| `routes.go` | Route registration — maps URL patterns to handler functions via `addRoutes` |
| `health.go` | Handler functions for the health domain endpoints |
| `responses.go` | Shared helpers: `decode` (read JSON body), `httpResponse` (write JSON response) |

## Handler Pattern

Each handler is a method on a handler struct that returns an `http.HandlerFunc`. This gives each handler its own closure, which means dependencies (logger, context, config) are captured once at registration time — not passed through every request.

```go
func (h *healthHandler) getPing(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Request type — declared right here, local to this handler
        type request struct {
            Message string `json:"message"`
        }

        // 2. Decode and validate input
        req, err := decode[request](r)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // If the HTTP response shape must differ from the domain output,
        // declare a local response type here
        //
        // type response struct {
        //     Msg       string `json:"msg"`
        //     Timestamp int64  `json:"timestamp"`
        // }

        // 3. Call the domain service
        out, err := h.s.Ping(ctx, health.PingInput{Message: req.Message})
        if err != nil {
            if errors.Is(err, health.ErrEmptyPingMessage) {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // 4. Write the response
        httpResponse(w, http.StatusOK, out)

        // If the HTTP response shape must differ from the domain output,
        // return the response object here
        // httpResponse(w, http.StatusOK, response{
        //     Msg:       out.Message,
        //     Timestamp: out.ReceivedAtUnix,
        // })
    }
}
```

Three things to notice:
- The **request type** is declared as a local `type` inside the closure — not in a shared file. It belongs to this handler.
- **Validation** happens right there, before the domain call. If the input is wrong the handler returns immediately.
- The **response** is the domain output passed directly to `httpResponse` when the shapes match. When the HTTP response needs a different shape, declare a local response type between decode and the domain call — same pattern as the request type — then use it at step 4.

## Naming Convention

Handler method names follow the domain method they call, lowercased:

| Domain method | Handler method |
|---|---|
| `service.GetHealth(...)` | `getHealth(...)` |
| `service.Ping(...)` | `getPing(...)` |

This makes it easy to trace from a route back to its domain operation.

## Adding a New Endpoint

1. Add a handler method on the relevant handler struct (or create a new struct if it is a different domain). Follow the closure pattern above: request type declared locally, validate, call service, write response.
2. Register the route in `addRoutes` in `routes.go`.
3. Keep the handler thin. If you find yourself writing conditional logic beyond HTTP error mapping, that logic belongs in `health/service.go`.

## Related

- [docs/architecture.md](../docs/architecture.md) — where `httpserver/` fits as a Driver in the hexagonal layout
- [health/README.md](../health/README.md) — domain service methods and input/output types available to call
