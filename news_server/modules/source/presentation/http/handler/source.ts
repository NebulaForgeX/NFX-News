/**
 * Source Handler
 * 类似 Sjgz-Backend 的 interfaces/http/handler/source.go
 */
import type { FastifyRequest, FastifyReply } from 'fastify'
import { GetterService } from '../../../application/getter/service'
import type { GetSourceDataQuery } from '../dto/req.dto'
import type { SourceDataResponse, ErrorResponse } from '../dto/res.dto'

export class SourceHandler {
  constructor(private readonly getterService: GetterService) {}

  /**
   * 获取源数据
   * 始终返回固定的JSON数据结构：
   * - 成功: { status: 'success' | 'cache', id: string, updatedTime: number, items: NewsItem[] }
   * - 错误: { error: string }
   */
  async getSourceData(
    request: FastifyRequest<{ Querystring: GetSourceDataQuery }>,
    reply: FastifyReply
  ) {
    const { id, latest } = request.query

    // 确保始终返回JSON格式
    reply.type('application/json')

    if (!id) {
      return reply.code(400).send({ error: 'Missing id parameter' } as ErrorResponse)
    }

    try {
      const result = await this.getterService.getSourceData.execute({
        sourceId: id,
        latest: latest === 'true',
      })

      // 确保返回的数据结构完整
      const response: SourceDataResponse = {
        status: result.status || 'success',
        id: result.id,
        updatedTime: result.updatedTime,
        items: result.items || [], // 确保items始终是数组
      }

      return reply.code(200).send(response)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Internal Server Error'
      // 错误时也返回JSON格式
      return reply.code(500).send({ error: message } as ErrorResponse)
    }
  }
}

