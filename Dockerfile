FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o petsite .


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/petsite .

COPY static ./static

RUN mkdir -p log

EXPOSE 6060

CMD ["./petsite"]