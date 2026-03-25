FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build
COPY api/ ./api/
WORKDIR /build/api
RUN CGO_ENABLED=1 go build -o /quran-api ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /quran-api /usr/local/bin/quran-api

# Copy frontend files
COPY quran.html /app/quran.html
COPY reset.css /app/reset.css
COPY styles.css /app/styles.css

ENV PORT=8080
ENV DB_PATH=/data/quran.db
ENV STATIC_DIR=/app

EXPOSE 8080
CMD ["quran-api"]
