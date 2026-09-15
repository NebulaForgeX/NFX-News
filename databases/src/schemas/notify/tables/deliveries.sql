CREATE TABLE IF NOT EXISTS "notify"."deliveries" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "channel_id" UUID NOT NULL,
  "report_id" UUID,
  "status" VARCHAR(32) NOT NULL DEFAULT 'pending',
  "error_message" TEXT,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "sent_at" TIMESTAMP
);
