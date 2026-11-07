FROM golang:1.25.5 AS builder

WORKDIR /app
COPY . .

# RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o NABARI main.go

#---
FROM scratch

WORKDIR /app
COPY --from=builder /app/NABARI .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# distroless-style non-root uid:gid。scratch なので /etc/passwd は無く数値で指定。
USER 65532:65532

EXPOSE 34165
CMD ["./NABARI"]
