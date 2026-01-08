/**
 * 36氪新闻源获取器
 * 从 newsnow/server/sources/_36kr.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'
import { parseRelativeDate } from './utils/date'

export const _36krQuickGetter: SourceGetter = async () => {
  const baseURL = 'https://www.36kr.com'
  const url = `${baseURL}/newsflashes`
  const response = await fetch(url, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch 36kr: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  
  const news: NewsItem[] = []
  const $items = $('.newsflash-item')
  
  $items.each((_, el) => {
    const $el = $(el)
    const $a = $el.find('a.item-title')
    const url = $a.attr('href')
    const title = $a.text() // 与 newsnow 保持一致，不使用 trim
    const relativeDate = $el.find('.time').text() // 与 newsnow 保持一致，不使用 trim
    
    if (url && title && relativeDate) {
      const parsedDate = parseRelativeDate(relativeDate, 'Asia/Shanghai')
      const dateValue = parsedDate instanceof Date ? parsedDate.valueOf() : new Date(parsedDate).valueOf()
      
      news.push(
        NewsItem.fromEntity({
          id: url,
          title,
          url: `${baseURL}${url}`,
          extra: {
            date: dateValue,
          },
        })
      )
    }
  })

  return news
}

