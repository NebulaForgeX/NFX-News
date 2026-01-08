/**
 * 金十数据新闻源获取器
 * 从 newsnow/server/sources/jin10.ts 迁移
 */
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'
import { parseRelativeDate } from './utils/date'

interface Jin10Item {
  id: string
  time: string
  type: number
  data: {
    pic?: string
    title?: string
    source?: string
    content?: string
    source_link?: string
    vip_title?: string
    lock?: boolean
    vip_level?: number
    vip_desc?: string
  }
  important: number
  tags: string[]
  channel: number[]
  remark: any[]
}

export const jin10Getter: SourceGetter = async () => {
  const timestamp = Date.now()
  const url = `https://www.jin10.com/flash_newest.js?t=${timestamp}`

  const response = await fetch(url, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Jin10: ${response.statusText}`)
  }

  const rawData = await response.text()

  // 解析 JavaScript 变量（与 newsnow 保持一致）
  const jsonStr = rawData
    .replace(/^var\s+newest\s*=\s*/, '') // 移除开头的变量声明
    .replace(/;*$/, '') // 移除末尾可能存在的分号
    .trim() // 移除首尾空白字符

  let data: Jin10Item[]
  try {
    data = JSON.parse(jsonStr) as Jin10Item[]
  } catch (error) {
    throw new Error(`无法解析金十数据：${error instanceof Error ? error.message : 'JSON解析失败'}`)
  }

  return data
    .filter((k) => (k.data.title || k.data.content) && !k.channel?.includes(5))
    .map((k) => {
      const text = (k.data.title || k.data.content)!.replace(/<\/?b>/g, '')
      // 使用与 newsnow 相同的解构赋值方式
      const [, title, desc] = text.match(/^【([^】]*)】(.*)$/) ?? []
      
      return NewsItem.fromEntity({
        id: k.id,
        title: title ?? text,
        url: `https://flash.jin10.com/detail/${k.id}`,
        pubDate: parseRelativeDate(k.time, 'Asia/Shanghai').valueOf(),
        extra: {
          hover: desc,
          // 使用与 newsnow 相同的逻辑：!!k.important && "✰"
          info: !!k.important && '✰',
        },
      })
    })
}

