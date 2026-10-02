CREATE TABLE IF NOT EXISTS "news"."profile_preferences" (
  "account_id" UUID NOT NULL,
  "profile_id" UUID NOT NULL,
  "column_order" JSONB NOT NULL DEFAULT '[]'::jsonb,
  "payload" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("account_id", "profile_id")
);
