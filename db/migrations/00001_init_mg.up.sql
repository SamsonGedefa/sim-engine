CREATE TABLE market_data (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    high FLOAT8 NOT NULL,
    low FLOAT8 NOT NULL,
    open FLOAT8 NOT NULL,
    close FLOAT8 NOT NULL,
    volume INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

CREATE TABLE order_books (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    bid_price FLOAT8 NOT NULL,
    bid_size INTEGER NOT NULL,
    ask_price FLOAT8 NOT NULL,
    ask_size INTEGER NOT NULL
);