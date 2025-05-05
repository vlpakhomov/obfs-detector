FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .. .
RUN CGO_ENABLED=0 go build -v -o ./obfs_detector ./cmd

FROM alpine:3.21
COPY --from=builder /app/config.json .
COPY --from=builder /app/obfs_detector .
ENTRYPOINT ["./obfs_detector"]