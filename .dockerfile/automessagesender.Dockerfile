FROM golang:1.25.4-trixie AS builder

WORKDIR /build

COPY . .

RUN go mod verify
RUN go build -ldflags "-s -w" -trimpath -o automessagesender ./cmd/

FROM debian:trixie-slim

WORKDIR /app

COPY --from=builder /build/automessagesender /app/automessagesender

CMD ["./automessagesender"]