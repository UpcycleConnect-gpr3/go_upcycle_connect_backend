CREATE TABLE IF NOT EXISTS PAYMENTS (
    id VARCHAR(36) PRIMARY KEY,
    object_id VARCHAR(36),
    user_id CHAR(36),
    stripe_session_id VARCHAR(255),
    amount_cents INT DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'eur',
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_payment_session (stripe_session_id),
    KEY idx_payment_object (object_id),
    KEY idx_payment_user (user_id)
)
-- One-time payment for the purchase of an annonce (OBJECTS row). user_id is the
-- buyer. user_id / object_id are not FK-constrained: user_id belongs to the auth
-- service (separate DB), and an object may be deleted while its payment history
-- must remain for accounting.
