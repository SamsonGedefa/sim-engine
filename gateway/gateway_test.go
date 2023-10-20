package gateway

// type MockCQG struct {
// 	mock.Mock
// }

// func (m *MockCQG) GetRealTimeBars(symbol string) (CandleStick, error) {
// 	args := m.Called(symbol)
// 	return args.Get(0).(CandleStick), args.Error(1)
// }

// func (m *MockCQG) GetOrderBook(symbol string) (OrderBook, error) {
// 	args := m.Called(symbol)
// 	return args.Get(0).(OrderBook), args.Error(1)
// }

// type MockRedisClient struct {
// 	mock.Mock
// }

// func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) {
// 	m.Called(key, value, expiration)
// }

// func (m *MockRedisClient) Publish(channel string, message interface{}) {
// 	m.Called(channel, message)
// }

// func TestGetRealTimeBars(t *testing.T) {
// 	mockCQG := new(MockCQG)
// 	mockRedis := new(MockRedisClient)
// 	gateway := NewGateway(mockCQG, mockRedis, "")

// 	mockCandleStick := CandleStick{
// 		Symbol:    "AAPL",
// 		Open:      130.8,
// 		High:      135,
// 		Low:       129,
// 		Close:     134,
// 		Volume:    5000,
// 		Timestamp: time.Now().Unix(),
// 	}

// 	mockCQG.On("GetRealTimeBars", "AAPL").Return(mockCandleStick, nil)

// 	marketData, err := gateway.GetRealTimeBars("AAPL")

// 	assert.NoError(t, err, "Expected no error getting real-time bars")
// 	assert.Equal(t, "AAPL", marketData.Symbol, "Expected symbol AAPL")

// 	mockCQG.AssertExpectations(t)
// 	mockRedis.AssertExpectations(t)
// }
