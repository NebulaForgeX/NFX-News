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
