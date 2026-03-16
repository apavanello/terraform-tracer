# Tarefas: F05 a F08 (Frontend / UI / UX)

Este plano concentra o desafio de renderização do grafo no navegador. Focaremos no Vue.js para componentização e manuseio do estado local retornado pelo Backend.

## F05 - Upload/Seleção de Diretório Local (Must)
- [ ] Criar modal ou seção de Configuração Inicial no Vue (`components/Setup.vue`).
- [ ] O componente terá um campo de *input text* indicando o caminho absoluto do diretório ou consumirá a rota natural da aplicação (iniciada pelo backend CLI no diretório atual).
- [ ] Fazer chamada GET para a rota `/api/v1/graph` do backend usando `fetch` ou `axios`.
- [ ] Implementar Loader visual no botão até receber a *Promise* de retorno e atualizar a *Store* Global.

## F06 - Visualização em Grafo Interativo (Must)
- [ ] Instalar a lib `vis-network` (ou D3/SigmaJS).
- [ ] Criar o componente Wrapper `GraphView.vue` que ocupa a `main-content` da tela.
- [ ] Manipular o JSON do backend separando as matrizes e serializando no formato do Vis.Js: `nodes: [{ id, label, group }]` e `edges: [{ from, to, dashes }]`.
- [ ] Aplicar customização visual (Design System roxo/moderno) nos nós e nas conexões.
- [ ] Adicionar navegação em rede (Zoom com o scroll, Pan/Arrastar canvas).

## F07 - Painel de Detalhes do Recurso (Must)
- [ ] Adicionar EventListener ao grafo no evento `'selectNode'`.
- [ ] Criar componente `PropsDrawer.vue` para a barra lateral direita que abre (slide-in) ao selecionar um nó.
- [ ] Recuperar da State Management o dicionário completo do Node clicado (Extra properties, Type, LineNumber).
- [ ] Montar tabela iterando no array de *Environments* (Mapeando valor de Stg vs Prod lado a lado).
- [ ] Botão de "fechar" (DeselectNode).

## F08 - Árvore de Arquivos (Should)
- [ ] Criar barra lateral esquerda estática (`components/FileTree.vue`).
- [ ] Agrupar o JSON retornado agrupando os nós lidos pelo seu arquivo de origem (`node.File`).
- [ ] Renderizar hierarquia visual (ícones de pasta e de arquivo HCL).
- [ ] Fazer navegação cruzada: Clicar em um arquivo na Tree View da esquerda dispara evento que joga o Zoom (Fit) na área do Canvas onde aqueles recursos estão mapeados.
