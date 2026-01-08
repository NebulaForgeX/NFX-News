/**
 * 法布财经新闻源获取器
 * 从 newsnow/server/sources/fastbull.ts 迁移
 * 修复：完全按照原始实现，去除多余的 trim()，确保解析逻辑一致
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'

export const fastbullExpressGetter: SourceGetter = async () => {
  const baseURL = 'https://www.fastbull.com'
  
  // 使用 fetch 获取 HTML，与 newsnow 的 myFetch 行为保持一致
  const response = await fetch(`${baseURL}/cn/express-news`, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Fastbull express: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  const $main = $('.news-list')
  const news: NewsItem[] = []

  // 完全按照原始 newsnow 实现：不使用 trim()，直接使用 text()
  $main.each((_, el) => {
    const $el = $(el)
    const a = $el.find('.title_name')
    const url = a.attr('href')
    const titleText = a.text()  // 原始代码没有 trim()
    const title = titleText.match(/【(.+)】/)?.[1] ?? titleText  // 使用可选链和空值合并
    const date = $el.attr('data-date')

    if (url && title && date) {
      news.push(
        NewsItem.fromEntity({
          id: url,
          title: title.length < 4 ? titleText : title,
          url: baseURL + url,
          pubDate: Number(date),
        })
      )
    }
  })

  return news
}

export const fastbullNewsGetter: SourceGetter = async () => {
  const baseURL = 'https://www.fastbull.com'
  
  const response = await fetch(`${baseURL}/cn/news`, {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Fastbull news: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  const $main = $('.trending_type')
  const news: NewsItem[] = []

  // 完全按照原始 newsnow 实现：不使用 trim()
  $main.each((_, el) => {
    const a = $(el)
    const url = a.attr('href')
    const title = a.find('.title').text()  // 原始代码没有 trim()
    const date = a.find('[data-date]').attr('data-date')

    if (url && title && date) {
      news.push(
        NewsItem.fromEntity({
          id: url,
          title,
          url: baseURL + url,
          pubDate: Number(date),
        })
      )
    }
  })

  return news
}

