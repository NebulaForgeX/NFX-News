CREATE OR REPLACE VIEW "mcp"."ToolCallsActiveView" AS
SELECT
  "id", "tool_name", "arguments", "ok", "error_message", "created_at"
FROM "mcp"."tool_calls";

COMMENT ON VIEW "mcp"."ToolCallsActiveView" IS 'MCP tool invocation log.';
