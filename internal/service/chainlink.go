package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/dev-araujo/chainlink-price-feed/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

type PriceData struct {
	Pair      string
	Price     *big.Float
	Timestamp int64
}

type ChainlinkService struct {
	client          *ethclient.Client
	catalog         *FeedCatalog
	exchangeService *ExchangeService
}

func NewChainlinkService(client *ethclient.Client, catalog *FeedCatalog, exchangeService *ExchangeService) *ChainlinkService {
	return &ChainlinkService{
		client:          client,
		catalog:         catalog,
		exchangeService: exchangeService,
	}
}

func (s *ChainlinkService) GetPriceUSD(ctx context.Context, asset string) (*PriceData, error) {
	return s.fetchPriceFromChainlink(ctx, asset)
}

func (s *ChainlinkService) GetPriceBRL(ctx context.Context, asset string) (*PriceData, error) {
	assetPriceData, err := s.fetchPriceFromChainlink(ctx, asset)
	if err != nil {
		return nil, err
	}

	brlRate, err := s.exchangeService.GetBRLRate()
	if err != nil {
		return nil, fmt.Errorf("não foi possível obter a taxa de câmbio do BRL: %w", err)
	}

	priceInBRL := new(big.Float).Mul(assetPriceData.Price, brlRate)

	return &PriceData{
		Pair:      fmt.Sprintf("%s/BRL", strings.ToUpper(asset)),
		Price:     priceInBRL,
		Timestamp: assetPriceData.Timestamp,
	}, nil
}

func (s *ChainlinkService) fetchPriceFromChainlink(ctx context.Context, asset string) (*PriceData, error) {
	entry, ok := s.catalog.Get(asset)
	if !ok {
		return nil, fmt.Errorf("ativo '%s' não encontrado no catálogo de feeds", asset)
	}

	address := common.HexToAddress(entry.ProxyAddress)
	priceFeed, err := contracts.NewAggregatorV3Interface(address, s.client)
	if err != nil {
		return nil, fmt.Errorf("falha ao instanciar contrato para %s: %w", asset, err)
	}

	callOpts := &bind.CallOpts{Context: ctx}

	decimals, err := priceFeed.Decimals(callOpts)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar decimais para %s: %w", asset, err)
	}

	latestRoundData, err := priceFeed.LatestRoundData(callOpts)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar dados para %s: %w", asset, err)
	}

	price := new(big.Float).SetInt(latestRoundData.Answer)
	price.Quo(price, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))

	return &PriceData{
		Pair:      fmt.Sprintf("%s/USD", strings.ToUpper(asset)),
		Price:     price,
		Timestamp: latestRoundData.UpdatedAt.Int64(),
	}, nil
}

func (s *ChainlinkService) fetchAllPrices(priceFetcher func(ctx context.Context, asset string) (*PriceData, error)) ([]*PriceData, error) {
	feeds := s.catalog.All()
	prices := make([]*PriceData, 0, len(feeds))
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(context.Background())

	for symbol := range feeds {
		symbol := symbol
		g.Go(func() error {
			priceData, err := priceFetcher(ctx, symbol)
			if err != nil {
				return fmt.Errorf("falha ao buscar preço para %s: %w", symbol, err)
			}
			mu.Lock()
			prices = append(prices, priceData)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (s *ChainlinkService) GetAllPricesUSD() ([]*PriceData, error) {
	return s.fetchAllPrices(s.GetPriceUSD)
}

func (s *ChainlinkService) GetAllPricesBRL() ([]*PriceData, error) {
	return s.fetchAllPrices(s.GetPriceBRL)
}
