# Tarefas: F01 a F04 (Core Backend / Parser Engine)

Esta frente concentra o desafio técnico de ler os arquivos Terraform locais (HCL) e montar a Árvore Sintática e o Grafo de dependências. A ordem de execução respeita a priorização MoSCoW.

## F01 - Parseamento de Recursos Locais (Must)
- [ ] Ler arquivos `.tf` no diretório informado pelo comando CLI.
- [ ] Implementar a library `hashicorp/hcl/v2/hclsyntax` para ler o `hcl.File`.
- [ ] Iterar sobre os blocos (Blocks) do tipo `resource` e `data`.
- [ ] Extrair os Labels de Tipo e Nome (Ex: `aws_vpc` e `main`).
- [ ] Montar e salvar a `struct Node` inicial no grafo interno.

## F02 - Parseamento de Módulos Locais (Must)
- [ ] Identificar os blocos do tipo `module`.
- [ ] Ler os atributos obrigatórios (`source` e `version`).
- [ ] Caso a *source* inicie com `./` ou `../` (local), identificar de qual subdiretório ele está puxando para gerar futuros nós dependentes (F03).
- [ ] Agrupar o módulo no modelo relacional do Grafo interno (`ID = module.vpc`).

## F03 - Identificação de Dependências (Must)
*Esta é a feature mais complexa do parser. Demanda testes unitários.*
- [ ] **Dependências Explícitas:** Ler o atributo nativo `depends_on = [ ... ]` de cada node. Se existir, criar um `Edge (Aresta)` do tipo *explicit* apontando pro ID respectivo.
- [ ] **Dependências Implícitas (Traversal):** Analisar o corpo (`Body`) do HCL de um recurso usando `hcl.Traversal`.
  - [ ] Ao encontrar `vpc_id = aws_vpc.main.id`, extrair a raiz `aws_vpc.main`.
  - [ ] Criar a aresta do tipo *implicit* entre o nó atual e o nó apontado.
- [ ] Eliminar dependências cíclicas (Circular Dependency Check), logando Warning se houver.

## F04 - Parseamento de Environments (Should)
- [ ] Identificar pastas conhecidas de *environments* ou buscar arquivos terminados em `*.auto.tfvars` ou `*.tfvars`.
- [ ] Subir o parser `hcl` focado em Atributos simples para extrair Key/Value.
- [ ] Correlacionar esses valores aos `variables` lidos nos arquivos principais.
- [ ] Devolver o `map[string]map[string]string` pronto pro Frontend carregar no painel lateral.
