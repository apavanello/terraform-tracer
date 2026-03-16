# Modelo de Dados do Parser (AST)

Aqui definimos o modelo lógico que o backend GO e a lib `hashicorp/hcl/v2` extrairão e serializarão em JSON para o Frontend.

```mermaid
erDiagram
    Project ||--o{ File : "contains"
    Project ||--o{ Environment : "has"
    File ||--o{ Resource : "defines"
    File ||--o{ Module : "defines"
    File ||--o{ Variable : "defines"
    Resource ||--o{ Edge : "depends on (from/to)"
    Module ||--o{ Edge : "depends on (from/to)"

    Resource {
        string ID "ex: aws_vpc.main"
        string Provider "ex: aws"
        string Type "ex: aws_vpc"
        string Name "ex: main"
        json Properties "Key-Value pairs (extraídas)"
        int LineStart "Linha de inicio"
    }

    Module {
        string ID "ex: module.vpc"
        string Source "Caminho ou Registry"
        string Version "ex: ~> 3.0"
        json Inputs "Variáveis passadas ao módulo"
    }

    Variable {
        string Name "Nome da variável"
        string Type "Tipo esperado (string, map, etc)"
        string Default "Valor padrão"
    }

    Environment {
        string Name "ex: prod, stg"
        string FilePath "ex: envs/prod.tfvars"
        json FlatValues "Valores resolvidos e mesclados"
    }

    Edge {
        string FromNode "ID de origem (quem depende)"
        string ToNode "ID de destino (de quem se depende)"
        string EdgeType "Implicit (ref) ou Explicit (depends_on)"
        string Label "ex: vpc_id"
    }
```

O Frontend receberá um JSON plano otimizado contendo as listas `nodes` (unindo *Resources* e *Modules*) e `edges`, combinando com os valores por *Environment* em formato de dicionário de lookup.
