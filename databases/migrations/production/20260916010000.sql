-- tenant columns for crawl / report / notify
ALTER TABLE "crawl"."sessions" ADD COLUMN IF NOT EXISTS "account_id" UUID;
ALTER TABLE "crawl"."sessions" ADD COLUMN IF NOT EXISTS "profile_id" UUID;
CREATE INDEX IF NOT EXISTS "idx_crawl_sessions_account" ON "crawl"."sessions"("account_id", "profile_id");

ALTER TABLE "report"."keywords" ADD COLUMN IF NOT EXISTS "account_id" UUID;
ALTER TABLE "report"."keywords" ADD COLUMN IF NOT EXISTS "profile_id" UUID;
CREATE INDEX IF NOT EXISTS "idx_report_keywords_account" ON "report"."keywords"("account_id", "profile_id");

ALTER TABLE "report"."snapshots" ADD COLUMN IF NOT EXISTS "account_id" UUID;
ALTER TABLE "report"."snapshots" ADD COLUMN IF NOT EXISTS "profile_id" UUID;
CREATE INDEX IF NOT EXISTS "idx_report_snapshots_account" ON "report"."snapshots"("account_id", "profile_id");

ALTER TABLE "notify"."channels" ADD COLUMN IF NOT EXISTS "account_id" UUID;
ALTER TABLE "notify"."channels" ADD COLUMN IF NOT EXISTS "profile_id" UUID;
CREATE INDEX IF NOT EXISTS "idx_notify_channels_account" ON "notify"."channels"("account_id", "profile_id");

ALTER TABLE "notify"."deliveries" ADD COLUMN IF NOT EXISTS "account_id" UUID;
ALTER TABLE "notify"."deliveries" ADD COLUMN IF NOT EXISTS "profile_id" UUID;
CREATE INDEX IF NOT EXISTS "idx_notify_deliveries_account" ON "notify"."deliveries"("account_id", "profile_id");

DROP VIEW IF EXISTS "crawl"."SessionsActiveView";
CREATE VIEW "crawl"."SessionsActiveView" AS
SELECT
  "id", "account_id", "profile_id", "source_id", "status", "item_count", "error_message",
  "started_at", "finished_at"
FROM "crawl"."sessions";

DROP VIEW IF EXISTS "report"."KeywordsActiveView";
CREATE VIEW "report"."KeywordsActiveView" AS
SELECT
  "id", "account_id", "profile_id", "group_name", "word", "kind", "count_limit", "created_at"
FROM "report"."keywords";

DROP VIEW IF EXISTS "report"."SnapshotsActiveView";
CREATE VIEW "report"."SnapshotsActiveView" AS
SELECT
  "id", "account_id", "profile_id", "mode", "title", "payload", "item_count", "created_at"
FROM "report"."snapshots";

DROP VIEW IF EXISTS "notify"."ChannelsActiveView";
CREATE VIEW "notify"."ChannelsActiveView" AS
SELECT
  "id", "account_id", "profile_id", "kind", "name", "enabled", "config", "created_at", "updated_at"
FROM "notify"."channels";

DROP VIEW IF EXISTS "notify"."DeliveriesActiveView";
CREATE VIEW "notify"."DeliveriesActiveView" AS
SELECT
  "id", "account_id", "profile_id", "channel_id", "report_id", "status", "error_message",
  "created_at", "sent_at"
FROM "notify"."deliveries";
