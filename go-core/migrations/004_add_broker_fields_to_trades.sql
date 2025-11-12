-- Add broker-specific fields to trades table
-- This migration adds columns for storing broker-specific trade information

-- Note: SQLite doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN
-- If columns already exist, this migration should be marked as applied manually
-- or the migration runner should be updated to handle column existence checks

-- Add broker-specific columns
ALTER TABLE trades ADD COLUMN trading_broker TEXT;
ALTER TABLE trades ADD COLUMN trader_broker_id TEXT;
ALTER TABLE trades ADD COLUMN exchange_order_id TEXT;
ALTER TABLE trades ADD COLUMN order_id TEXT;
ALTER TABLE trades ADD COLUMN product_type TEXT;
ALTER TABLE trades ADD COLUMN transaction_type TEXT;

-- Create indexes for faster lookups by broker ID
CREATE INDEX IF NOT EXISTS idx_trades_trading_broker ON trades(trading_broker);
CREATE INDEX IF NOT EXISTS idx_trades_exchange_order_id ON trades(exchange_order_id);
CREATE INDEX IF NOT EXISTS idx_trades_order_id ON trades(order_id);
