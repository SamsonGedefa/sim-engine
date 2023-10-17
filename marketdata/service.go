package marketdata

import (
	"sync"
	"time"

	"github.com/SamsonGedefa/simulator/main.go/internal/order"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
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
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"price"`
		Timestamp time.Time `json:"timestamp"`
	}

	OrderBook struct {
		mu    sync.Mutex
		Symbol string `json:"symbol"`
		Bids   []order.Order `json:"bids"`
		Asks   []order.Order `json:"asks"`
	}

	// WebSocketHandler handles WebSocket connections for streaming market data.
	WebSocketHandler struct {
		upgrader *websocket.Upgrader
		marketDataProvider MarketDataProvider
	}

)

func NewWebSocketHandler(upgrader *websocket.Upgrader, marketDataProvider MarketDataProvider) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader:            upgrader,
		marketDataProvider:  marketDataProvider,
	}
}

func (wh *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	conn, err := wh.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	marketDataTicker := time.NewTicker(1 * time.Second)
	defer marketDataTicker.Stop()

	for {
		select {
		case <-marketDataTicker.C:
			marketData := wh.marketDataProvider.GetMarketData()
			if err := conn.WriteJSON(marketData); err != nil {
				return err
			}
		}
	}
}