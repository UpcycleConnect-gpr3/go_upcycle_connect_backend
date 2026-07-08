CREATE TABLE IF NOT EXISTS SUBSCRIPTIONS (
    id VARCHAR(36) PRIMARY KEY,
    user_id CHAR(36),
    stripe_session_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    price_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_stripe_session (stripe_session_id),
    KEY idx_user (user_id)
)
-- Note: user_id references a user owned by the auth service (separate DB),
-- so no cross-service foreign key is enforced here.
