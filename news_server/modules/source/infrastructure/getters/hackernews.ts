/**
 * Hacker News 新闻源获取器
 * 从 newsnow/server/sources/hackernews.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'

export const hackernewsGetter: SourceGetter = async () => {
  const baseURL = 'https://news.ycombinator.com'
  const response = await fetch(baseURL)
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Hacker News: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  
  const news: NewsItem[] = []
  const $main = $('.athing')
  
  $main.each((_, el) => {
    const $el = $(el)
    const $a = $el.find('.titleline a').first()
    const title = $a.text().trim()
    const id = $el.attr('id')
    const score = $(`#score_${id}`).text().trim()
    const url = `${baseURL}/item?id=${id}`
    
    if (id && title) {
      news.push(
        NewsItem.fromEntity({
          id,
          title,
          url,
          extra: {
            info: score || undefined,
          },
        })
      )
    }
  })
  
  return news
}

