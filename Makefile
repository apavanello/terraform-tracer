.PHONY: dev-ui dev-api build clean

# ── Development ──────────────────────────────────────────────
dev-ui:
	cd ui && bun run dev

dev-api:
	go run ./cmd/tracer start ./examples/simple-aws

# ── Production Build ─────────────────────────────────────────
build: build-ui build-go

build-ui:
	cd ui && bun install && bun run build
	rm -rf internal/api/dist
	cp -r ui/dist internal/api/dist

build-go: build-ui
	go build -o tracer ./cmd/tracer

# ── Housekeeping ─────────────────────────────────────────────
clean:
	rm -f tracer
	rm -rf internal/api/dist
	mkdir -p internal/api/dist
	echo '<html><body>placeholder</body></html>' > internal/api/dist/index.html

tidy:
	go mod tidy
	cd ui && bun install
