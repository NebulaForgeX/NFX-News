package mcp

import "nfxnews/pkgs/errx"

var (
	ErrMCPToolUnknown     = errx.InvalidArg("MCP_TOOL_UNKNOWN", "unknown MCP tool")
	ErrMCPToolFailed      = errx.Internal("MCP_TOOL_FAILED", "MCP tool execution failed")
	ErrMCPArgsInvalid     = errx.InvalidArg("MCP_ARGS_INVALID", "invalid MCP tool arguments")
	ErrMCPClientMissing   = errx.Internal("MCP_CLIENT_MISSING", "required gRPC client is unavailable")
)

/*
!MCP_TOOL_UNKNOWN
*en<unknown MCP tool>
*zh<未知 MCP 工具>
*fr<outil MCP inconnu>

!MCP_TOOL_FAILED
*en<MCP tool execution failed>
*zh<MCP 工具执行失败>
*fr<échec de l'exécution de l'outil MCP>

!MCP_ARGS_INVALID
*en<invalid MCP tool arguments>
*zh<MCP 工具参数无效>
*fr<arguments d'outil MCP invalides>

!MCP_CLIENT_MISSING
*en<required gRPC client is unavailable>
*zh<所需 gRPC 客户端不可用>
*fr<client gRPC requis indisponible>
*/
