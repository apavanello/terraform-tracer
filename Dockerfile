# Stage 1: Build Frontend (Vue/Vite)
FROM oven/bun:latest AS ui-builder
WORKDIR /app/ui
# Copy dependency manifests first for caching
COPY ui/package.json ui/bun.lock ./
RUN bun install --frozen-lockfile
# Copy the full UI source and build
COPY ui/ ./
RUN bun run build

# Stage 2: Build Backend (Go)
FROM golang:alpine AS go-builder
WORKDIR /app
# Cache Go module downloads
COPY go.mod go.sum ./
RUN go mod download
# Copy project source
COPY . .
# Inject compiled frontend into the embed directory
RUN rm -rf internal/api/dist
COPY --from=ui-builder /app/ui/dist ./internal/api/dist
# Compile static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tracer ./cmd/tracer

# Stage 3: Minimal runner
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata git
COPY --from=go-builder /app/tracer /tracer
ENTRYPOINT ["/tracer"]
CMD ["start", "/workspace", "--port", "9999"]
EXPOSE 9999
