-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user TEXT NOT NULL,
    environment TEXT NOT NULL,
    ticker TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY,
    user TEXT NOT NULL,
    symbol TEXT NOT NULL,
    side TEXT NOT NULL, -- buy or sell
    type TEXT NOT NULL, -- market, limit, stop
    quantity REAL NOT NULL,
    price REAL NOT NULL,
    status TEXT NOT NULL,
    session_id TEXT,
    source TEXT, -- user or ai
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trades table
CREATE TABLE IF NOT EXISTS trades (
    id INTEGER PRIMARY KEY,
    buy_order_id INTEGER NOT NULL,
    sell_order_id INTEGER NOT NULL,
    symbol TEXT NOT NULL,
    quantity REAL NOT NULL,
    price REAL NOT NULL,
    session_id TEXT,
    source TEXT, -- user or ai
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Positions table
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY,
    user TEXT NOT NULL,
    symbol TEXT NOT NULL,
    quantity REAL NOT NULL,
    average_price REAL NOT NULL
);

-- Environments table
CREATE TABLE IF NOT EXISTS environments (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    volatility TEXT,
    trend TEXT,
    liquidity TEXT
    -- Add more columns as needed
);

-- Tickers table
CREATE TABLE IF NOT EXISTS tickers (
    symbol TEXT PRIMARY KEY,
    name TEXT NOT NULL
    -- Add more columns as needed
);
