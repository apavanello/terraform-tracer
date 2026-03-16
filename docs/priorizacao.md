# Priorização de Features (MoSCoW)

Abaixo estão listadas as estimativas de esforço e a classificação MoSCoW para as features descritas no projeto **Terraform Tracer**.

## Critérios
- **M (Must Have)**: Funcionalidades essenciais para que o MVP funcione e resolva o problema principal.
- **S (Should Have)**: Funcionalidades muito importantes, mas o sistema pode sobreviver ao primeiro dia sem elas.
- **C (Could Have)**: Funcionalidades que agregam valor e experiência de uso, mas só devem ser feitas se houver tempo hábil.
- **W (Won't Have / Would Like)**: Fora do escopo do atual MVP, ficarão para versões futuras.

- **Esforços**: Simples, Médio, Complexo.

## Tabela de Funcionalidades

| Feature | Descrição Curta | Escopo / Tipo | MoSCoW | Esforço | Skills Necessárias |
| :--- | :--- | :--- | :---: | :---: | :--- |
| **F01** | Parseamento de Recursos Locais | Backend (Regra de Negócio) | **Must** | Simples | Golang (hashicorp/hcl/v2) |
| **F02** | Parseamento de Módulos Locais | Backend (Regra de Negócio) | **Must** | Simples | Golang |
| **F03** | Identificação de Dependências (Edges) | Backend (Regra de Negócio) | **Must** | Médio | Golang, Teoria de Grafos |
| **F04** | Parseamento de Environments | Backend (Regra de Negócio) | **Should** | Simples | Golang, Estrutura de dados |
| **F05** | Upload/Seleção de Diretório Local | Frontend / Integração API | **Must** | Simples | Javascript/Vue, Vite, Bun |
| **F06** | Visualização em Grafo Interativo | Frontend (UX/UI) | **Must** | Médio | Vue, vis-network |
| **F07** | Painel de Detalhes do Recurso | Frontend (UX/UI) | **Must** | Simples | Vue, Stores/State |
| **F08** | Árvore de Arquivos (File Explorer) | Frontend (UX/UI) | **Should** | Simples | Vue, CSS |
| **F09** | Provider Git (Clonar repo privado) | Backend (Integração) | **Could** | Complexo | Golang, APIs do Github, SecOps |
| **F10** | Suporte a Monorepos | Backend / Frontend | **Won't** | Complexo | Golang, JS |

## Tabela de Requisitos Não-Funcionais (Arquitetura)

| RNF | Descrição Curta | Escopo / Tipo | MoSCoW | Esforço |
| :--- | :--- | :--- | :---: | :---: |
| **RNF01** | Ferramenta Local / CLI / Desktop | Infraestrutura / Deploy | **Must** | Simples |
| **RNF02** | Processamento Síncrono Padrão | Performance | **Must** | Simples |
| **RNF03** | Cache Temporário em Memória (Local) | Armazenamento / Estado | **Should** | Médio |
| **RNF04** | Sistema Aberto Local (Sem Auth) | Segurança / Rede | **Must** | Simples |

### Resumo Técnico Atualizado

Conforme evoluído e validado em revisões de arquitetura:
- **Backend**: Será construído em **Golang**. O ecossistema Go nos permite embutir a biblioteca oficial *hashicorp/hcl/v2*, fornecendo extrema resiliência e facilidade no parse (diminuindo o Esforço de recursos de *Médio/Complexo* para *Simples/Médio*, já que não precisamos reinventar o AST). O Build distribuirá um arquivo binário único (embedded loader).
- **Frontend**: Utilizaremos **Vue.js + Vite**. O pacote final será servido unicamente via `bun run build`. O uso do Vue permite a melhor manipulação e reatividade do grafo e janelas modais lateral/dashboard em uma Single Page Application sólida.
