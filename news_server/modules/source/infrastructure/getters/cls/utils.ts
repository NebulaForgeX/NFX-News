/**
 * 财联社工具函数
 * 从 newsnow/server/sources/cls/utils.ts 迁移
 * https://github.com/DIYgod/RSSHub/blob/master/lib/routes/cls/utils.ts
 */
import md5 from 'md5'
import { createHash } from 'crypto'

const params = {
  appName: 'CailianpressWeb',
  os: 'web',
  sv: '7.7.5',
}

export async function getSearchParams(moreParams?: Record<string, string>) {
  const searchParams = new URLSearchParams({ ...params, ...moreParams })
  searchParams.sort()
  
  // 计算 SHA-1
  const sha1Hash = createHash('sha1').update(searchParams.toString()).digest('hex')
  
  // 计算 MD5
  const sign = md5(sha1Hash)
  
  searchParams.append('sign', sign)
  return searchParams
}

