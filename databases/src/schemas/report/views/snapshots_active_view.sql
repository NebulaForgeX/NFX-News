CREATE OR REPLACE VIEW "report"."SnapshotsActiveView" AS
SELECT
  "id", "mode", "title", "payload", "item_count", "created_at"
FROM "report"."snapshots";

COMMENT ON VIEW "report"."SnapshotsActiveView" IS 'Generated report snapshot rows.';
