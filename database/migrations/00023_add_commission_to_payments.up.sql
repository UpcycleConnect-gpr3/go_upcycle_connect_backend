ALTER TABLE PAYMENTS
    ADD COLUMN commission_cents INT DEFAULT 0 AFTER amount_cents,
    ADD COLUMN net_cents INT DEFAULT 0 AFTER commission_cents
-- UpcycleConnect keeps a 10% commission (commission_cents); net_cents is what is
-- owed to the seller. Bookkeeping only: the full amount lands on the platform
-- Stripe account, so the seller payout is handled separately (or via Connect).
