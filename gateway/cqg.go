package gateway

import (
	"log"

	"github.com/gorilla/websocket"
)

// cgg is an example data provider.
type cqg struct {
	conn *websocket.Conn
}

type CQG interface {
	StartWebSocketConnection(url string, gateway *gateway) error
	GetRealTimeBars(symbol string) (MarketData, error)
	GetOrderBook(symbol string) (OrderBook, error)
}

func NewCQG() CQG {
	return &cqg{}
}

func (c *cqg) StartWebSocketConnection(url string, gateway *gateway) error {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	// Start listening for data from the data provider
	go c.listenForData(gateway)

	return nil
}

func (c *cqg) listenForData(gateway *gateway) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from data provider:", err)
			// Add reconnection logic here if needed.
			return
		}

		// Update the gateway's data structures
		gateway.processProviderData(message)
	}
}

func (c cqg) GetRealTimeBars(symbol string) (MarketData, error) {
	return MarketData{}, nil
}

func (c cqg) GetOrderBook(symbol string) (OrderBook, error) {
	return OrderBook{}, nil
}
