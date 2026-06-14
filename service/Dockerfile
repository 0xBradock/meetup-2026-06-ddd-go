FROM golang:1.26 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/task_alive ./cmd/task_alive

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /out/server /app/server
COPY --from=build /out/task_alive /app/task_alive

USER nonroot

EXPOSE 8080 9090

ENV DRIVER=inmemory
ENV HTTP_ADDR=:8080

ENTRYPOINT ["/app/server"]
