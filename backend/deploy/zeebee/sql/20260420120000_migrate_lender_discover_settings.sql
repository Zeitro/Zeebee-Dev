-- Migrate lender_discover_settings data into organization_settings
INSERT INTO "public"."organization_settings" (organization_id, purpose, enabled, setting)
SELECT 1, 'lenderDiscover', false,
  jsonb_build_object(
    'schedule_enabled', schedule_enabled,
    'interval_days', interval_days,
    'scope_state', scope_state,
    'scope_extra', scope_extra,
    'user_prompt', user_prompt,
    'openai_model', openai_model
  )
FROM "public"."lender_discover_settings"
ORDER BY id ASC LIMIT 1
ON CONFLICT DO NOTHING;

-- Insert default values if lender_discover_settings was empty (aligned with 20260404120000_lender_discover.sql seed row)
INSERT INTO "public"."organization_settings" (organization_id, purpose, enabled, setting)
SELECT 1, 'lenderDiscover', false,
  jsonb_build_object(
    'schedule_enabled', false,
    'interval_days', 7,
    'scope_state', 'CA',
    'scope_extra', 'private mortgage lenders, hard money lenders, and private lending companies operating in this state',
    'user_prompt', 'Return only factual-sounding entries; leave unknown fields empty. Output must follow the requested JSON shape.',
    'openai_model', 'gpt-4o-mini'
  )
WHERE NOT EXISTS (
  SELECT 1 FROM "public"."organization_settings"
  WHERE organization_id = 1 AND purpose = 'lenderDiscover'
);
