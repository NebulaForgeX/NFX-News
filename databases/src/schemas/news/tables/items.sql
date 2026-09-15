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
