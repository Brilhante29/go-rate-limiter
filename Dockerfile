FROM golang:1.26.5-alpine3.24 AS build

WORKDIR /src
RUN apk add --no-cache redis
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN redis-server --save "" --appendonly no --daemonize yes && \
    REDIS_TEST_ADDR=127.0.0.1:6379 go test ./... && \
    redis-cli shutdown nosave && \
    go vet ./...
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/go-rate-limiter ./cmd/go-rate-limiter

FROM scratch AS app
COPY --from=build /out/go-rate-limiter /go-rate-limiter
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/go-rate-limiter"]
CMD ["serve"]

FROM redis:8.8.0-alpine@sha256:9d317178eceac8454a2284a9e6df2466b93c745529947f0cd42a0fa9609d7005 AS demo
COPY --from=build /out/go-rate-limiter /usr/local/bin/go-rate-limiter
USER redis
ENTRYPOINT ["go-rate-limiter"]
CMD ["demo", "--output", "/tmp/rate-limiter-baseline.json"]
