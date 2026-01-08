/**
 * HTTP 服务器创建
 * 类似 Sjgz-Backend 的 interfaces/http/server.go
 */
import type { FastifyInstance } from 'fastify'
import fastify from 'fastify'
import cors from '@fastify/cors'
import { Router } from './router'
import { Registry } from './registry'
import { SourceHandler } from './handler/source'
import type { GetterService } from '../../application/getter/service'

export interface NewHTTPServerOptions {
  getterService: GetterService
}

export async function NewHTTPServer(options: NewHTTPServerOptions): Promise<FastifyInstance> {
  const app = fastify({
    logger: true,
  })

  // 注册 CORS 插件，允许跨域请求
  await app.register(cors, {
    origin: true, // 允许所有来源
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    allowedHeaders: ['Content-Type', 'Authorization'],
  })

  // 创建 Handler
  const sourceHandler = new SourceHandler(options.getterService)

  // 创建 Registry
  const registry = new Registry(sourceHandler)

  // 创建 Router 并注册路由
  const router = new Router(app, registry)
  router.RegisterRoutes()

  return app
}

