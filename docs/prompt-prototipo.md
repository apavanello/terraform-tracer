# Prompt para Geração de Protótipo

## Descrição do Projeto
Uma interface gráfica interativa (dashboard web) orientada a fluxos visuais que mapeia e rastreia recursos do Terraform. O usuário fornece uma pasta local contendo os arquivos *.tf, e o sistema exibe visualmente (nós e arestas) as dependências daquele recurso em diversas *environments* (inventories), facilitando análises de impacto e troubleshooting por DevOps e SREs.

## Tipo de Aplicação
Dashboard interativo / Ferramenta de visualização de fluxo (estilo mapa mental ou diagrama de nós).

## Especificações Visuais
- **Cor primária**: #a30eff (Roxo vibrante)
- **Estilo**: Moderno, limpo, responsivo, com uso de *glassmorphism* leve nas laterais, cantos arredondados, fontes sem serifa (Inter ou Roboto).
- **Fundo**: Claro (branco/cinza muito claro, ex: #F8F9FA).

## Funcionalidades Principais
1. **Seleção de Origem**: Um botão simulando o carregamento de uma pasta local com arquivos Terraform.
2. **Barra Lateral de Recursos**: Uma lista/tree view onde o usuário pode selecionar módulos e arquivos `.tf` encontrados.
3. **Mapeamento Visual (Graph View)**: Área central ampla que exibe um grafo interativo (nós conectáveis) representando as dependências de infraestrutura de quem depende de quem.
4. **Painel de Detalhes Relacionados**: Ao clicar em um nó na área central, abre-se um painel direito (drawer ou sidebar direita) exibindo as *environments* em que o recurso está presente, além de um snippet do código (mockado).
5. **Filtros rápidos**: Top bar com campo de busca e filtro por tipo de recurso (AWS, Azure, Módulos locais).

## Estrutura de Telas
- **Tela Única (Single Page Application)**:
  - **Header**: Logo "Terraform Tracer", campo global de busca, botão "Carregar Pasta".
  - **Left Sidebar**: Árvore de arquivos e módulos (Ex: `main.tf`, `variables.tf`, `modules/vpc`).
  - **Main Content**: Canvas/Painel infinito com o grafo das dependências. Deve passar a sensação de uma tela de "canvas" viva.
  - **Right Sidebar (Dinâmica)**: Exibida apenas quando um recurso/nó é selecionado. Mostra os metadados do recurso (nome, type, variáveis nas envs de prod, stg e dev).

## Prompt Completo para a Ferramenta
Crie um dashboard web responsivo chamado "Terraform Tracer" utilizando HTML, TailwindCSS (ou CSS puro) e JavaScript. Utilize a cor roxa #a30eff como cor primária (botões, detalhes, bordas ativas) e um fundo branco/cinza claro estilo "clean". 

A tela deve ser dividida em 3 grandes seções:
1. Uma barra lateral esquerda de 250px simulando uma árvore de arquivos Terraform (`vpc.tf`, `rds.tf`, `eks.tf` etc).
2. Uma área central expansiva que simule um grafo de rede de dependências, contendo caixas (nós) interligadas por setas (por exemplo: um nó "aws_vpc" ligado a "aws_subnet"). Use uma biblioteca de grafos mockada ou desenhe caixinhas estilizadas usando flexbox/grid para simular a topologia interativa.
3. Um painel direito (oculto por padrão e que aparece ao interagir) exibindo as propriedades do recurso selecionado: abas para "Código", "Environments (Dev/Stg/Prod)" e "Dependências".

Adicione um botão de destaque "Selecionar Pasta" no topo da sidebar esquerda. O design deve parecer uma ferramenta de desenvolvedor moderna (estilo Vercel ou Supabase, porém em tema claro), com sombras suaves, bordas arredondadas e micro-interações (hover nos nós e arquivos). Use a fonte Inter. Não adicione autenticação; exiba diretamente a tela principal com dados mocados.
