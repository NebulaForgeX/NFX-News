/**
 * 酷安新闻源获取器
 * 从 newsnow/server/sources/coolapk/index.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../../domain/getter/getter'
import { NewsItem } from '../../../domain/newsitem/entity'
import { genHeaders } from './utils'

interface CoolapkResponse {
  data: {
    id: string
    message: string
    editor_title: string
    url: string
    entityType: string
    pubDate: string
    dateline: number
    targetRow: {
      subTitle: string
    }
  }[]
}

export const coolapkGetter: SourceGetter = async () => {
  const url =
    'https://api.coolapk.com/v6/page/dataList?url=%2Ffeed%2FstatList%3FcacheExpires%3D300%26statType%3Dday%26sortField%3Ddetailnum%26title%3D%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&title=%E4%BB%8A%E6%97%A5%E7%83%AD%E9%97%A8&subTitle=&page=1'

  const response = await fetch(url, {
    headers: await genHeaders(),
  })

  if (!response.ok) {
    throw new Error(`Failed to fetch Coolapk: ${response.statusText}`)
  }

  const r = (await response.json()) as CoolapkResponse

  if (!r.data || !r.data.length) {
    throw new Error('Failed to fetch Coolapk data')
  }

  return r.data
    .filter((k) => k.id)
    .map((i) => {
      // 使用 cheerio 解析 message（可能包含 HTML）
      let title = i.editor_title
      if (!title && i.message) {
        const $ = cheerio.load(i.message)
        title = $('body').text().split('\n')[0].trim() || i.message.split('\n')[0].trim()
      }

      return NewsItem.fromEntity({
        id: i.id,
        title: title || i.id,
        url: `https://www.coolapk.com${i.url}`,
        extra: {
          info: i.targetRow?.subTitle,
        },
      })
    })
}

