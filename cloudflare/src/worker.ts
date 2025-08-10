export default {
  async fetch(request: Request, env: any): Promise<Response> {
    const url = new URL(request.url)
    if (request.headers.get('x-forwarded-proto') !== 'https') {
      url.protocol = 'https:'
      return Response.redirect(url.toString(), 301)
    }

    const isAuth = url.pathname.startsWith('/api/v1/auth')
    const isWebhooks = url.pathname.startsWith('/api/v1/webhooks')

    const upstreamBase = env.UPSTREAM_ORIGIN as string
    if (!upstreamBase) {
      return new Response('Missing UPSTREAM_ORIGIN', { status: 500 })
    }
    const upstream = new URL(upstreamBase)
    upstream.pathname = url.pathname
    upstream.search = url.search

    // Rebuild headers without forwarding Host to avoid upstream mismatch
    const forwardHeaders = new Headers(request.headers)
    forwardHeaders.delete('host')
    const init: RequestInit = {
      method: request.method,
      headers: forwardHeaders,
    }
    if (!['GET', 'HEAD'].includes(request.method)) {
      init.body = await request.clone().arrayBuffer()
    }
    const req = new Request(upstream.toString(), init)

    const cache = caches.default
    const cacheKey = new Request(request.url, request)
    if (!isAuth && !isWebhooks && request.method === 'GET') {
      const hit = await cache.match(cacheKey)
      if (hit) return hit
    }

    const response = await fetch(req, {
      cf: {
        cacheTtl: (!isAuth && !isWebhooks && request.method === 'GET') ? 60 : 0,
        cacheEverything: (!isAuth && !isWebhooks && request.method === 'GET'),
      }
    } as any)

    const headers = new Headers(response.headers)
    headers.set('Strict-Transport-Security', 'max-age=31536000; includeSubDomains; preload')
    headers.set('X-Content-Type-Options', 'nosniff')
    headers.set('X-Frame-Options', 'DENY')
    headers.set('Referrer-Policy', 'no-referrer')
    headers.set('Server', 'edge')

    if (request.method === 'HEAD') {
      return new Response(null, { status: response.status, headers })
    }
    const copied = new Response(response.body, { status: response.status, headers })

    if (!isAuth && !isWebhooks && request.method === 'GET' && response.ok) {
      await cache.put(cacheKey, copied.clone())
    }

    return copied
  }
} satisfies ExportedHandler


