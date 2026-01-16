package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dev-araujo/chainlink-price-feed/internal/config"
	"github.com/dev-araujo/chainlink-price-feed/internal/handler"
	"github.com/dev-araujo/chainlink-price-feed/internal/service"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Web structs
type PriceViewModel struct {
	Pair           string
	Price          string
	ImageURL       string
	LastUpdate     string
	CurrencySymbol string
}

var currencySymbols = map[string]string{
	"brl": "R$",
	"usd": "$",
}

var templates = template.Must(template.ParseFiles("./web/templates/prices.html"))

func main() {
	cfg := config.Load()
	if cfg.RpcURL == "" {
		log.Println("AVISO: RPC_URL não definida, usando padrão ou erro pode ocorrer.")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	client, err := ethclient.Dial(cfg.RpcURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao nó da rede principal da Ethereum: %v", err)
	}
	log.Println("Conectado com sucesso à rede principal da Ethereum!")

	exchangeService := service.NewExchangeService()
	chainlinkService := service.NewChainlinkService(client, exchangeService)
	assetService := service.NewAssetService()

	priceHandler := handler.NewPriceHandler(chainlinkService, assetService)

	router := gin.Default()
	router.Use(cors.Default())

	router.Static("/styles", "./web/styles")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	router.GET("/ready", func(c *gin.Context) {

		html := `
		<div id="price-list-container" class="price-list" 
			 hx-get="/prices/all/brl" 
			 hx-trigger="load, change from:#currency-select" 
			 hx-swap="innerHTML" 
			 hx-indicator="#loading-indicator">
		</div>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})

	router.GET("/prices/all/:currency", func(c *gin.Context) {
		currency := strings.ToLower(c.Param("currency"))

		var priceData []*service.PriceData
		var err error

		if currency == "brl" {
			priceData, err = chainlinkService.GetAllPricesBRL()
		} else {
			priceData, err = chainlinkService.GetAllPricesUSD()
		}

		if err != nil {
			log.Printf("Erro ao obter preços: %v", err)
			c.String(http.StatusInternalServerError, "<div class='error'>Erro ao processar dados de preços</div>")
			return
		}

		currencySymbol, ok := currencySymbols[currency]
		if !ok {
			currencySymbol = "$"
		}

		viewModels := make([]PriceViewModel, len(priceData))
		for i, p := range priceData {
			assetSymbol := strings.ToLower(strings.Split(p.Pair, "/")[0])
			imageURL, _ := assetService.GetAssetImageURL(assetSymbol)

			viewModels[i] = PriceViewModel{
				Pair:           p.Pair,
				Price:          p.Price.Text('f', 2),
				ImageURL:       imageURL,
				LastUpdate:     time.Unix(p.Timestamp, 0).Format("15:04:05"),
				CurrencySymbol: currencySymbol,
			}
		}

		err = templates.ExecuteTemplate(c.Writer, "prices.html", viewModels)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error executing template: %v", err)
		}
	})

	priceHandler.RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ATIVO"})
	})

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Iniciando applicação unificada (Web + API) na porta %s", cfg.ServerPort)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
