FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o booking ./cmd/booking/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/booking .

CMD [ "./booking", "-c", "config.yml" ]