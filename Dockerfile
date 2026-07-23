# syntax=docker/dockerfile:1

FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
COPY migrations/ migrations/
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api
RUN go build -trimpath -ldflags="-s -w" -o /out/import-curriculum ./cmd/import-curriculum
RUN go build -trimpath -ldflags="-s -w" -o /out/import-lessons ./cmd/import-lessons

FROM alpine:3.21 AS final
RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/api /out/import-curriculum /out/import-lessons ./
# Migrations are also go:embed'd into the api binary; the .sql files are
# copied here too so import-curriculum/import-lessons's -dir default and any
# manual inspection on the server have them on disk. migrations.go itself
# (source) is deliberately not copied.
COPY migrations/*.sql ./migrations/
COPY content/ ./content/
RUN chown -R app:app /app
USER app
EXPOSE 8080
ENTRYPOINT ["/app/api"]
