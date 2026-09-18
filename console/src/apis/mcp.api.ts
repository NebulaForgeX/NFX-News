import type { DataResponse } from "nfx-ui/types";

import { protectedClient } from "./clients";
import { URL_PATHS } from "./ip";

export const ListMCPTools = async (): Promise<string[]> => {
  const { data } = await protectedClient.get<DataResponse<string[]>>(URL_PATHS.MCP.tools);
  return data.data;
};

export const RunMCPTool = async (name: string, argumentsJson?: Record<string, unknown>): Promise<unknown> => {
  const { data } = await protectedClient.post<DataResponse<unknown>>(URL_PATHS.MCP.byName(name), { arguments: argumentsJson ?? {} });
  return data.data;
};

export const RunMCPJSON = async (toolName: string, args?: Record<string, unknown>): Promise<unknown> => {
  const { data } = await protectedClient.post<DataResponse<unknown>>(URL_PATHS.MCP.run, { tool_name: toolName, arguments: args ?? {} });
  return data.data;
};
