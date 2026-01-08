/**
 * News 映射器
 * 将数据库模型转换为领域模型
 */
import type { NewsRO } from '../../../domain/news/read_model'
import type { News as NewsModel } from '@models/news/news'
import { consola } from 'consola'

/**
 * 将数据库模型转换为读模型
 */
export function ToNewsRO(model: NewsModel): NewsRO {
  // 处理 extra 字段：可能是字符串或已经是对象
  let extra: any = undefined
  if (model.extra) {
    if (typeof model.extra === 'string') {
      try {
        extra = JSON.parse(model.extra)
      } catch (error) {
        // 如果解析失败，记录错误但继续处理
        consola.warn('Failed to parse extra field as JSON:', error)
        extra = undefined
      }
    } else {
      // 如果已经是对象，直接使用
      extra = model.extra
    }
  }

  return {
    id: model.id,
    sourceId: model.sourceId,
    originalId: model.originalId,
    title: model.title,
    url: model.url,
    mobileUrl: model.mobileUrl ?? undefined,
    pubDate: model.pubDate ?? undefined,
    extra,
    createdAt: model.createdAt,
    updatedAt: model.updatedAt,
  }
}

