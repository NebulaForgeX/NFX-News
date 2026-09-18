CREATE OR REPLACE VIEW "report"."KeywordsActiveView" AS
SELECT
  "id", "account_id", "profile_id", "group_name", "word", "kind", "count_limit", "created_at"
FROM "report"."keywords";

COMMENT ON VIEW "report"."KeywordsActiveView" IS 'Report keyword rows.';
