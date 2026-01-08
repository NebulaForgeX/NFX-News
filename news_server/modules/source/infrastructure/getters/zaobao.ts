/**
 * 联合早报新闻源获取器
 * 从 newsnow/server/sources/zaobao.ts 迁移
 * 
 * 编码处理：
 * - 网站返回的是 GB2312/GBK 编码
 * - 优先使用 GBK 解码（GBK 是 GB2312 的超集，能更好地处理中文）
 * - 如果失败，依次尝试 GB2312 和 GB18030
 * - 最终输出 UTF-8 格式，确保前端正常显示中文
 */
import { Buffer } from 'node:buffer'
import * as cheerio from 'cheerio'
import iconv from 'iconv-lite'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'
import { parseRelativeDate } from './utils/date'

export const zaobaoGetter: SourceGetter = async () => {
  const response = await fetch('https://www.zaochenbao.com/realtime/', {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Zaobao: ${response.statusText}`)
  }

  // 获取 ArrayBuffer 以正确处理编码
  const arrayBuffer = await response.arrayBuffer()
  const base = 'https://www.zaochenbao.com'
  
  // 使用 GBK 编码解码（GBK 是 GB2312 的超集，能更好地处理中文）
  // 优先尝试 GBK，因为它能处理更多字符，包括 GB2312 的所有字符
  let utf8String: string
  try {
    utf8String = iconv.decode(Buffer.from(arrayBuffer), 'gbk')
  } catch (gbkError) {
    // 如果 GBK 失败，尝试 GB2312
    try {
      utf8String = iconv.decode(Buffer.from(arrayBuffer), 'gb2312')
    } catch (gb2312Error) {
      // 如果 GB2312 也失败，尝试 GB18030（最全面的中文编码）
      try {
        utf8String = iconv.decode(Buffer.from(arrayBuffer), 'gb18030')
      } catch (gb18030Error) {
        throw new Error(
          `无法解码联合早报数据：尝试了 GBK、GB2312、GB18030 均失败。` +
          `GBK错误：${gbkError instanceof Error ? gbkError.message : String(gbkError)}`
        )
      }
    }
  }
  
  // 使用 cheerio 解析 HTML，确保正确解码 HTML 实体为 UTF-8
  const $ = cheerio.load(utf8String, {
    decodeEntities: true, // 解码 HTML 实体，确保中文正确显示
  })
  const news: NewsItem[] = []
  const $main = $('div.list-block>a.item')
  
  $main.each((_, el) => {
    const $a = $(el)
    const url = $a.attr('href')
    // 与 newsnow 保持一致，不使用 trim()
    const title = $a.find('.eps')?.text()
    // 与 newsnow 保持一致，不使用 trim()
    const date = $a.find('.pdt10')?.text().replace(/-\s/g, ' ')
    
    if (url && title && date) {
      news.push(
        NewsItem.fromEntity({
          id: url,
          title,
          url: base + url,
          pubDate: parseRelativeDate(date, 'Asia/Shanghai').valueOf(),
        })
      )
    }
  })
  
  // 使用与 newsnow 相同的排序逻辑
  return news.sort((m, n) => (n.pubDate || 0) > (m.pubDate || 0) ? 1 : -1)
}

