# 🚀 Desafio Técnico: Performance e Análise de Dados via API

Este repositório contém a implementação de uma API utilizando **Go (Golang)** para ingestão e análise de dados de usuários em massa. Essa é a resolução de um desafio proposto pela Codecon que pode ser achado [aqui](https://github.com/codecon-dev/desafio-1-1s-vs-3j)
.

---

## 🧠 Desafio

- Receber um JSON com **100.000 usuários**.
- Manter os dados na **memória** (sem persistência em banco).
- Expor **endpoints** para análise.

### 📥 Estrutura esperada do JSON de entrada:

```json
{
  "id": "uuid",
  "name": "string",
  "age": "int",
  "score": "int",
  "active": "bool",
  "country": "string",
  "team": {
    "name": "string",
    "leader": "bool",
    "projects": [{ "name": "string", "completed": "bool" }]
  },
  "logs": [{ "date": "YYYY-MM-DD", "action": "login/logout" }]
}
```

---

## 🔧 Endpoints da API

### `POST /users`
- Recebe um JSON com os usuários.
- Armazena os dados em memória.

### `GET /superusers`
- Retorna os usuários com `score >= 900` e `active = true`.
- Inclui tempo de processamento da requisição.

### `GET /top-countries`
- Agrupa os **superusuários** por país.
- Retorna os **5 países com mais superusuários**.

### `GET /team-insights`
- Agrupa os usuários por `team.name`.
- Para cada time, retorna:
  - Total de membros.
  - Quantidade de líderes.
  - Total de projetos concluídos.
  - Porcentagem de membros ativos.

### `GET /active-users-per-day`
- Conta o número de **logins por dia** com base nos logs.

### `GET /evaluation`
- Executa uma **autoavaliação** dos endpoints principais.
- Para cada endpoint avaliado, retorna:
  - Status da resposta (esperado: 200).
  - Tempo de resposta em milissegundos.
  - Se o retorno é um JSON válido.

---

## 🚀 Como Executar

1. Clone o repositório:
   ```bash
   git clone https://github.com/feribeirods/desafio-json-go.git
   cd desafio-json-go
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Rode a API:
   ```bash
   go run main.go
   ```

4. Acesse via: `http://localhost:8081`

---

## ✅ Avaliação Automática

Acesse o endpoint `/evaluation` após importar os dados para validar automaticamente:

```bash
curl http://localhost:8081/evaluation
```

O JSON retornado mostrará a performance e validade dos principais endpoints.

---

## 🛠 Tecnologias

- [Go (Golang)](https://golang.org/)
- Armazenamento em memória com estrutura otimizada

---

## 📌 Considerações

- Foco em performance e estrutura clara dos dados.
- Nenhum banco de dados externo é necessário.
- A API está pronta para testes em escala e validações automatizadas.

---

## 🧑‍💻 Autor

Desenvolvido por [Luiz Fernando](https://github.com/feribeirods) como solução para o desafio técnico de performance e análise de dados via API.
