# syntax=docker/dockerfile:1
#
# Builds the Next.js static export and bakes it into a caddy:2 image. The
# Caddyfile itself is bind-mounted by docker-compose (not copied here) so it
# can be edited without a rebuild.

FROM node:22-alpine AS build
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
# Empty base -> lib/api.ts calls relative /api/v1/... paths, so Caddy proxies
# them same-origin in prod; no browser CORS needed.
ENV NEXT_PUBLIC_API_BASE=""
RUN npm run build

FROM caddy:2 AS final
COPY --from=build /web/out /srv
