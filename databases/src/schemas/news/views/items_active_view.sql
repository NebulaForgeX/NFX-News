CREATE OR REPLACE VIEW "news"."ItemsActiveView" AS
SELECT
  "id", "source_id", "original_id", "title", "url", "mobile_url",
  "pub_date", "extra", "created_at", "updated_at"
FROM "news"."items";

COMMENT ON VIEW "news"."ItemsActiveView" IS 'News items for list and search reads.';
