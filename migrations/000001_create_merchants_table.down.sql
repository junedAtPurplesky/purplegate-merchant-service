-- Drop trigger first
DROP TRIGGER IF EXISTS update_merchants_updated_at ON merchants;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_merchants_email;
DROP INDEX IF EXISTS idx_merchants_created_at;

-- Drop merchants table
DROP TABLE IF EXISTS merchants; 