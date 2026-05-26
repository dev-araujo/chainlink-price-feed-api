package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dev-araujo/chainlink-price-feed/internal/config"
	"github.com/dev-araujo/chainlink-price-feed/internal/handler"
	"github.com/dev-araujo/chainlink-price-feed/internal/service"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.RpcURL == "" {
		log.Fatal("RPC_URL não pode ser vazia.")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	client, err := ethclient.Dial(cfg.RpcURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao nó Ethereum: %v", err)
	}
	log.Println("Conectado com sucesso à rede principal da Ethereum!")

	catalog, err := service.NewFeedCatalog()
	if err != nil {
		log.Fatalf("Falha ao inicializar catálogo de feeds: %v", err)
	}

	exchangeService := service.NewExchangeService()
	chainlinkService := service.NewChainlinkService(client, catalog, exchangeService)
	assetService := service.NewAssetService()

	priceHandler := handler.NewPriceHandler(chainlinkService, assetService)

	router := gin.Default()
	router.Use(cors.Default())

	priceHandler.RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ATIVO"})
	})

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Iniciando servidor na porta %s", cfg.ServerPort)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
