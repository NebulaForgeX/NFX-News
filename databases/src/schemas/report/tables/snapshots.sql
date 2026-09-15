CREATE TABLE IF NOT EXISTS "report"."snapshots" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "mode" VARCHAR(32) NOT NULL,
  "title" VARCHAR(255) NOT NULL,
  "payload" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "item_count" INTEGER NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_report_snapshots_created_at" ON "report"."snapshots"("created_at" DESC);
