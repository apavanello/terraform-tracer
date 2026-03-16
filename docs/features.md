# Levantamento de Requisitos e Funcionalidades (Features)

Este documento descreve as funcionalidades principais para o produto **Terraform Tracer**, baseado no modelo MVP definido no documento de intenção (`docs/intencao.md`).

## 1. Módulo Parser (Análise Estática de Código)

O core da ferramenta precisa ser capaz de ler e entender arquivos `.tf` e `.tfvars`.
- **F01 - Parseamento de Recursos Locais**: Ler recursivamente uma pasta escolhida pelo usuário, identificando arquivos `.tf`, blocos `resource` e blocos `data`.
- **F02 - Parseamento de Módulos Locais**: Identificar blocos `module` e mapear as fontes locais para validar de onde aquele módulo está sendo importado.
- **F03 - Identificação de Dependências (Edges)**: Buscar relações estáticas explícitas (`depends_on`) e implícitas (uso de referências via interpolação, exemplo: `aws_vpc.main.id`) para criar as arestas do grafo.
- **F04 - Parseamento de Environments/Inventories**: Ler e mapear de forma correta arquivos de definições de variáveis (ex: `.tfvars` separados por diretórios de ambientes como Prod, Stg, Dev).

## 2. Módulo de Interface Visual (Frontend)

Com a arquitetura de dados pronta, os mesmos deverão ser plotados visualmente.
- **F05 - Upload/Seleção de Diretório Local**: Uma tela/botão onde o usuário, via navegador, consiga apontar um diretório da sua máquina e o browser envie ao backend temporariamente os arquivos de texto para parse.
- **F06 - Visualização em Grafo Interativo**: Exibir a relação de todos os componentes plotados pelo parser utilizando uma biblioteca de nós e arestas. Deve suportar zoom in, zoom out e *fit to screen*.
- **F07 - Painel de Detalhes do Recurso**: Ao interagir (clicar) num nó do grafo, abrir um painel lateral dinâmico mostrando seu nome, tipo, e cruzamento de suas variáveis pelas *envs* mapeadas (mostrar que em Prod ele é tamanho *large* e em Dev é *small*).
- **F08 - Árvore de Arquivos (File Explorer)**: Painel lateral reproduzindo a estrutura de diretórios do projeto Terraform submetido.

## 3. Módulo de Integração com Git/Remoto (Should Have / Could Have)

Uma vez validadas as pastas locais, a experiência do SRE pode ser melhorada com integrações diretas.
- **F09 - Provider Git (Clonagem Temporária)**: O usuário irá inserir uma URL de um repositório git e um *Personal Access Token* (se privado). O backend clonará o repositório em `/tmp`, fará o parser e em seguida apagará a pasta local por segurança.
- **F10 - Suporte a repósitórios multi-diretórios (Mono-repos)**: Se a URL do Git for um monorepositorio, o backend deve ser capaz de apresentar a primeira camada para que o usuário direcione qual subpasta contém os *tf* do projeto em questão.
