# Arquitetura de Componentes

Este documento descreve a arquitetura do sistema **Terraform Tracer** utilizando o modelo C4 (Contexto e Containers). A decisão arquitetural atualizada foca em um **Monorepo Acoplado**, utilizando **Golang** no backend (para aproveitar bibliotecas nativas como `hashicorp/hcl`) e **Vue.js** empacotado pelo **Bun** no frontend.

## Diagrama de Contexto (C4 Nível 1)

```mermaid
C4Context
    title Diagrama de Contexto - Terraform Tracer

    Person(dev, "DevOps / SRE", "Engenheiro que necessita analisar dependências de infraestrutura em código.")
    System(tracer, "Terraform Tracer", "Ferramenta CLI/Local que lê um diretório e exibe o grafo visual das dependências.")
    System_Ext(tfFiles, "Arquivos Terraform Locais", "Pasta contendo .tf, .tfvars, modules, etc.")

    Rel(dev, tracer, "Executa CLI e interage pela UI Web", "localhost:8080")
    Rel(tracer, tfFiles, "Lê e faz parse do código estático (HCL)", "File System")
```

## Diagrama de Containers (C4 Nível 2)

```mermaid
C4Container
    title Diagrama de Container - Acoplamento Backend (Go) + Frontend (Vue)

    Person(dev, "DevOps / SRE", "Usuário executando a ferramenta")

    System_Boundary(tracer_system, "Terraform Tracer") {
        Container(cli, "Go CLI Entrypoint", "Golang", "Inicia o parse, sobe o web server e (opcionalmente) abre o browser.")
        Container(hcl_parser, "HCL Parser Engine", "Golang (hashicorp/hcl/v2)", "Motor nativo que lê os arquivos locais e monta a árvore Abstract Syntax Tree (AST).")
        Container(api_server, "API Server", "Golang (net/http ou Fiber)", "Fornece os endpoints REST com o grafo JSON para a UI.")
        Container(spa, "Single Page App", "Vue.js + Vite + Bun", "Interface rica. O build gerado é embarcado (embedded) no binário do Go.")
    }

    ContainerDb(file_sys, "Local File System", "Disk", "Arquitetura stateless, lê os .tf na hora.")

    Rel(dev, cli, "Dispara comando", "CLI")
    Rel(cli, api_server, "Inicia servidor web")
    Rel(cli, hcl_parser, "Chama parser em background")
    Rel(hcl_parser, file_sys, "Lê arquivos .tf", "I/O")
    
    Rel(dev, spa, "Acessa dashboard", "HTTPS/Browser")
    Rel(spa, api_server, "Requisita dados do grafo /api/v1/graph", "JSON/REST")
    Rel(api_server, spa, "Serve arquivos estáticos embutidos", "HTTP")
```

## Decisões Arquiteturais:
1. **Linguagem Backend (Golang):** Substituímos o Python por Go. O ecossistema Go é nativo da Hashicorp. Usar a library `hashicorp/hcl/v2` elimina a precisão frágil de *Regex* e a necessidade de reinventar a roda em Python.
2. **Frontend (Vue via Bun):** O frontend será reativo, utilizando o Vue.js por sua clareza e separação limpa de componentes. O gerenciador de pacotes e empacotador será o **Bun**, para builds ultra-rápidos e leves. 
3. **Distribuição (Embed):** Com as funcionalidades do Go 1.16+ (`go:embed`), a pasta de build do *Vue* (`dist/`) poderá ser compilada dentro de um único arquivo executável binário final. Isso é excelente para a distribuição (o SRE só precisa baixar um arquivo `tracer-linux-amd64` e rodar).
