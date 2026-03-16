# Requisitos Não-Funcionais (RNF)

Este documento detalha as restrições e premissas técnicas de infraestrutura, performance e ambiente necessárias para o **Terraform Tracer**, conforme conversas com a arquitetura.

## RNF01 - Modelo de Deploy e Distribuição
- **Natureza:** Ferramenta Local / CLI / Desktop.
- **Descrição:** A aplicação será desenhada primariamente para rodar localmente no computador do próprio engenheiro (DevOps/SRE) como um servidor temporário. Ao iniciar o processo (ex: `tracer start`), o backend em Python sobe uma porta local (ex: `localhost:8080`) e expõe o frontend da SPA via browser. Sem necessidade imediata de *containers orchestration* (k8s) ou *cloud deployment*.

## RNF02 - Performance e Tempo de Resposta
- **Natureza:** Processamento Síncrono (Padrão).
- **Descrição:** O tempo de resposta para a leitura (*parse*) do diretório será síncrono. Assumimos que, em média, as pastas de projeto dos usuários retornarão os grafos em poucos segundos. O frontend lidará com a espera no cliente (loader screen), não requerendo no MVP arquiteturas complexas como filas assíncronas (Celery, Redis queues) ou WebSockets. 

## RNF03 - Gestão de Estado e Caching
- **Natureza:** Cache Temporário de Sessão (Em Memória Local).
- **Descrição:** Embora a aplicação seja *Stateless* sem banco de dados (sem Postgres ou MongoDB persistentes no disco do usuário), implementaremos um sistema de *cache temporário*. Enquanto o servidor Python estiver rodando na máquina ou o frontend estiver aberto, o último projeto "parseado" pode ficar mantido na memória (Dicionário Python ou LocalStorage do JS). Se a ferramenta reiniciar, o estado é completamente perdido e refeito.

## RNF04 - Autenticação e Segurança de Rede
- **Natureza:** Sistema Aberto (Sem Auth).
- **Descrição:** Por rodar puramente na interface de *loopback* local (`localhost` ou `127.0.0.1`), a aplicação não terá camada de autenticação, login com perfis nem criptografia de tráfego de borda (HTTPS). Se futuramente o cliente a hospedar em ambiente compartilhado da empresa, ele ficará responsável por envelopar a ferramenta atrás de um proxy corporativo.
