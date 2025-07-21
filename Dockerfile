FROM golang:1.24.0-alpine AS builder
WORKDIR /marketplace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /marketplace/bin/app ./cmd/app/main.go

FROM scratch
COPY --from=builder /marketplace/bin/app /marketplace/app
COPY --from=builder /marketplace/migrations /marketplace/migrations
WORKDIR /marketplace
CMD ["/marketplace/app"]
