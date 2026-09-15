CREATE OR REPLACE VIEW "news"."ProfilePreferencesActiveView" AS
SELECT
  "account_id", "profile_id", "column_order", "payload", "updated_at"
FROM "news"."profile_preferences";

COMMENT ON VIEW "news"."ProfilePreferencesActiveView" IS 'Reader preferences keyed by Identity account+profile.';
