package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	coinGeckoURL        = "https://api.coingecko.com/api/v3/coins/list?include_platform=false"
	coinGeckoDetailURL  = "https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=false&community_data=false&developer_data=false"
	assetCacheRefreshTTL = 24 * time.Hour
	assetHTTPTimeout    = 10 * time.Second
)

type coinListEntry struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type coinDetail struct {
	Image struct {
		Small string `json:"small"`
		Large string `json:"large"`
	} `json:"image"`
}

type AssetService struct {
	mu         sync.RWMutex
	imageCache map[string]string
	symbolToID map[string]string
	httpClient *http.Client
}

func NewAssetService() *AssetService {
	s := &AssetService{
		imageCache: make(map[string]string),
		symbolToID: make(map[string]string),
		httpClient: &http.Client{Timeout: assetHTTPTimeout},
	}
	go s.loadCoinList()
	return s
}

func (s *AssetService) GetAssetImageURL(symbol string) (string, error) {
	key := strings.ToLower(symbol)

	s.mu.RLock()
	if url, ok := s.imageCache[key]; ok {
		s.mu.RUnlock()
		return url, nil
	}
	s.mu.RUnlock()

	coinID, err := s.resolveCoinID(key)
	if err != nil {
		return fallbackLogoURL(key), nil
	}

	imageURL, err := s.fetchImageFromCoinGecko(coinID)
	if err != nil {
		return fallbackLogoURL(key), nil
	}

	s.mu.Lock()
	s.imageCache[key] = imageURL
	s.mu.Unlock()

	return imageURL, nil
}

func (s *AssetService) resolveCoinID(symbol string) (string, error) {
	s.mu.RLock()
	id, ok := s.symbolToID[symbol]
	s.mu.RUnlock()
	if ok {
		return id, nil
	}
	return "", fmt.Errorf("coin ID não encontrado para símbolo '%s'", symbol)
}

func (s *AssetService) fetchImageFromCoinGecko(coinID string) (string, error) {
	url := fmt.Sprintf(coinGeckoDetailURL, coinID)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro HTTP ao buscar imagem para %s: %w", coinID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("rate limit CoinGecko")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CoinGecko retornou status %d para %s", resp.StatusCode, coinID)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var detail coinDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return "", err
	}

	if detail.Image.Small == "" {
		return "", fmt.Errorf("imagem não encontrada para %s", coinID)
	}
	return detail.Image.Small, nil
}

func (s *AssetService) loadCoinList() {
	resp, err := s.httpClient.Get(coinGeckoURL)
	if err != nil {
		log.Printf("[AssetService] falha ao buscar lista de coins: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[AssetService] falha ao ler lista de coins: %v", err)
		return
	}

	var list []coinListEntry
	if err := json.Unmarshal(body, &list); err != nil {
		log.Printf("[AssetService] falha ao decodificar lista de coins: %v", err)
		return
	}

	symbolToID := make(map[string]string, len(list))
	for _, c := range list {
		sym := strings.ToLower(c.Symbol)
		if _, exists := symbolToID[sym]; !exists {
			symbolToID[sym] = c.ID
		}
	}

	s.mu.Lock()
	s.symbolToID = symbolToID
	s.mu.Unlock()

	log.Printf("[AssetService] mapeados %d símbolos de coins do CoinGecko", len(symbolToID))

	time.AfterFunc(assetCacheRefreshTTL, s.loadCoinList)
}

func fallbackLogoURL(symbol string) string {
	knownFallbacks := map[string]string{
		"btc":   "https://assets.coingecko.com/coins/images/1/small/bitcoin.png",
		"eth":   "https://assets.coingecko.com/coins/images/279/small/ethereum.png",
		"link":  "https://assets.coingecko.com/coins/images/877/small/chainlink-new-logo.png",
		"uni":   "https://assets.coingecko.com/coins/images/12504/small/uni.jpg",
		"paxg":  "https://assets.coingecko.com/coins/images/9519/small/paxg.PNG",
		"stx":   "https://assets.coingecko.com/coins/images/2069/small/Stacks_logo_full.png",
		"1inch": "https://assets.coingecko.com/coins/images/13469/small/1inch-token.png",
	}
	if url, ok := knownFallbacks[symbol]; ok {
		return url
	}
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=375BD2&color=fff&size=64&bold=true", strings.ToUpper(symbol))
}
