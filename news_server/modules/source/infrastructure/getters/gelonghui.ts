/**
 * 格隆汇新闻源获取器
 * 从 newsnow/server/sources/gelonghui.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'
import { parseRelativeDate } from './utils/date'

export const gelonghuiGetter: SourceGetter = async () => {
  const baseURL = 'https://www.gelonghui.com'
  const response = await fetch('https://www.gelonghui.com/news/', {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Gelonghui: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  
  const news: NewsItem[] = []
  const $main = $('.article-content')
  
  $main.each((_, el) => {
    const $el = $(el)
    const $a = $el.find('.detail-right>a')
    const url = $a.attr('href')
    const title = $a.find('h2').text().trim()
    const info = $el.find('.time > span:nth-child(1)').text().trim()
    const relativeTime = $el.find('.time > span:nth-child(3)').text().trim()
    
    if (url && title && relativeTime) {
      news.push(
        NewsItem.fromEntity({
          id: url,
          title,
          url: baseURL + url,
          extra: {
            date: parseRelativeDate(relativeTime, 'Asia/Shanghai').valueOf(),
            info,
          },
        })
      )
    }
  })

  return news
}

