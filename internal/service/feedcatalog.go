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
	feedCatalogURL     = "https://reference-data-directory.vercel.app/feeds-mainnet.json"
	catalogRefreshTTL  = 6 * time.Hour
	catalogHTTPTimeout = 15 * time.Second
)

type feedEntry struct {
	Name         string   `json:"name"`
	ProxyAddress string   `json:"proxyAddress"`
	AssetName    string   `json:"assetName"`
	FeedType     string   `json:"feedType"`
	FeedCategory string   `json:"feedCategory"`
	Decimals     int      `json:"decimals"`
	Docs         feedDocs `json:"docs"`
}

type feedDocs struct {
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
	Hidden     bool   `json:"hidden"`
}

type FeedEntry struct {
	Symbol       string
	AssetName    string
	ProxyAddress string
	Decimals     int
}

type FeedCatalog struct {
	mu          sync.RWMutex
	feeds       map[string]FeedEntry
	lastFetched time.Time
	httpClient  *http.Client
}

func NewFeedCatalog() (*FeedCatalog, error) {
	c := &FeedCatalog{
		httpClient: &http.Client{Timeout: catalogHTTPTimeout},
	}
	if err := c.refresh(); err != nil {
		return nil, fmt.Errorf("falha ao carregar catálogo de feeds: %w", err)
	}
	go c.startAutoRefresh()
	return c, nil
}

func (c *FeedCatalog) All() map[string]FeedEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]FeedEntry, len(c.feeds))
	for k, v := range c.feeds {
		out[k] = v
	}
	return out
}

func (c *FeedCatalog) Get(symbol string) (FeedEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.feeds[strings.ToLower(symbol)]
	return e, ok
}

func (c *FeedCatalog) startAutoRefresh() {
	ticker := time.NewTicker(catalogRefreshTTL)
	defer ticker.Stop()
	for range ticker.C {
		if err := c.refresh(); err != nil {
			log.Printf("[FeedCatalog] falha ao atualizar catálogo: %v", err)
		} else {
			log.Printf("[FeedCatalog] catálogo atualizado — %d feeds USD/Crypto", c.count())
		}
	}
}

func (c *FeedCatalog) count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.feeds)
}

func (c *FeedCatalog) refresh() error {
	resp, err := c.httpClient.Get(feedCatalogURL)
	if err != nil {
		return fmt.Errorf("erro HTTP ao buscar catálogo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("catálogo retornou status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler corpo do catálogo: %w", err)
	}

	var raw []feedEntry
	if err := json.Unmarshal(body, &raw); err != nil {
		return fmt.Errorf("erro ao decodificar catálogo JSON: %w", err)
	}

	feeds := make(map[string]FeedEntry)
	for _, f := range raw {
		if !isValidCryptoUSDFeed(f) {
			continue
		}
		symbol := strings.ToLower(f.Docs.BaseAsset)
		if symbol == "" {
			continue
		}
		if _, exists := feeds[symbol]; exists {
			continue
		}
		feeds[symbol] = FeedEntry{
			Symbol:       symbol,
			AssetName:    f.AssetName,
			ProxyAddress: f.ProxyAddress,
			Decimals:     f.Decimals,
		}
	}

	c.mu.Lock()
	c.feeds = feeds
	c.lastFetched = time.Now()
	c.mu.Unlock()

	log.Printf("[FeedCatalog] carregados %d feeds Crypto/USD do Chainlink", len(feeds))
	return nil
}

func isValidCryptoUSDFeed(f feedEntry) bool {
	if f.ProxyAddress == "" || f.ProxyAddress == "0x0000000000000000000000000000000000000000" {
		return false
	}
	if f.Docs.Hidden {
		return false
	}
	if !strings.EqualFold(f.Docs.QuoteAsset, "USD") {
		return false
	}
	if !strings.EqualFold(f.FeedType, "Crypto") {
		return false
	}
	if f.Docs.BaseAsset == "" {
		return false
	}
	return true
}
