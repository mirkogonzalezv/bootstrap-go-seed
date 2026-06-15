CREATE TABLE IF NOT EXISTS health_checks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status VARCHAR(50) NOT NULL DEFAULT 'ok',
    checked_at TIMESTAMPZ NOT NULL DEFAULT NOW()
)