-- Lender org provenance: manual (admin UI create) vs search (OpenAI discover).
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS source VARCHAR(16) NOT NULL DEFAULT 'manual';

COMMENT ON COLUMN organizations.source IS 'manual | search';
