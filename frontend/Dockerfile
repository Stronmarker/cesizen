# ── Dev : Vite dev server (hot-reload) ─────────────────────────────────────
FROM node:20-alpine AS dev
WORKDIR /app
EXPOSE 5173
CMD ["sh", "-c", "npm install && npm run dev -- --host"]

# ── Builder : build de production (dist/) ──────────────────────────────────
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# VITE_API_URL est injecté au build (Vite fige les VITE_* dans le bundle).
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# ── Prod : nginx sert le build statique (SPA) ──────────────────────────────
FROM nginx:alpine AS prod
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
