package gateway

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxRetries = 5
	retryDelay = 5 * time.Second
)

// cgg is an example data provider.
type cqg struct {
	conn *websocket.Conn
}

type CQG interface {
	StartListening(conn *websocket.Conn, gateway *gateway) error
	GetRealTimeBars(symbol string) (CandleStick, error)
	GetOrderBook(symbol string) (OrderBook, error)
}

func NewCQG() CQG {
	return &cqg{}
}

func (c *cqg) StartListening(conn *websocket.Conn, gateway *gateway) error {
	c.conn = conn

	// Start listening for data from the data provider
	go c.listenForData(gateway)

	return nil
}

func (c *cqg) listenForData(gateway *gateway) error {
	retries := 0

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from data provider:", err)

			if retries >= maxRetries {
				return fmt.Errorf("exceeded max retries for reading from data provider: %v", err)
			}

			// TODO: add exponential backoff

			time.Sleep(retryDelay)
			retries++

			continue
		}

		// Reset the retry counter after a successful read
		retries = 0

		// Update the gateway's data structures
		gateway.processProviderData(message)
	}
}

func (c cqg) GetRealTimeBars(symbol string) (CandleStick, error) {
	return CandleStick{}, nil
}

func (c cqg) GetOrderBook(symbol string) (OrderBook, error) {
	return OrderBook{}, nil
}
