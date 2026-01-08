/**
 * 快手新闻源获取器
 * 从 newsnow/server/sources/kuaishou.ts 迁移
 */
import type { SourceGetter } from '../../domain/getter/getter'
import { NewsItem } from '../../domain/newsitem/entity'

interface KuaishouResponse {
  defaultClient: {
    ROOT_QUERY: {
      'visionHotRank({"page":"home"})': {
        type: string
        id: string
        typename: string
      }
      [key: string]: any
    }
    [key: string]: any
  }
}

interface HotRankData {
  result: number
  pcursor: string
  webPageArea: string
  items: {
    type: string
    generated: boolean
    id: string
    typename: string
  }[]
}

export const kuaishouGetter: SourceGetter = async () => {
  // 获取快手首页HTML
  const response = await fetch('https://www.kuaishou.com/?isHome=1', {
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36',
    },
  })
  
  if (!response.ok) {
    throw new Error(`Failed to fetch Kuaishou: ${response.statusText}`)
  }

  const html = await response.text()
  
  // 提取window.__APOLLO_STATE__中的数据
  // 使用与 newsnow 相同的正则表达式：(\{.+?\});
  // 注意：.+? 是非贪婪匹配，会匹配到第一个 }; 为止
  const matches = html.match(/window\.__APOLLO_STATE__\s*=\s*(\{.+?\});/)
  
  if (!matches || !matches[1]) {
    throw new Error('无法获取快手热榜数据：未找到 __APOLLO_STATE__')
  }

  let data: KuaishouResponse
  try {
    // 解析JSON数据
    data = JSON.parse(matches[1]) as KuaishouResponse
  } catch (error) {
    throw new Error(`无法解析快手热榜数据：${error instanceof Error ? error.message : 'JSON解析失败'}`)
  }

  return parseKuaishouData(data)
}

function parseKuaishouData(data: KuaishouResponse): NewsItem[] {
  try {
    // 获取热榜数据ID - 使用与 newsnow 相同的键名格式
    const hotRankKey = 'visionHotRank({"page":"home"})'
    const hotRankQuery = data.defaultClient.ROOT_QUERY[hotRankKey]
    
    if (!hotRankQuery || !hotRankQuery.id) {
      throw new Error('无法找到热榜查询ID')
    }

    const hotRankId = hotRankQuery.id

    // 获取热榜列表数据
    const hotRankData = data.defaultClient[hotRankId] as HotRankData
    
    if (!hotRankData || !hotRankData.items || !Array.isArray(hotRankData.items)) {
      throw new Error('热榜数据格式不正确')
    }
    
    // 转换数据格式（与 newsnow 保持一致）
    return hotRankData.items
      .filter((k) => {
        const itemData = data.defaultClient[k.id]
        return itemData && itemData.tagType !== '置顶'
      })
      .map((item) => {
        // 从id中提取实际的热搜词
        const hotSearchWord = item.id.replace('VisionHotRankItem:', '')

        // 获取具体的热榜项数据
        const hotItem = data.defaultClient[item.id]
        
        if (!hotItem || !hotItem.name) {
          return null
        }

        return NewsItem.fromEntity({
          id: hotSearchWord,
          title: hotItem.name,
          url: `https://www.kuaishou.com/search/video?searchKey=${encodeURIComponent(hotItem.name)}`,
          extra: {
            icon: hotItem.iconUrl
              ? {
                  url: hotItem.iconUrl,
                  scale: 1,
                }
              : undefined,
          },
        })
      })
      .filter((item): item is NewsItem => item !== null)
  } catch (error) {
    throw new Error(`解析快手热榜数据失败：${error instanceof Error ? error.message : '未知错误'}`)
  }
}

