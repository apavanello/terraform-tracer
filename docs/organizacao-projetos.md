# Organização de Projeto (Monorepo Acoplado)

Em virtude de definirmos a distribuição em pacote único com Go Embarcando (Embedded) os estáticos do Vue, usaremos um *Standard Go Project Layout*. 

O repositório estará organizado assim:

```text
terraform-tracer/
├── cmd/
│   └── tracer/               # Ponto de entrada do CLI (main.go, comandos cobra)
├── internal/
│   ├── parser/               # Core business do parse de HCL (onde o esforço grande estará)
│   ├── api/                  # Endpoints REST e start do Fiber/Gin pra servir json/estáticos
│   └── models/               # Structs de Graph, Nodes, Edges, Environments
├── ui/                       # Frontend SPA (Vue 3 + Vite) gerido inteiramente pelo Bun
│   ├── src/                  # Código Vue, Stores e Assets
│   ├── package.json          # Dependências do frontend
│   └── bun.lockb             # Lockfile do Bun garantindo instâncias ultrarrápidas
├── Makefile / Taskfile       # Scripts para orquestração de Build
├── go.mod                    # Dependências do backend (hashicorp/hcl)
├── go.sum                    # Hashes segurança Go
└── docs/                     # Repositório de documentações, arquiteturas e C4 (aonde estamos!)
```

## Como a build (Pipeline Local/CI) funcionará:
1. **Frontend Build**: O script roda `cd ui && bun install && bun run build`. O Vite colocará todos os HTMLs, CSS e mini-scripts (vis-network, Vue runtime) dentro de `ui/dist`.
2. **Go Embed**: No código do Go (`internal/api/server.go`), usaremos:
   ```go
   //go:embed ../../ui/dist/*
   var uiFiles embed.FS
   ```
3. **Backend Compile**: O script roda `go build -o tracer ./cmd/tracer`. Ele vai compilar o script do backend binário que já carrega o zip comprimido de todo front dentro do próprio corpo (standalone).
4. O DevOps recebe `tracer` e joga em `/usr/local/bin`. Fim. Menor atrito possível!
