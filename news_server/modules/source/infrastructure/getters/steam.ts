/**
 * Steam 新闻源获取器
 * 从 newsnow/server/sources/steam.ts 迁移
 */
import * as cheerio from 'cheerio'
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'

export const steamGetter: SourceGetter = async () => {
  const response = await fetch('https://store.steampowered.com/stats/stats/')
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Steam: ${response.statusText}`)
  }

  const html = await response.text()
  const $ = cheerio.load(html)
  
  const news: NewsItem[] = []
  const $rows = $('#detailStats tr.player_count_row')

  $rows.each((_, el) => {
    const $el = $(el)
    const $a = $el.find('a.gameLink')
    const url = $a.attr('href')
    const gameName = $a.text().trim()
    const currentPlayers = $el.find('td:first-child .currentServers').text().trim()

    if (url && gameName && currentPlayers) {
      news.push(
        NewsItem.fromEntity({
          id: url,
          title: gameName,
          url,
          pubDate: Date.now(),
          extra: {
            info: currentPlayers,
          },
        })
      )
    }
  })

  return news
}

