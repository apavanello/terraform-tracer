# Terraform Tracer

A visual dependency tracer for Terraform resources. Select a directory, pick a resource, and see its environments, modules, and dependent resources rendered as an interactive graph.

## Quick Start

```bash
# Build everything
make build

# Run the tracer on a terraform directory
./tracer start ./path/to/terraform
```

## Development

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- [Bun](https://bun.sh/)

### Dev Mode
```bash
# Start the API server (with hot-reload via air if installed)
make dev-api

# In another terminal, start the Vue dev server
make dev-ui
```

## Project Structure
```
terraform-tracer/
├── cmd/tracer/       # CLI entrypoint (Cobra)
├── internal/
│   ├── parser/       # HCL parse engine (hashicorp/hcl/v2)
│   ├── api/          # HTTP server + REST endpoints
│   └── models/       # Graph, Node, Edge, Environment structs
├── ui/               # Vue 3 + Vite frontend (managed by Bun)
├── examples/         # Sample Terraform projects for testing
└── docs/             # Architecture & planning docs
```

## License
MIT
