/**
 * GitHub 新闻源获取器
 * 从 newsnow/server/sources/github.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'

export const githubTrendingGetter: SourceGetter = async () => {
  const baseURL = 'https://github.com'
  const response = await fetch('https://github.com/trending?spoken_language_code=')
  
  if (!response.ok) {
    throw new Error(`Failed to fetch GitHub trending: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  
  const news: NewsItem[] = []
  const $main = $('main .Box div[data-hpc] > article')
  
  $main.each((_, el) => {
    const $el = $(el)
    const $a = $el.find('>h2 a')
    const title = $a.text().replace(/\n+/g, '').trim()
    const url = $a.attr('href')
    const star = $el.find('[href$=stargazers]').text().replace(/\s+/g, '').trim()
    const desc = $el.find('>p').text().replace(/\n+/g, '').trim()
    
    if (url && title) {
      news.push(
        NewsItem.fromEntity({
          id: url,
          title,
          url: `${baseURL}${url}`,
          extra: {
            info: `✰ ${star}`,
            hover: desc,
          },
        })
      )
    }
  })
  
  return news
}

