# Diagrama de Sequência - Extração de Dependências

Este diagrama detalha a parte mais crítica e complexa: como o backend de Go interpreta o HCL para achar de quem um recurso depende (*Edges* implícitas e explícitas).

```mermaid
sequenceDiagram
    participant CLI as Cmd/Tracer
    participant Parser as Internal/Parser (GO)
    participant HCL as hashicorp/hcl/v2
    participant Memory as Graph Store (RAM)

    CLI->>Parser: Iniciar parse no path "./infra"
    Parser->>HCL: Carregar files `*.tf` no Parser
    HCL-->>Parser: Retorna HCL File/Body blocks
    
    loop Para cada Bloco Resource
        Parser->>HCL: Extrair bloco (Type, Labels)
        Parser->>Memory: Criar Node(Resource)
        
        Parser->>HCL: Extrair bloco `depends_on` (Explícitas)
        alt Encontrou `depends_on`
            Parser->>Memory: Criar Aresta (Edge type: Explicit)
        end
        
        Parser->>HCL: Extrair Expressões de Atributos (ex: `aws_vpc.main.id`)
        alt Encontrou Interpolações (Implícitas)
            Parser->>Parser: Resolver HCL Traversal Variables
            Parser->>Memory: Criar Aresta (Edge type: Implicit, target: aws_vpc.main)
        end
    end
    
    Parser-->>CLI: Parse completo (Nodes + Edges)
```
