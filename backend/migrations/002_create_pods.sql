
CREATE TABLE IF NOT EXISTS pods (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id  VARCHAR(255) UNIQUE NOT NULL,
    status     VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Index for faster lookups when checking for expired pods or a user's pods
CREATE INDEX IF NOT EXISTS idx_pods_status_expires ON pods(status, expires_at);
CREATE INDEX IF NOT EXISTS idx_pods_user_id ON pods(user_id);
