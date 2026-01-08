/**
 * IT之家新闻源获取器
 * 从 newsnow/server/sources/ithome.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'
import { parseRelativeDate } from './utils/date'

export const ithomeGetter: SourceGetter = async () => {
  const response = await fetch('https://www.ithome.com/list/', {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Ithome: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  
  const news: NewsItem[] = []
  const $main = $('#list > div.fl > ul > li')
  
  $main.each((_, el) => {
    const $el = $(el)
    const $a = $el.find('a.t')
    const url = $a.attr('href')
    const title = $a.text().trim()
    const date = $el.find('i').text().trim()
    
    if (url && title && date) {
      const isAd = url?.includes('lapin') || ['神券', '优惠', '补贴', '京东'].find((k) => title.includes(k))
      if (!isAd) {
        news.push(
          NewsItem.fromEntity({
            id: url,
            title,
            url,
            pubDate: parseRelativeDate(date, 'Asia/Shanghai').valueOf(),
          })
        )
      }
    }
  })
  
  return news.sort((m, n) => (n.pubDate || 0) > (m.pubDate || 0) ? 1 : -1)
}

