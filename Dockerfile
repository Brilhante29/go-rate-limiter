FROM golang:1.26.5-alpine3.24 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test ./... && go vet ./...
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/go-rate-limiter ./cmd/go-rate-limiter

FROM scratch AS app
COPY --from=build /out/go-rate-limiter /go-rate-limiter
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/go-rate-limiter"]
CMD ["serve"]

FROM redis:8.8.0-alpine AS demo
COPY --from=build /out/go-rate-limiter /usr/local/bin/go-rate-limiter
USER redis
ENTRYPOINT ["go-rate-limiter"]
CMD ["demo", "--output", "/tmp/rate-limiter-baseline.json"]
