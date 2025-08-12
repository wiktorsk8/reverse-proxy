FROM golang:1.24-bookworm AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o reverse-proxy cmd/main.go

FROM gcr.io/distroless/base-debian12 AS prod

WORKDIR /prod

COPY --from=base /app/reverse-proxy .
COPY --from=base /app/internal/config/proxy.yml .
COPY --from=base /app/.env .

USER nonroot:nonroot

EXPOSE 8000

CMD ["./reverse-proxy", "proxy.yml"]
