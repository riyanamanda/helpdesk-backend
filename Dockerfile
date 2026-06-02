FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed

FROM alpine:latest

ARG TARGETARCH
ARG GOOSE_VERSION=v3.24.2

RUN apk --no-cache add ca-certificates tzdata curl tar

RUN case "$TARGETARCH" in \
			amd64) GOOSE_ARCH="x86_64" ;; \
			arm64) GOOSE_ARCH="arm64" ;; \
			*) echo "unsupported TARGETARCH: $TARGETARCH" && exit 1 ;; \
		esac && \
		curl -fsSL -o /tmp/goose.tar.gz "https://github.com/pressly/goose/releases/download/${GOOSE_VERSION}/goose_linux_${GOOSE_ARCH}.tar.gz" && \
		tar -xzf /tmp/goose.tar.gz -C /usr/local/bin goose && \
		chmod +x /usr/local/bin/goose && \
		rm -f /tmp/goose.tar.gz

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/seed .
COPY --from=builder /app/migrations ./migrations
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
