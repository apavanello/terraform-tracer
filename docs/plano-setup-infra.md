# Setup de Infra / CI / CD e Entrega

Como o MVP foca em ser uma ferramenta local (Estilo CLI rodando no loopback), nossa "Infraestrutura" se baseará pesadamente no release e no pipeline de distribuição dos binários do Golang para as máquinas dos DevOps.

- [ ] **Configuração de Pipeline CI (Github Actions)**
  - [ ] Criar `.github/workflows/ci.yml`.
  - [ ] Adicionar "Job: Linting" (golangci-lint).
  - [ ] Adicionar "Job: Tests" (`go test ./...` para o backend unit test).
  - [ ] Adicionar "Job: UI Build Test" (`bun run build` para checar se gera dist limpo).
- [ ] **Pipeline de Release Automático (CD)**
  - [ ] Criar `.github/workflows/release.yml` acionado apenas em [Tags].
  - [ ] Compilar *Cross-Platform* binaries:
    - [ ] Linux amd64/arm64.
    - [ ] macOS (Darwin) amd64/arm64.
    - [ ] Windows amd64.
  - [ ] Empacotar os binários e publicar no "Github Releases" permitindo que o usuário dê um simples `wget` e execute.
- [ ] **Integração de Qualidade (Opcional Future-Proof)**
  - [ ] Adicionar *SonarCloud* ou *CodeClimate* na verificação do `go.mod`.
  - [ ] Definir versão semântica do projeto (ex: `v0.1.0-beta`).
- [ ] **Ambientes e QA**
  - [ ] Criar uma pasta `/examples/` na raiz do repositório contendo casos conhecidos de projetos Terraform mockados (VPC simples, EKS com Locals, Monorepo HCL). Esta pasta será vital para testes automatizados unitários no Go Parser.
