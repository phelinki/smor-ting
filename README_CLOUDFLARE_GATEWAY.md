## Cloudflare API Gateway and Edge Security for Smor-Ting

This repo includes a Cloudflare Worker that acts as an API gateway and security layer in front of the backend.

### Prerequisites

- Cloudflare account ID: 79b2d1445dc45e3d44d66dd40b6d8474
- Zone ID: efed3ad992752985706390ce9429b77e
- A Cloudflare API token with permissions:
  - Account: Workers Scripts: Edit
  - Zone: Workers Routes: Edit
  - Zone: Cache Purge: Purge
  - Zone: DNS: Edit (if you will also manage DNS)

### Create API Token

1. Go to Cloudflare Dashboard → My Profile → API Tokens → Create Token
2. Start with template: Edit Cloudflare Workers
3. Add permissions above, scope to your account and the target zone
4. Save the token securely

### Configure wrangler

Edit `cloudflare/wrangler.toml`:

```
account_id = "79b2d1445dc45e3d44d66dd40b6d8474"
[[routes]]
pattern = "api.<your-domain>.com/*"
zone_id = "efed3ad992752985706390ce9429b77e"
```

Set the upstream origin (Railway or your backend URL):

```bash
cd cloudflare
npx wrangler secret put UPSTREAM_ORIGIN
# Enter e.g. https://smor-ting-api.up.railway.app
```

Login and deploy:

```bash
npm i -D wrangler
npx wrangler login
npx wrangler deploy
```

### What the Worker does

- Forces HTTPS
- Adds security headers (HSTS, X-Content-Type-Options, X-Frame-Options, Referrer-Policy)
- Caches safe GET endpoints (excludes `/api/v1/auth/*` and `/api/v1/webhooks/*`) for 60s at the edge
- Forwards all other traffic to your origin

### Backend adjustments

- Rate limiter now respects `CF-Connecting-IP` for real client IPs behind Cloudflare.

### DNS

- Create `api.<your-domain>.com` proxied (orange cloud) A/CNAME record pointing to your origin (Railway domain) in the Cloudflare zone.

### Purging cache

To purge Workers cache programmatically you can call Cloudflare API or let TTL expire. For manual purge:

```bash
npx wrangler kv:key delete <key>
# or via Dashboard → Caching → Configuration → Purge Everything
```


