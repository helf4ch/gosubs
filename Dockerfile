FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd/api


FROM scratch

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
