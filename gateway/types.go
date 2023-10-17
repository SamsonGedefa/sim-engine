package gateway

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/SamsonGedefa/simulator/main.go/internal/order"
	"github.com/gorilla/websocket"
)

type (
	// MarketDataProvider defines an interface for retrieving market data.
	MarketDataProvider interface {
		GetMarketData() MarketData
		GetOrderBook() OrderBook
	}

	// DataStorer interface {
	// 	StoreData(data MarketData) error
	// 	RetrieveData() (MarketData, error)
	// }
	// MarketData represents live market data.
	MarketData struct {
		mu        sync.Mutex
		Symbol    string    `json:"symbol"`
		Price     float64   `json:"price"`
		Timestamp time.Time `json:"timestamp"`
	}

	OrderBook struct {
		mu     sync.Mutex
		Symbol string        `json:"symbol"`
		Bids   []order.Order `json:"bids"`
		Asks   []order.Order `json:"asks"`
	}

	// WebSocketHandler handles WebSocket connections for streaming market data.
	WebSocketHandler struct {
		upgrader           *websocket.Upgrader
		marketDataProvider MarketDataProvider
	}
)

func (ob *OrderBook) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Symbol string        `json:"symbol"`
		Bids   []order.Order `json:"bids"`
		Asks   []order.Order `json:"asks"`
	}{
		Symbol: ob.Symbol,
		Bids:   ob.Bids,
		Asks:   ob.Asks,
	})
}

func (md *MarketData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Symbol    string    `json:"symbol"`
		Price     float64   `json:"price"`
		Timestamp time.Time `json:"timestamp"`
	}{
		Symbol:    md.Symbol,
		Price:     md.Price,
		Timestamp: md.Timestamp,
	})
}
