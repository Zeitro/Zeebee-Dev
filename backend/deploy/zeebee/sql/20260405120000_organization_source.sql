-- Lender org provenance: Manual (admin UI create) vs Search (OpenAI discover).
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS source VARCHAR(16) NOT NULL DEFAULT 'Manual';

COMMENT ON COLUMN organizations.source IS 'Manual | Search';
