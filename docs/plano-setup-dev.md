# Setup do Ambiente de Desenvolvimento

Este plano foca em deixar a máquina do engenheiro pronta para codificar, testar e empacotar o Terraform Tracer.

- [ ] **Configuração do Repositório Base**
  - [ ] Executar `git init` e criar estrutura de pastas (`cmd`, `internal`, `ui`, `docs`).
  - [ ] Criar `.gitignore` para Go e Node/Bun (`/ui/node_modules`, `tracer`, `*.exe`).
- [ ] **Setup do Backend (Golang)**
  - [ ] Inicializar módulo (`go mod init github.com/sua-org/terraform-tracer`).
  - [ ] Adicionar dependência principal Hashicorp: `go get github.com/hashicorp/hcl/v2`.
  - [ ] Adicionar framework web leve (se aplicável), ex: `go get github.com/gofiber/fiber/v2`.
  - [ ] Configurar Linter: Instalar e configurar `golangci-lint`.
- [ ] **Setup do Frontend (Vue + Bun)**
  - [ ] Entrar na pasta `ui/` e rodar o bootstrap do Vite+Vue via Bun: `bun create vite . --template vue-ts`.
  - [ ] Instalar dependências de UI: `bun add vis-network @tabler/icons-vue`.
  - [ ] Configurar as variáveis de ambiente base (.env local apontando para `:8080`).
- [ ] **Scripts de Desenvolvimento (Makefile / Taskfile)**
  - [ ] Criar comando `make dev-ui` que executa `cd ui && bun run dev`.
  - [ ] Criar comando `make dev-api` que sobe `go run ./cmd/tracer start`.
  - [ ] Criar comando `make build` que empacota o front com `bun run build` e em seguida `go build`.
