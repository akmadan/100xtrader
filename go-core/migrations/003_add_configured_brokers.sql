-- Migration to add configured_brokers field to users table
-- This field stores broker configurations (like Dhan access tokens) as JSON

ALTER TABLE users ADD COLUMN configured_brokers TEXT DEFAULT '{}';

