CREATE TABLE IF NOT EXISTS "notify"."channels" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "account_id" UUID,
  "profile_id" UUID,
  "kind" VARCHAR(32) NOT NULL,
  "name" VARCHAR(128) NOT NULL,
  "enabled" BOOLEAN NOT NULL DEFAULT true,
  "config" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_notify_channels_account" ON "notify"."channels"("account_id", "profile_id");
