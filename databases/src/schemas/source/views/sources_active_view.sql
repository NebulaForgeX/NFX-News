CREATE OR REPLACE VIEW "source"."SourcesActiveView" AS
SELECT
  "id", "name", "title", "column_key", "home", "color",
  "interval_ms", "source_type", "redirect", "metadata",
  "created_at", "updated_at"
FROM "source"."sources";

COMMENT ON VIEW "source"."SourcesActiveView" IS 'Catalog of news sources.';
