FROM golang:1.25.4-trixie AS builder

WORKDIR /build

COPY . .

RUN go mod verify
RUN go build -ldflags "-s -w" -trimpath -o automessagesender ./cmd/

FROM debian:trixie-slim

RUN apt-get update && apt-get install -y curl \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/automessagesender /app/automessagesender

EXPOSE 8080

CMD ["./automessagesender"]