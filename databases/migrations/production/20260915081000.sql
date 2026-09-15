-- pgcrypto.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto"
WITH
  SCHEMA "public" VERSION "1.3";

-- btree_gist.sql
-- Active: 1768427124487@@192.168.1.64@10105@postgres
CREATE EXTENSION IF NOT EXISTS "btree_gist"
WITH
  SCHEMA "public";

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "source";
COMMENT ON SCHEMA "source" IS 'News source catalog';

-- sources.sql
CREATE TABLE IF NOT EXISTS "source"."sources" (
  "id" VARCHAR(64) PRIMARY KEY,
  "name" VARCHAR(128) NOT NULL,
  "title" VARCHAR(255),
  "column_key" VARCHAR(64) NOT NULL DEFAULT 'tech',
  "home" VARCHAR(512),
  "color" VARCHAR(32),
  "interval_ms" INTEGER NOT NULL DEFAULT 600000,
  "source_type" VARCHAR(32) NOT NULL DEFAULT 'hottest',
  "redirect" VARCHAR(64),
  "metadata" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- sources_active_view.sql
CREATE OR REPLACE VIEW "source"."SourcesActiveView" AS
SELECT
  "id", "name", "title", "column_key", "home", "color",
  "interval_ms", "source_type", "redirect", "metadata",
  "created_at", "updated_at"
FROM "source"."sources";

COMMENT ON VIEW "source"."SourcesActiveView" IS 'Catalog of news sources.';

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "news";
COMMENT ON SCHEMA "news" IS 'Persisted news items';

-- items.sql
CREATE TABLE IF NOT EXISTS "news"."items" (
  "id" VARCHAR(255) PRIMARY KEY,
  "source_id" VARCHAR(64) NOT NULL,
  "original_id" VARCHAR(255) NOT NULL,
  "title" TEXT NOT NULL,
  "url" TEXT NOT NULL,
  "mobile_url" TEXT,
  "pub_date" TIMESTAMP,
  "extra" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_news_items_source_id" ON "news"."items"("source_id");
CREATE INDEX IF NOT EXISTS "idx_news_items_updated_at" ON "news"."items"("updated_at" DESC);
CREATE INDEX IF NOT EXISTS "idx_news_items_title" ON "news"."items" USING gin (to_tsvector('simple', "title"));

-- profile_preferences.sql
CREATE TABLE IF NOT EXISTS "news"."profile_preferences" (
  "account_id" UUID NOT NULL,
  "profile_id" UUID NOT NULL,
  "column_order" JSONB NOT NULL DEFAULT '[]'::jsonb,
  "payload" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("account_id", "profile_id")
);

-- items_active_view.sql
CREATE OR REPLACE VIEW "news"."ItemsActiveView" AS
SELECT
  "id", "source_id", "original_id", "title", "url", "mobile_url",
  "pub_date", "extra", "created_at", "updated_at"
FROM "news"."items";

COMMENT ON VIEW "news"."ItemsActiveView" IS 'News items for list and search reads.';

-- profile_preferences_active_view.sql
CREATE OR REPLACE VIEW "news"."ProfilePreferencesActiveView" AS
SELECT
  "account_id", "profile_id", "column_order", "payload", "updated_at"
FROM "news"."profile_preferences";

COMMENT ON VIEW "news"."ProfilePreferencesActiveView" IS 'Reader preferences keyed by Identity account+profile.';

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "crawl";
COMMENT ON SCHEMA "crawl" IS 'Crawl sessions and scheduling';

-- sessions.sql
CREATE TABLE IF NOT EXISTS "crawl"."sessions" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "source_id" VARCHAR(64),
  "status" VARCHAR(32) NOT NULL DEFAULT 'running',
  "item_count" INTEGER NOT NULL DEFAULT 0,
  "error_message" TEXT,
  "started_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "finished_at" TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_crawl_sessions_started_at" ON "crawl"."sessions"("started_at" DESC);

-- sessions_active_view.sql
CREATE OR REPLACE VIEW "crawl"."SessionsActiveView" AS
SELECT
  "id", "source_id", "status", "item_count", "error_message",
  "started_at", "finished_at"
FROM "crawl"."sessions";

COMMENT ON VIEW "crawl"."SessionsActiveView" IS 'Crawl session rows for list and get.';

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "report";
COMMENT ON SCHEMA "report" IS 'Keyword filters and report snapshots';

-- keywords.sql
CREATE TABLE IF NOT EXISTS "report"."keywords" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "group_name" VARCHAR(128) NOT NULL DEFAULT 'default',
  "word" VARCHAR(255) NOT NULL,
  "kind" VARCHAR(16) NOT NULL DEFAULT 'include',
  "count_limit" INTEGER NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_report_keywords_group" ON "report"."keywords"("group_name");

-- snapshots.sql
CREATE TABLE IF NOT EXISTS "report"."snapshots" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "mode" VARCHAR(32) NOT NULL,
  "title" VARCHAR(255) NOT NULL,
  "payload" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "item_count" INTEGER NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_report_snapshots_created_at" ON "report"."snapshots"("created_at" DESC);

-- keywords_active_view.sql
CREATE OR REPLACE VIEW "report"."KeywordsActiveView" AS
SELECT
  "id", "group_name", "word", "kind", "count_limit", "created_at"
FROM "report"."keywords";

COMMENT ON VIEW "report"."KeywordsActiveView" IS 'Report keyword rows.';

-- snapshots_active_view.sql
CREATE OR REPLACE VIEW "report"."SnapshotsActiveView" AS
SELECT
  "id", "mode", "title", "payload", "item_count", "created_at"
FROM "report"."snapshots";

COMMENT ON VIEW "report"."SnapshotsActiveView" IS 'Generated report snapshot rows.';

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "notify";
COMMENT ON SCHEMA "notify" IS 'Notification channels and deliveries';

-- channels.sql
CREATE TABLE IF NOT EXISTS "notify"."channels" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "kind" VARCHAR(32) NOT NULL,
  "name" VARCHAR(128) NOT NULL,
  "enabled" BOOLEAN NOT NULL DEFAULT true,
  "config" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- deliveries.sql
CREATE TABLE IF NOT EXISTS "notify"."deliveries" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "channel_id" UUID NOT NULL,
  "report_id" UUID,
  "status" VARCHAR(32) NOT NULL DEFAULT 'pending',
  "error_message" TEXT,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "sent_at" TIMESTAMP
);

-- channels_active_view.sql
CREATE OR REPLACE VIEW "notify"."ChannelsActiveView" AS
SELECT
  "id", "kind", "name", "enabled", "config", "created_at", "updated_at"
FROM "notify"."channels";

COMMENT ON VIEW "notify"."ChannelsActiveView" IS 'Notify channel rows.';

-- deliveries_active_view.sql
CREATE OR REPLACE VIEW "notify"."DeliveriesActiveView" AS
SELECT
  "id", "channel_id", "report_id", "status", "error_message",
  "created_at", "sent_at"
FROM "notify"."deliveries";

COMMENT ON VIEW "notify"."DeliveriesActiveView" IS 'Notify delivery rows.';

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "mcp";
COMMENT ON SCHEMA "mcp" IS 'MCP tool invocation audit';

-- tool_calls.sql
CREATE TABLE IF NOT EXISTS "mcp"."tool_calls" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "tool_name" VARCHAR(128) NOT NULL,
  "arguments" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "ok" BOOLEAN NOT NULL DEFAULT true,
  "error_message" TEXT,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- tool_calls_active_view.sql
CREATE OR REPLACE VIEW "mcp"."ToolCallsActiveView" AS
SELECT
  "id", "tool_name", "arguments", "ok", "error_message", "created_at"
FROM "mcp"."tool_calls";

COMMENT ON VIEW "mcp"."ToolCallsActiveView" IS 'MCP tool invocation log.';

-- schema.sql
CREATE SCHEMA IF NOT EXISTS "system";
COMMENT ON SCHEMA "system" IS 'System-level state';

-- system_state.sql
CREATE TABLE IF NOT EXISTS "system"."system_state" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "initialized" BOOLEAN NOT NULL DEFAULT false,
  "initialized_at" TIMESTAMP,
  "initialization_version" VARCHAR(50),
  "reset_count" INTEGER NOT NULL DEFAULT 0,
  "metadata" JSONB DEFAULT '{}'::jsonb,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_system_state_created_at" ON "system"."system_state"("created_at" DESC);

-- system_state_active_view.sql
CREATE OR REPLACE VIEW "system"."SystemStateActiveView" AS
SELECT
  "id", "initialized", "initialized_at", "initialization_version",
  "reset_count", "metadata", "created_at", "updated_at"
FROM "system"."system_state";

COMMENT ON VIEW "system"."SystemStateActiveView" IS 'System bootstrap state rows.';
