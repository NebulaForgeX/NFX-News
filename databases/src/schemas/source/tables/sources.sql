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
