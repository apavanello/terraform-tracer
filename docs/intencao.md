# Intenção do Projeto: Terraform Tracer

## Resumo da Ideia
Uma interface gráfica e interativa (website) orientada a fluxos visuais que faz o rastreio de recursos do Terraform. O objetivo é permitir que o usuário aponte para um repositório Git ou uma pasta local contendo código Terraform, selecione um recurso específico e visualize de forma gráfica e clara suas dependências, módulos associados e as diferentes *environments* (ambientes/inventories) onde ele está configurado.

## Maturidade e Escopo
**MVP (Produto Mínimo Viável):** A meta para a primeira versão é construir uma ferramenta funcional e utilizável no dia a dia, com interface amigável, provendo valor imediato para o time ao facilitar a consulta e mapeamento de dependências.

## Problema a ser Resolvido e Benefícios
Em operações de Infraestrutura como Código (IaC), entender o raio de impacto de uma alteração pode ser complexo apenas lendo código. A falta de visibilidade rápida sobre como os recursos do Terraform se conectam e afetam os ambientes pode levar a incidentes. A ferramenta trará os seguintes benefícios:
- Acelerar o processo de **troubleshooting**.
- Facilitar a **análise de impacto** antes de aplicar mudanças.
- Evitar quebras acidentais em ambientes produtivos compartilhados.

## Público-Alvo
**Engenheiros DevOps e SREs** que administram e evoluem constantemente a infraestrutura, precisando de confiabilidade e agilidade para analisar o estado e a topologia dos recursos nas variadas *envs*.

## Mecânica de Extração de Dados
A aplicação efetuará a análise estática (parse) do código Terraform a partir de duas origens possíveis:
1. **Conexão com repositório Git:** A interface recebe um link, faz o download ou clone temporário do repositório, efetua a leitura e exibe os dados.
2. **Pasta Local:** O usuário pode apontar ou dar o upload de um diretório já existente na máquina.

## Interface e Experiência do Usuário (UX)
A visualização deve ser guiada por **Fluxos Gráficos (nós e arestas no estilo de mapa mental ou diagramas de componentes)**. Fugindo de visualizações puramente em texto ou tabelas maçantes, o usuário enxergará de forma interativa como blocos, dependências (*depends_on*, referências lógicas) e os módulos interagem entre as separações lógicas de ambientes (*inventories*).
