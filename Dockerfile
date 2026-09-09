FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o healthcheck ./cmd/healthcheck

FROM gcr.io/distroless/static-debian12

COPY --from=builder /build/server /server
COPY --from=builder /build/healthcheck /healthcheck

ENV DB_PATH=/data/status.db
EXPOSE 8080

ENTRYPOINT ["/server"]
