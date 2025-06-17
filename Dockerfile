FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd
RUN ls -lh /app/main

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .
RUN ls -lh ./main
RUN chmod +x ./main
CMD ["./main"]