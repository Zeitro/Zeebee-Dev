-- Insert default values if lender_discover_settings was empty
INSERT INTO "public"."organization_settings" (organization_id, purpose, enabled, setting)
SELECT 1, 'LenderDiscover', false,
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
  WHERE organization_id = 1 AND purpose = 'LenderDiscover'
);
