FROM golang:1.25-alpine AS build
RUN apk add --no-cache git ca-certificates
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /build/server ./cmd/server && \
    CGO_ENABLED=0 go build -o /build/ingest ./cmd/ingest

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /build/server /usr/local/bin/server
COPY --from=build /build/ingest /usr/local/bin/ingest
EXPOSE 8080
ENTRYPOINT ["server"]
