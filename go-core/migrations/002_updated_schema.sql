-- Updated database schema migration
-- This file contains the updated schema matching the new Go models

-- Drop old tables that are no longer needed
DROP TABLE IF EXISTS trade_screenshots;
DROP TABLE IF EXISTS trade_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS trade_journals;
DROP TABLE IF EXISTS trade_actions;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS trade_setups;

-- Update trades table to match new Trade model
DROP TABLE IF EXISTS trades;
CREATE TABLE IF NOT EXISTS trades (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    symbol TEXT NOT NULL,
    market_type TEXT NOT NULL CHECK (market_type IN ('indian', 'us', 'crypto', 'forex', 'commodities')),
    entry_date TIMESTAMP NOT NULL,
    entry_price DECIMAL NOT NULL,
    quantity INTEGER NOT NULL,
    total_amount DECIMAL NOT NULL,
    exit_price DECIMAL,
    direction TEXT NOT NULL CHECK (direction IN ('long', 'short')),
    stop_loss DECIMAL,
    target DECIMAL,
    strategy TEXT NOT NULL,
    outcome_summary TEXT NOT NULL CHECK (outcome_summary IN ('profitable', 'loss', 'breakeven', 'partial_profit', 'partial_loss')),
    trade_analysis TEXT,
    rules_followed TEXT, -- JSON array stored as TEXT
    screenshots TEXT, -- JSON array stored as TEXT
    psychology TEXT, -- JSON object stored as TEXT
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create strategies table
CREATE TABLE IF NOT EXISTS strategies (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create rules table
CREATE TABLE IF NOT EXISTS rules (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('entry', 'exit', 'stop_loss', 'take_profit', 'risk_management', 'psychology', 'other')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create mistakes table
CREATE TABLE IF NOT EXISTS mistakes (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('psychological', 'behavioral')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_trades_user_id ON trades(user_id);
CREATE INDEX IF NOT EXISTS idx_trades_market_type ON trades(market_type);
CREATE INDEX IF NOT EXISTS idx_trades_symbol ON trades(symbol);
CREATE INDEX IF NOT EXISTS idx_trades_direction ON trades(direction);
CREATE INDEX IF NOT EXISTS idx_trades_outcome_summary ON trades(outcome_summary);
CREATE INDEX IF NOT EXISTS idx_trades_created_at ON trades(created_at);

CREATE INDEX IF NOT EXISTS idx_strategies_user_id ON strategies(user_id);
CREATE INDEX IF NOT EXISTS idx_strategies_created_at ON strategies(created_at);

CREATE INDEX IF NOT EXISTS idx_rules_user_id ON rules(user_id);
CREATE INDEX IF NOT EXISTS idx_rules_category ON rules(category);
CREATE INDEX IF NOT EXISTS idx_rules_created_at ON rules(created_at);

CREATE INDEX IF NOT EXISTS idx_mistakes_user_id ON mistakes(user_id);
CREATE INDEX IF NOT EXISTS idx_mistakes_category ON mistakes(category);
CREATE INDEX IF NOT EXISTS idx_mistakes_created_at ON mistakes(created_at);
