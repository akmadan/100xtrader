CREATE TABLE IF NOT EXISTS algorithms (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    code TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft',
    symbol TEXT NOT NULL DEFAULT '',
    timeframe TEXT NOT NULL DEFAULT '1m',
    execution_mode TEXT NOT NULL DEFAULT 'paper_trading',
    broker TEXT,
    enabled INTEGER NOT NULL DEFAULT 0,
    config TEXT DEFAULT '{}',
    state TEXT DEFAULT '{}',
    last_run_at TEXT,
    last_signal TEXT,
    total_trades INTEGER NOT NULL DEFAULT 0,
    win_rate REAL NOT NULL DEFAULT 0.0,
    total_pnl REAL NOT NULL DEFAULT 0.0,
    version INTEGER NOT NULL DEFAULT 1,
    tags TEXT DEFAULT '[]',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_algorithms_user_id ON algorithms(user_id);
CREATE INDEX IF NOT EXISTS idx_algorithms_status ON algorithms(status);
CREATE INDEX IF NOT EXISTS idx_algorithms_enabled ON algorithms(enabled);
CREATE INDEX IF NOT EXISTS idx_algorithms_symbol ON algorithms(symbol);

