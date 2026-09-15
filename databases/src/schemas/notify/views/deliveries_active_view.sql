CREATE OR REPLACE VIEW "notify"."DeliveriesActiveView" AS
SELECT
  "id", "channel_id", "report_id", "status", "error_message",
  "created_at", "sent_at"
FROM "notify"."deliveries";

COMMENT ON VIEW "notify"."DeliveriesActiveView" IS 'Notify delivery rows.';
