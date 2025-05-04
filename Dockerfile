FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o person-enricher ./cmd/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/person-enricher .

COPY .env .

EXPOSE 8080

ENTRYPOINT ["./person-enricher"]
