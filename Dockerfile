FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./client ./client/.
RUN CGO_ENABLED=0 go build -o ./server ./server/.

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/client /app/server ./
CMD ["./server"]
