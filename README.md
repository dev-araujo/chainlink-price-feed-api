<img src="https://img.shields.io/static/v1?label=license&message=MIT&color=5965E0&labelColor=121214" alt="License">

<img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go"> <img src="https://img.shields.io/badge/Gin-0077B5?style=for-the-badge&logo=gin&logoColor=white" alt="Gin"> <img src="https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=ethereum&logoColor=white" alt="Go-Ethereum"> <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker"> <img src="https://img.shields.io/badge/Chainlink-375BD2?style=for-the-badge&logo=chainlink&logoColor=white" alt="Chainlink">

# Chainlink Price Feed com GO

Esta API, desenvolvida em **Go**, atua como uma ponte para os **[Chainlink Data Feeds](https://docs.chain.link/data-feeds/price-feeds/addresses?page=1&testnetPage=1&testnetSearch=)**, permitindo que aplicações acessem dados de preços da **blockchain Ethereum** de forma simples e eficiente.

A aplicação se conecta a um **nó da rede Ethereum**, interage com os **contratos inteligentes da Chainlink** para buscar os preços de ativos e os expõe através de uma API RESTful. Além disso, a aplicação inclui uma **interface web simples**(feita com HTMX) para visualizar esses preços.

## 🎨 Demo

Acesse a aplicação através dos links abaixo:

> Os serviços online estão em plataformas gratuitas, portanto podem estar indisponíveis.


- **API:** https://chainlink-api.onrender.com
- **Web:** https://chainlink-golang-web.onrender.com

<img src='./assets/gopher-link.png' width='300'>

## 🛠️ Stack

* [Go](https://golang.org/)
* [Gin](https://github.com/gin-gonic/gin)
* [Go-Ethereum](https://github.com/ethereum/go-ethereum)
* [Docker](https://www.docker.com/)
* [HTMX](https://htmx.org/)

## 🚀 Executando a aplicação

Siga as instruções abaixo para ter uma cópia do projeto rodando em sua máquina.

**Instalação**

1.  Clone o repositório:
    ```sh
    git clone https://github.com/dev-araujo/chainlink-price-feed.git
    ```

2.  Crie e configure o arquivo `.env`:
    ```sh
    cp .env.example .env
    ```

    Edite o arquivo `.env` com sua URL de RPC da Ethereum:

```

RPC_URL="https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID" # URL do nó RPC da Ethereum
SERVER_PORT="8080"
GIN_MODE="release"
WEB_PORT="8081"
API_URL="http://localhost:8080"

```

   > **💡 Dica:** Para um RPC gratuito, considere usar a [Public Node](https://ethereum.publicnode.com/).

---

### Opção 1: Docker (Recomendado)

**Pré-requisitos:** [Docker](https://docs.docker.com/get-docker/)

Para iniciar a aplicação, execute:
```sh
docker-compose up --build
```

A API estará disponível em `http://localhost:8080` e a aplicação web em `http://localhost:8081`.

-----

### Opção 2: Localmente

**Pré-requisitos:** [Go](https://golang.org/doc/install) (1.24.4+)

#### Rodando a API

Para iniciar a API, execute:

```sh
go run ./cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`.

#### Rodando a aplicação Web

Em um terminal separado, para iniciar a aplicação web, execute:

```sh
go run ./cmd/web/main.go
```

A aplicação web estará disponível em `http://localhost:8081`.

-----

## Endpoints da API

A API fornece os seguintes endpoints para consulta:

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `GET` | `/health` | Verifica o status da API. |
| `GET` | `/api/price/:asset/usd` | Retorna o preço do ativo especificado em USD. |
| `GET` | `/api/price/:asset/brl` | Retorna o preço do ativo especificado em BRL. |
| `GET` | `/api/price/all/usd` | Retorna o preço de todos os ativos suportados em USD. |
| `GET` | `/api/price/all/brl` | Retorna o preço de todos os ativos suportados em BRL. |

**Parâmetro de Path:**

  * `:asset`: O símbolo do ativo a ser consultado (ex: `btc`, `eth`).
      - Atualmente os seguintes ativos podem ser consultados: `1inch`, `link`, `btc`, `eth`, `paxg`, `stx`, `uni`

**Exemplo 1: Preço de um único ativo em USD**

*Requisição:*

```http
GET /api/price/eth/usd
```

*Resposta:*

```json
{
    "pair": "ETH/USD",
    "price": 3000.00,
    "timestamp": 1678886400,
    "imageUrl": "https://cryptologos.cc/logos/ethereum-eth-logo.png?v=040"
}
```

**Exemplo 2: Preço de todos os ativos em BRL**

*Requisição:*

```http
GET /api/price/all/brl
```

*Resposta:*

```json
[
    {
        "pair": "ETH/BRL",
        "price": 15000.00,
        "timestamp": 1678886400,
        "imageUrl": "https://cryptologos.cc/logos/ethereum-eth-logo.png?v=040"
    },
    {
        "pair": "BTC/BRL",
        "price": 225000.00,
        "timestamp": 1678886400,
        "imageUrl": "https://cryptologos.cc/logos/bitcoin-btc-logo.png?v=040"
    }
]
```

-----

## Interface Web

<img src="./assets/interface.png" alt="Interface web"/>

Uma interface web simples foi incluída no projeto para consumir os endpoints da API e exibir os preços de forma visualmente agradável.

A interface utiliza o **HTMX** para carregar os dados dinamicamente, permitindo que o usuário alterne entre as moedas (USD e BRL) sem a necessidade de recarregar a página.

**Características:**

  * **HTML/CSS:** Frontend leve e moderno.
  * **HTMX:** Para requisições assíncronas e atualização de conteúdo.
  * **Dinâmica:** Permite visualizar os preços de todos os ativos suportados tanto em USD quanto em BRL.




-----

#### Autor 👷

<img src="https://avatars.githubusercontent.com/u/97068163?v=4" width=120>

[Adriano P Araujo](https://www.linkedin.com/in/araujocode/)
