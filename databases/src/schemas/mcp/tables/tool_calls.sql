CREATE TABLE IF NOT EXISTS "mcp"."tool_calls" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "tool_name" VARCHAR(128) NOT NULL,
  "arguments" JSONB NOT NULL DEFAULT '{}'::jsonb,
  "ok" BOOLEAN NOT NULL DEFAULT true,
  "error_message" TEXT,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
