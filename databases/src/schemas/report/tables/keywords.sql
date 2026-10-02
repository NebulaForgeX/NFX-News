CREATE TABLE IF NOT EXISTS "report"."keywords" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "account_id" UUID,
  "profile_id" UUID,
  "group_name" VARCHAR(128) NOT NULL DEFAULT 'default',
  "word" VARCHAR(255) NOT NULL,
  "kind" VARCHAR(16) NOT NULL DEFAULT 'include',
  "count_limit" INTEGER NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_report_keywords_group" ON "report"."keywords"("group_name");
CREATE INDEX IF NOT EXISTS "idx_report_keywords_account" ON "report"."keywords"("account_id", "profile_id");
