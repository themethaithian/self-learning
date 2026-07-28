# Self-Improve Web — frontend

Next.js (App Router, TypeScript, Tailwind) static export, served by Caddy on
the VPS. All data-fetching is client-side at runtime — the Go API is the
only backend.

- Fonts: Inter (UI/English), Noto Sans Thai Looped (Thai lesson text),
  JetBrains Mono (code) — wired via `next/font/google` in `app/layout.tsx`.
- Talks to the Go API at `NEXT_PUBLIC_API_BASE` (default
  `http://localhost:8080`) using a bearer token stored in `localStorage`;
  see `lib/api.ts` and `lib/token.ts`.
- Package manager: npm only.

## Develop

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000). Paste a bearer token
(matching the Go API's `API_BEARER_TOKEN`) on the token entry page.

The root `docker-compose.yml`'s `web` service also publishes port 3000 — run
`make dev-down` first if that stack is up, or `next dev` fails to bind.

## Build

```bash
npm run build
```

Produces the static export in `out/`.

### Prod build (same-origin through Caddy)

`web/Dockerfile` builds with `NEXT_PUBLIC_API_BASE=""`, so `lib/api.ts` calls
relative `/api/v1/...` paths and Caddy reverse-proxies them to the Go API —
no CORS involved. To reproduce that build manually:

```bash
NEXT_PUBLIC_API_BASE= npm run build
```
