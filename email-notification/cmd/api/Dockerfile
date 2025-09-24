FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /api-server .

EXPOSE 9919

CMD ["/api-server"]
