FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY main.go .
RUN go mod init best-api && go build -o best-api main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/best-api .
EXPOSE 9999
CMD ["./best-api"]
