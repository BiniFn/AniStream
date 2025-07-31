FROM golang:1.24.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /proxy ./cmd/proxy

FROM alpine:3.20
WORKDIR /app

RUN adduser -D aniways
USER aniways

COPY --from=builder /proxy /app/proxy

CMD ["/app/proxy"]
