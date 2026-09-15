CREATE OR REPLACE VIEW "crawl"."SessionsActiveView" AS
SELECT
  "id", "source_id", "status", "item_count", "error_message",
  "started_at", "finished_at"
FROM "crawl"."sessions";

COMMENT ON VIEW "crawl"."SessionsActiveView" IS 'Crawl session rows for list and get.';
