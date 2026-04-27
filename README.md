
<div align="left">
  <img src="https://img.shields.io/static/v1?label=license&message=MIT&color=5965E0&labelColor=121214" alt="License">
  <br>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Gin-0077B5?style=for-the-badge&logo=gin&logoColor=white" alt="Gin">
  <img src="https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=ethereum&logoColor=white" alt="Go-Ethereum">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/Chainlink-375BD2?style=for-the-badge&logo=chainlink&logoColor=white" alt="Chainlink">
</div>

<br>

<h1>Chainlink Price Feed API (Golang)</h1>

> 📖 **Artigo de referência:** Se você se interessa pelo desenvolvimento de integrações Web3 com Go, confira este tutorial:
> [Consultando Preços de Criptomoedas com Chainlink e Golang](https://dev.to/dev-araujo/como-integrar-chainlink-data-feeds-em-go-para-multiplos-tokens-ekb)

<br/>

Esta é uma API RESTful de alta performance desenvolvida em **Go** que atua como um gateway seguro e eficiente para os **[Chainlink Data Feeds](https://docs.chain.link/data-feeds/price-feeds/addresses)**. A aplicação abstrai a complexidade da comunicação direta com a blockchain Ethereum, permitindo que clientes e outros microsserviços obtenham cotações de criptoativos em tempo real de forma simples e escalável.

A arquitetura do backend foi projetada com separação de responsabilidades, isolando a camada de roteamento HTTP, a lógica de negócios e as chamadas RPC via contratos inteligentes.

<div align="center">
  <img src='./assets/gopher-link.png' width='300'>
</div>

## 🛠️ Stack Tecnológico do Backend

* **[Go (1.24.4+)](https://golang.org/)**: Linguagem principal, garantindo alta concorrência e tipagem forte.
* **[Gin Web Framework](https://github.com/gin-gonic/gin)**: Roteador HTTP focado em máxima performance e gerenciamento de middlewares.
* **[Go-Ethereum (Geth)](https://github.com/ethereum/go-ethereum)**: Cliente oficial utilizado para estabelecer a conexão RPC com nós da rede Ethereum e executar chamadas de leitura nos Smart Contracts.
* **[Golang Sync](https://pkg.go.dev/golang.org/x/sync)**: Utilizado para gerenciar concorrência avançada e otimizar buscas assíncronas.

> *Nota: O projeto também inclui um serviço web leve (desenvolvido com Go Templates e HTMX) para visualização rápida dos dados, servido de forma independente.*
> 
> Acesse a aplicação rodando em produção:
>
>  **Web Interface:** [https://crypto.dev-araujo.com.br/](https://crypto.dev-araujo.com.br/)


<img src="./assets/interface.png" alt="Interface web" width="100%" style="border-radius: 10px; box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);"/>


---

## 🏗️ Estrutura do Projeto

O código-fonte segue diretrizes sólidas de design, facilitando a manutenção e a escalabilidade:

```text
.
├── cmd/
│   ├── api/          # Ponto de entrada (main) do servidor REST da API
│   └── web/          # Ponto de entrada do serviço frontend (HTMX)
├── contracts/        # Bindings gerados via abigen para os contratos da Chainlink
├── internal/
│   ├── config/       # Carregamento de variáveis de ambiente e configurações (RPC, Portas)
│   ├── handler/      # Controladores (Handlers) HTTP do Gin
│   └── service/      # Regras de negócio e comunicação com os Smart Contracts (go-ethereum)
```

---

## 🚀 Executando a Aplicação

### Pré-requisitos

* Go 1.24.4 ou superior
* Docker e Docker Compose (Opcional, para execução em containers)
* Uma URL RPC da Ethereum Mainnet (via [Infura](https://infura.io/), [Alchemy](https://www.alchemy.com/) ou [Public Node](https://ethereum.publicnode.com/))

### 1. Configuração de Ambiente

Clone o repositório e configure as variáveis de ambiente necessárias para o backend:

```sh
git clone [https://github.com/dev-araujo/chainlink-price-feed.git](https://github.com/dev-araujo/chainlink-price-feed.git)
cd chainlink-price-feed
cp .env.example .env
```

Edite o arquivo `.env`:

```env
# Configurações Essenciais do Backend
RPC_URL="[https://mainnet.infura.io/v3/SEU_ID_AQUI](https://mainnet.infura.io/v3/SEU_ID_AQUI)"
SERVER_PORT="8080"
GIN_MODE="release" # Altere para "debug" durante o desenvolvimento local

# Configurações do Serviço Web Secundário
WEB_PORT="8081"
API_URL="http://localhost:8080"
```

### 2. Inicialização (Opção via Docker)

A maneira mais prática de subir todo o ecossistema (API + Web UI) em instâncias isoladas:

```sh
docker-compose up --build
```
* A API estará escutando e aceitando requisições em `http://localhost:8080`.

### 3. Inicialização (Opção Local - Apenas Backend)

Se o objetivo for apenas testar ou desenvolver na API:

```sh
# Baixa e verifica as dependências (gin, go-ethereum, etc)
go mod tidy

# Inicia o servidor da API
go run ./cmd/api/main.go
```

---

## 📡 Endpoints e Integração

O backend suporta a consulta direta aos oráculos da Chainlink para uma lista selecionada de ativos.

**Ativos Suportados (parâmetro `:asset`):** `btc`, `eth`, `link`, `uni`, `1inch`, `paxg`, `stx`

| Método | Endpoint | Descrição |
| --- | --- | --- |
| `GET` | `/health` | Liveness probe para verificação de status do servidor. |
| `GET` | `/api/price/:asset/usd` | Consulta o contrato na blockchain e retorna o valor em Dólar (USD). |
| `GET` | `/api/price/:asset/brl` | Consulta e converte o valor do ativo para Real Brasileiro (BRL). |
| `GET` | `/api/price/all/usd` | Retorna um array com a cotação de todos os ativos suportados em USD. |
| `GET` | `/api/price/all/brl` | Retorna um array com a cotação de todos os ativos suportados em BRL. |

### Exemplo de Uso (cURL)

```bash
curl -X GET http://localhost:8080/api/price/eth/usd -H "Accept: application/json"
```

**Resposta (JSON):**

```json
{
    "pair": "ETH/USD",
    "price": "3000.00",
    "timestamp": 1678886400,
    "imageUrl": "[https://cryptologos.cc/logos/ethereum-eth-logo.png?v=040](https://cryptologos.cc/logos/ethereum-eth-logo.png?v=040)"
}
```

---

## Author 👷

<img src="https://user-images.githubusercontent.com/97068163/149033991-781bf8b6-4beb-445a-913c-f05a76a28bfc.png" width="10%" alt="caricatura do autor desse repositório"/>

**Adriano P Araujo** 
<br>
[![LinkedIn](https://img.shields.io/badge/LinkedIn-0A66C2?logo=linkedin&logoColor=white&style=for-the-badge)](https://www.linkedin.com/in/araujocode/)
<br>
[![GitHub](https://img.shields.io/badge/GitHub-181717?logo=github&logoColor=white&style=for-the-badge)](https://github.com/dev-araujo)
