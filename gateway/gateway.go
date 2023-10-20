package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SamsonGedefa/simulator/main.go/internal/order"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/shopspring/decimal"
)

// CREATE TABLE market_data (
//     id SERIAL PRIMARY KEY,
//     symbol VARCHAR(10) NOT NULL,
//     high FLOAT8 NOT NULL,
//     low FLOAT8 NOT NULL,
//     open FLOAT8 NOT NULL,
//     close FLOAT8 NOT NULL,
//     volume INTEGER NOT NULL,
//     timestamp TIMESTAMP NOT NULL
// );

// CREATE TABLE order_books (
//     id SERIAL PRIMARY KEY,
//     symbol VARCHAR(10) NOT NULL,
//     bid_price FLOAT8 NOT NULL,
//     bid_size INTEGER NOT NULL,
//     ask_price FLOAT8 NOT NULL,
//     ask_size INTEGER NOT NULL
// );

type (

	// OrderBook represents the current order book for a symbol.
	OrderBook struct {
		Symbol string        `json:"symbol"`
		Bids   []order.Order `json:"bids"`
		Asks   []order.Order `json:"asks"`
	}

	CandleStick struct {
		High      decimal.Decimal `json:"high"`
		Low       decimal.Decimal `json:"low"`
		Open      decimal.Decimal `json:"open"`
		Close     decimal.Decimal `json:"close"`
		Volume    decimal.Decimal `json:"volume"`
		Timestamp time.Time       `json:"timestamp"`
	}
)

type Gateway interface {
	GetRealTimeBars(symbol string) (CandleStick, error)
	GetOrderBook(symbol string) (OrderBook, error)
	WebSocketHandler(w http.ResponseWriter, r *http.Request)
}

// Gateway connect us to the verious data providers we could be working with.
// it also manages the storage and notification of market data to connected clients.
type gateway struct {
	DataProvider CQG

	// redis cache
	Redis *redis.Client

	// TODO: postgres

	// Connected websocket clients
	Clients map[*websocket.Conn]bool

	// tracks shared subscriptions
	subscription *redis.PubSub

	mutex sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewGateway(c CQG, r *redis.Client, dataProviderURL string) Gateway {
	g := &gateway{
		DataProvider: c,
		Redis:        r,
		Clients:      make(map[*websocket.Conn]bool),
	}

	conn, _, err := websocket.DefaultDialer.Dial(dataProviderURL, nil)
	if err != nil {
		log.Fatal("Failed to connect to data provider:", err)
	}

	err = g.DataProvider.StartListening(conn, g)
	if err != nil {
		log.Fatal("Could not start listening to data provider:", err)
	}

	g.initializeRedisSubscription()
	return g
}

func (g *gateway) processProviderData(data []byte) {
	// TODO: Decode and process the incoming data.
	// TODO: Update the internal OrderBook, MarketData, etc.

	// Notify app clients, either via Redis
	_, err := g.Redis.Publish(context.Background(), "notification_channel", data).Result()
	if err != nil {
		log.Println("Failed to publish to Redis:", err)
	}
}

func (g *gateway) cacheInRedis(key string, data interface{}) error {
	_, err := g.Redis.Set(context.Background(), key, data, 0).Result()
	return err
}

func (g *gateway) DeregisterClient(client *websocket.Conn) {
	g.mutex.Lock()
	delete(g.Clients, client)
	g.mutex.Unlock()
}

func (g *gateway) GetRealTimeBars(symbol string) (CandleStick, error) {
	candle, err := g.DataProvider.GetRealTimeBars(symbol)
	if err != nil {
		return CandleStick{}, err
	}

	serializedData, err := json.Marshal(candle)
	if err != nil {
		log.Println("Failed to serialize market data:", err)
		return CandleStick{}, err
	}
	g.Redis.Set(context.Background(), "realtimebars:"+symbol, serializedData, 0)

	// TODO: store market data in postgres

	if err = g.cacheInRedis("realtimebars:"+symbol, serializedData); err != nil {
		return CandleStick{}, err
	}

	return candle, nil
}

func (g *gateway) GetOrderBook(symbol string) (OrderBook, error) {
	orderBook, err := g.DataProvider.GetOrderBook(symbol)
	if err != nil {
		return OrderBook{}, err
	}

	serializedData, err := json.Marshal(orderBook)
	if err != nil {
		log.Println("Failed to serialize order book:", err)
		return OrderBook{}, err
	}

	// TODO: store order book in postgres

	if err = g.cacheInRedis("orderbook:"+symbol, serializedData); err != nil {
		return OrderBook{}, err
	}
	return orderBook, nil
}

func (g *gateway) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		http.Error(w, "Failed to upgrade", http.StatusInternalServerError)
		return
	}

	g.mutex.Lock()
	g.Clients[conn] = true
	g.mutex.Unlock()

	go g.handleClientMessages(conn)
}

func (g *gateway) handleClientMessages(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		g.DeregisterClient(conn)
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			var request map[string]string
			if err := json.Unmarshal(p, &request); err != nil {
				log.Println("Failed to unmarshal client message:", err)
				continue
			}
			action, exists := request["action"]
			if !exists {
				continue
			}
			if action == "request_data" {
				ticker, exists := request["ticker"]
				if !exists {
					continue
				}
				go g.sendMarketDataForTicker(conn, ticker)
			}
		}
	}
}

func (g *gateway) sendMarketDataForTicker(conn *websocket.Conn, ticker string) {
	data, err := g.GetRealTimeBars(ticker)
	if err != nil {
		log.Println("Failed to fetch market data for ticker:", ticker, err)
		return
	}
	responseData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal market data:", err)
		return
	}
	g.mutex.Lock()
	if _, exists := g.Clients[conn]; exists {
		if err := conn.WriteMessage(websocket.TextMessage, responseData); err != nil {
			log.Println("Failed to send market data to client:", err)
		}
	}
	g.mutex.Unlock()
}

func (g *gateway) initializeRedisSubscription() {
	g.subscription = g.Redis.Subscribe(context.Background(), "realtimebars_notification_channel", "orderbook_notification_channel")

	go func() {
		channel := g.subscription.Channel()
		for msg := range channel {
			g.mutex.RLock()
			for client := range g.Clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if err != nil {
					log.Println("Failed to write message to WebSocket:", err)
					client.Close()
					g.DeregisterClient(client)
				}
			}
			g.mutex.RUnlock()
		}
	}()
}
