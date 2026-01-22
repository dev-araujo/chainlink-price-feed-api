
<div >
  <img src="https://img.shields.io/static/v1?label=license&message=MIT&color=5965E0&labelColor=121214" alt="License">
  <br>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Gin-0077B5?style=for-the-badge&logo=gin&logoColor=white" alt="Gin">
  <img src="https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=ethereum&logoColor=white" alt="Go-Ethereum">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/Chainlink-375BD2?style=for-the-badge&logo=chainlink&logoColor=white" alt="Chainlink">
</div>

<br>

<h1>Chainlink Price Feed com GO</h1>

> Se você se interessa pelo tema dá uma olhadinha nesse tutorial que escrevi:
> [Consultando Preços de Criptomoedas com Chainlink e Golang](https://dev.to/dev-araujo/como-integrar-chainlink-data-feeds-em-go-para-multiplos-tokens-ekb)

<br/>
<p>
  Esta API, desenvolvida em <strong>Go</strong>, atua como uma ponte para os <strong><a href="https://docs.chain.link/data-feeds/price-feeds/addresses?page=1&testnetPage=1&testnetSearch=1">Chainlink Data Feeds</a></strong>, permitindo que aplicações acessem dados de preços da <strong>blockchain Ethereum</strong> de forma simples e eficiente.
</p>

<p >
  A aplicação se conecta a um <strong>nó da rede Ethereum</strong>, interage com os <strong>contratos inteligentes da Chainlink</strong> para buscar os preços de ativos e os expõe através de uma API RESTful. Além disso, a aplicação inclui uma <strong>interface web simples</strong> (feita com HTMX) para visualizar esses preços.
</p>



<div >
  <img src='./assets/gopher-link.png' width='300'>
</div>

## 🎨 Demo

Acesse a aplicação rodando em produção:

- **Web Interface:** [https://crypto.dev-araujo.com.br/](https://crypto.dev-araujo.com.br/)


<img src="./assets/interface.png" alt="Interface web" width="100%" style="border-radius: 10px; box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);"/>


---

## 🛠️ Stack

* **[Go](https://golang.org/)**: Linguagem principal.
* **[Gin](https://github.com/gin-gonic/gin)**: Framework web de alta performance.
* **[Go-Ethereum](https://github.com/ethereum/go-ethereum)**: Cliente para interação com a blockchain.
* **[Docker](https://www.docker.com/)**: Containerização da aplicação.
* **[HTMX](https://htmx.org/)**: Interatividade no frontend sem complexidade de SPAs.

## 🚀 Executando a aplicação

Siga as instruções abaixo para ter uma cópia do projeto rodando em sua máquina.

### Pré-requisitos

* [Go](https://golang.org/doc/install) (1.24.4+)
* [Docker](https://docs.docker.com/get-docker/) (Opcional, mas recomendado)

### Instalação

1. Clone o repositório:
```sh
git clone [https://github.com/dev-araujo/chainlink-price-feed.git](https://github.com/dev-araujo/chainlink-price-feed.git)
cd chainlink-price-feed

```

2. Configure as variáveis de ambiente:
```sh
cp .env.example .env

```


3. Edite o arquivo `.env` inserindo sua URL RPC (Infura/Alchemy):
```json
RPC_URL="[https://mainnet.infura.io/v3/SEU_ID_DO_INFURA](https://mainnet.infura.io/v3/SEU_ID_DO_INFURA)"
SERVER_PORT="8080"
GIN_MODE="release"
WEB_PORT="8081"
API_URL="http://localhost:8080"

```


> **💡 Dica:** Para testes, você pode obter um RPC gratuito em [Public Node](https://ethereum.publicnode.com/).



---

### Opção 1: Docker (Recomendado)

Para iniciar todo o ambiente (API + Web) com um único comando:

```sh
docker-compose up --build

```

* A **API** estará disponível em `http://localhost:8080`
* A **Aplicação Web** estará disponível em `http://localhost:8081`

---

### Opção 2: Rodando Localmente (Sem Docker)

Se preferir rodar os serviços manualmente:

#### 1. Iniciando a API (Backend)

```sh
go run ./cmd/api/main.go

```

*A API iniciará na porta definida em `SERVER_PORT` (padrão: 8080).*

#### 2. Iniciando a Aplicação Web (Frontend)

Em um novo terminal:

```sh
go run ./cmd/web/main.go

```

*O site iniciará na porta definida em `WEB_PORT` (padrão: 8081).*

---

## 📡 Endpoints da API

A API fornece os seguintes endpoints para consulta:

| Método | Endpoint | Descrição |
| --- | --- | --- |
| `GET` | `/health` | Verifica o status da API. |
| `GET` | `/api/price/:asset/usd` | Retorna o preço do ativo em USD. |
| `GET` | `/api/price/:asset/brl` | Retorna o preço do ativo em BRL. |
| `GET` | `/api/price/all/usd` | Retorna todos os ativos em USD. |
| `GET` | `/api/price/all/brl` | Retorna todos os ativos em BRL. |

**Ativos Suportados (`:asset`):**
`btc`, `eth`, `link`, `uni`, `1inch`, `paxg`, `stx`

### Exemplos de Resposta

**GET** `/api/price/eth/usd`

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
  [![LinkedIn](https://img.shields.io/badge/LinkedIn-0A66C2?logo=linkedin&logoColor=white&style=for-the-badge)](https://www.linkedin.com/in/araujocode/) [![GitHub](https://img.shields.io/badge/GitHub-181717?logo=github&logoColor=white&style=for-the-badge)](https://github.com/dev-araujo)


