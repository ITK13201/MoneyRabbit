ARG GOOSE_VERSION=v3.28.0

FROM golang:1.26-alpine AS builder
ARG GOOSE_VERSION
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server \
    && CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@${GOOSE_VERSION}

FROM builder AS development
RUN wget -qO /usr/local/bin/atlas https://release.ariga.io/atlas/atlas-linux-amd64-latest \
    && chmod +x /usr/local/bin/atlas
EXPOSE 8080
ENTRYPOINT ["/app/server"]

FROM gcr.io/distroless/static:debug AS production
ENV GIN_MODE=release
COPY --from=builder /app/server /server
COPY --from=builder /app/db/migrations /db/migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
EXPOSE 8080
ENTRYPOINT ["/server"]
