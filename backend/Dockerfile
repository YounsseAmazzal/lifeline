
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/server ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app/backend

COPY --from=builder /out/server ./server
COPY assets ./assets

RUN mkdir -p ./uploads

EXPOSE 8080

CMD ["./server"]
