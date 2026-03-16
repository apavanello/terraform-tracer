# Fluxograma Principal do Sistema

Como o Terraform Tracer é executado e processado pelo usuário do início ao fim usando a estratégia Go + Vue embedado.

```mermaid
flowchart TD
    Start([Usuário executa: tracer start .]) --> Init(Golang inicializa CLI)
    
    Init --> CheckDir{Diretório válido?}
    CheckDir -- Não --> Error[Retorna erro no terminal]
    
    CheckDir -- Sim --> Parse[Go/hcl v2: Lê *.tf recursivamente]
    Parse --> AST[Extrai Resources e Modules]
    AST --> Vars[Carrega *.tfvars por pasta (Environments)]
    Vars --> DepMap[Cruza Interpolações para gerar Edges/Arestas]
    DepMap --> JSONCache[(Guarda grafo em Memória RAM)]
    
    JSONCache --> Server[Sobe Web Server :8080]
    Server --> OpenBrowser([Abre navegador padrão no SO])
    
    OpenBrowser --> UILoad[Vue SPA carrega estáticos embutidos]
    UILoad --> APIReq[Fetch GET /api/v1/graph]
    APIReq --> Draw[Rendeniza vis-network com os Nós]
    
    Draw --> WaitEvent((Aguardando Interação))
    WaitEvent -- "Clica no Nó" --> Sidebar[Abre Sidebar e cruza valor da Variável com o Envs]
```
