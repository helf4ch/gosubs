FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd/api


FROM alpine

WORKDIR /root/

COPY --from=builder /app/main .
COPY .env .

CMD ["./main"]
