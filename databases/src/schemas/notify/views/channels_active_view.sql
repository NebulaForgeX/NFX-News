CREATE OR REPLACE VIEW "notify"."ChannelsActiveView" AS
SELECT
  "id", "account_id", "profile_id", "kind", "name", "enabled", "config", "created_at", "updated_at"
FROM "notify"."channels";

COMMENT ON VIEW "notify"."ChannelsActiveView" IS 'Notify channel rows.';
