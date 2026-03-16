# Estrutura Core Funcional

Este plano define as funcionalidades mandatórias que darão sustentação a todo o projeto e farão a integração ponta-a-ponta funcionar antes mesmo de parsearmos lógica complexa do HCL.

- [ ] **Configuração da Estrutura Base (Backend e CLI)**
  - [ ] Implementar o pacote `cmd/tracer` usando [Cobra](https://github.com/spf13/cobra) para gerenciar os argumentos CLI (ex: rodar `tracer start . --port 9000`).
  - [ ] Criar o pacote `internal/api` e inicializar o servidor Web (Fiber ou standard `net/http`).
  - [ ] Configurar o arquivo `server.go` para permitir chamadas locais (CORS se necessário) e criar rotas de health-check (`/ping`).
- [ ] **Integração Go Embed (Backend servindo Frontend)**
  - [ ] Compilar uma build estática *dummy* do Vue (`/ui/dist`).
  - [ ] No `internal/api`, configurar a diretiva `//go:embed` para ler o FS comprimido.
  - [ ] Criar o *handler* HTTP que intercepta requisições não listadas na API e devolve o `index.html` do SPA.
- [ ] **Base da Single Page Application (UI)**
  - [ ] Limpar o scaffold do Vue gerado pelo Vite.
  - [ ] Criar estrutura base de componentes: `components/Sidebar.vue`, `components/GraphView.vue` e `components/PropsPanel.vue`.
  - [ ] Criar a `store` (se usar Pinia ou ref globais) para manter o Estado: *Graph Nodes*, *Selected Node*, e estado de *Loading*.
- [ ] **Comunicação Cliente-Servidor**
  - [ ] No Vue, criar uma camada de service (`services/api.ts`) que fará o fetch nativo para o path `/api/v1/graph`.
  - [ ] Tratar cenários onde o backend retorna `202 Accepted` (Parse ainda rodando) ou `500 Internal Error`.
  - [ ] Desenhar e fixar o Type/Strucuture que a store do Vue deve receber baseado no nosso arquivo `docs/modelo-dados.md`.
