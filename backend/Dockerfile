FROM golang:1.26-alpine AS builder
WORKDIR /app
# Download atlas CLI
RUN wget -qO /usr/local/bin/atlas https://release.ariga.io/atlas/atlas-linux-amd64-latest \
    && chmod +x /usr/local/bin/atlas
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server \
    && CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

FROM gcr.io/distroless/static:debug
COPY --from=builder /app/server /server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /usr/local/bin/atlas /usr/local/bin/atlas
EXPOSE 8080
ENTRYPOINT ["/server"]
