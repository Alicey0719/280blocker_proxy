FROM golang:1.24.3 AS builder

WORKDIR /app
COPY . .

# RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o NABARI main.go

#---
FROM scratch

WORKDIR /app
COPY --from=builder /app/NABARI .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 34165
CMD ["./NABARI"]
