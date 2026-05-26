package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dev-araujo/chainlink-price-feed/internal/service"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		log.Fatal("A variável de ambiente RPC_URL é necessária.")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao nó Ethereum: %v", err)
	}
	defer client.Close()

	exchangeService := service.NewExchangeService()

	catalog, err := service.NewFeedCatalog()
	if err != nil {
		log.Fatalf("Falha ao inicializar catálogo de feeds: %v", err)
	}

	chainlinkService := service.NewChainlinkService(client, catalog, exchangeService)

	asset := "xau" // eth| link| btc | aud | eur | jpy | ftse| xau |

	fmt.Printf("Buscando preço para %s/USD...\n", asset)
	priceDataUSD, err := chainlinkService.GetPriceUSD(context.Background(), asset)
	if err != nil {
		log.Fatalf("Erro ao buscar preço em USD: %v", err)
	}
	fmt.Printf("Par: %s\n", priceDataUSD.Pair)
	fmt.Printf("Preço: %s\n", priceDataUSD.Price.String())
	fmt.Printf("Última atualização (Timestamp): %d\n", priceDataUSD.Timestamp)
	fmt.Println("---------------------------------")

	fmt.Printf("Buscando preço para %s/BRL...\n", asset)
	priceDataBRL, err := chainlinkService.GetPriceBRL(context.Background(), asset)
	if err != nil {
		log.Fatalf("Erro ao buscar preço em BRL: %v", err)
	}
	fmt.Printf("Par: %s\n", priceDataBRL.Pair)
	fmt.Printf("Preço: %s\n", priceDataBRL.Price.String())
	fmt.Printf("Última atualização (Timestamp): %d\n", priceDataBRL.Timestamp)
}
