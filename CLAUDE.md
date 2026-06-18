# CESIZen

Mental health web application (CESI CDA formation project). Helps users track emotions, access mental health information, and manage stress.

## Stack

| Layer | Technology | Notes |
|---|---|---|
| Frontend | React 18 + Vite + Tailwind v4 | Port 5173. Mobile via Capacitor (Android/iOS), même base de code |
| Backend | Go | Port 8080, REST API |
| Database | PostgreSQL 16 | Port 5432 |
| Migrations | Liquibase | XML changelogs in `db/changelog/migrations/` |
| Infrastructure | Docker Compose | All services containerized |
| Build | Makefile | See commands below |

## Project Structure

```
cesizen/
├── CLAUDE.md
├── Makefile
├── docker-compose.yml
├── .env.example
├── backend/
│   ├── cmd/api/          # main.go entrypoint
│   └── internal/
│       ├── config/       # env/config loading
│       ├── db/           # DB connection pool
│       ├── domain/       # structs/models (pure Go, no DB deps)
│       ├── handler/      # HTTP handlers (thin layer, call services)
│       ├── middleware/    # auth, rate-limit, CORS, logging
│       ├── repository/   # SQL queries (pgx/sqlx)
│       ├── router/       # route registration
│       └── service/      # business logic
├── db/
│   └── changelog/
│       ├── db.changelog-master.xml
│       └── migrations/   # individual SQL changesets
└── frontend/
    └── src/
        ├── api/          # API client/calls
        ├── components/   # reusable UI components
        ├── contexts/     # React contexts (auth, etc.)
        ├── hooks/        # custom hooks
        ├── pages/
        │   ├── admin/    # back-office pages
        │   ├── info/     # public information pages
        │   └── tracker/  # emotion tracker pages
        └── types/        # TypeScript types
```

## Docker / Makefile Commands

```bash
make up             # start all services (detached)
make down           # stop all services
make build          # rebuild images
make restart        # down + up
make logs           # tail all logs
make ps             # service status
make clean          # down + remove volumes

make backend-logs   # tail backend only
make frontend-logs  # tail frontend only
make db-logs        # tail postgres only

make backend-shell  # sh into backend container
make frontend-shell # sh into frontend container
make db-shell       # psql into postgres

make db-migrate     # run pending Liquibase migrations
make db-rollback    # rollback last changeset
make db-status      # show migration status

# Tests
make test           # Go unit/integration tests (in backend container)
make recette        # functional end-to-end recette (scripts/recette.sh, all TC cases)
make TC-INFO-01     # run a single test case by ID (e.g. TC-AUTH-04)

# Web build
make web-build      # production web build (dist/), reads frontend/.env.production
make web-preview    # preview the production build locally

# Mobile (Capacitor)
make mobile-host    # regenerate frontend/.env.mobile + Android cleartext config from .env MOBILE_HOST
make mobile-build   # build:mobile (runs mobile-host first)
make mobile-sync    # mobile-build + npx cap sync
make mobile-android # mobile-sync + open Android Studio
make mobile-ios     # mobile-sync + open Xcode
```

Environment variables are in `.env` (copy from `.env.example`).

## Actors & Roles

| Role | Access |
|---|---|
| **Visiteur anonyme** | Public pages, information content, register |
| **Utilisateur connecté** | Emotion tracker (journal + stats), account settings |
| **Administrateur** | Back-office: user management, content editing, emotion referential config |

## Modules

**Obligatoires:**
- **Comptes utilisateurs** — register, login, logout, profile, password reset, JWT auth
- **Informations** — public content pages (CRUD by admin, read by all)

**Module au choix (chosen: Tracker d'émotions):**
- **Tracker d'émotions** — emotion journal (CRUD entries), stats by period, admin manages emotion referential

## Database Schema

Based on MCD/MLD from the spec. Migrations live in `db/changelog/migrations/`.

### Tables

**users**
- `id` UUID PK
- `email` varchar(255) UNIQUE NOT NULL
- `password_hash` varchar(255) NOT NULL
- `first_name` varchar(255)
- `nickname` varchar(255)
- `is_active` bool DEFAULT true
- `role` varchar(50) DEFAULT 'user' — `user` | `admin`
- `login_attempts` integer DEFAULT 0 — failed login counter
- `locked_until` timestamptz — temporary lock after 3 failed attempts
- `created_at` timestamptz
- `updated_at` timestamptz

> A default admin (`admin@cesizen.fr` / `Admin1234!`) is seeded automatically by a Liquibase changeset (idempotent). NULL `first_name`/`nickname` are handled via `COALESCE` in repository queries.

**contents** (information pages)
- `id` serial PK
- `title` text NOT NULL
- `content` text NOT NULL
- `author` text
- `is_published` bool DEFAULT true
- `created_at` timestamptz
- `updated_at` timestamptz

**primary_emotions** (emotion referential level 1)
- `id` serial PK
- `label` varchar(255) NOT NULL UNIQUE
- `is_active` bool DEFAULT true

**emotions** (emotion referential level 2)
- `id` serial PK
- `label` varchar(255) NOT NULL
- `primary_emotion_id` integer FK → primary_emotions.id
- `is_active` bool DEFAULT true

**entries** (tracker journal)
- `id` serial PK
- `user_id` UUID FK → users.id ON DELETE CASCADE
- `emotion_id` integer FK → emotions.id
- `intensity` integer CHECK (1..10)
- `comment` text
- `entry_date` timestamptz NOT NULL
- `created_at` timestamptz

**refresh_tokens**
- `token` text PK
- `user_id` UUID FK → users.id ON DELETE CASCADE
- `expires_at` timestamptz NOT NULL (7 days)
- `created_at` timestamptz

**password_reset_tokens**
- `token` text PK
- `user_id` UUID FK → users.id ON DELETE CASCADE
- `expires_at` timestamptz NOT NULL (1 hour)
- `created_at` timestamptz

### Emotion Referential (seed data)

Primary emotions: Joie, Colère, Peur, Tristesse, Surprise, Dégoût

Examples of secondary emotions per primary:
- Joie → Fierté, Contentement, Enchantement, Excitation, Émerveillement, Gratitude
- Colère → Frustration, Irritation, Rage, Ressentiment, Agacement, Hostilité
- Peur → Inquiétude, Anxiété, Terreur, Appréhension, Panique, Crainte
- Tristesse → Chagrin, Mélancolie, Abattement, Désespoir, Solitude, Dépression
- Surprise → Étonnement, Stupéfaction, Sidération, Incrédule, Émerveillement, Confusion
- Dégoût → Répulsion, Déplaisir, Nausée, Dédain, Horreur, Dégoût profond

## API Design

Base path: `/api/v1`

Authentication: JWT Bearer tokens. Access token (short-lived) + refresh token.
Max 3 failed login attempts → temporary account lock.

Key endpoint groups:
- `POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`, `POST /auth/logout`
- `GET/PUT /users/me` — own profile
- `GET/PUT/POST /admin/users` — admin user management
- `GET /contents`, `GET /contents/:id` — public
- `POST/PUT/DELETE /admin/contents/:id` — admin
- `GET /emotions`, `GET /primary-emotions` — public referential
- `POST/PUT/DELETE /admin/emotions` — admin manages referential
- `GET/POST/PUT/DELETE /tracker/entries` — authenticated user's journal
- `GET /tracker/stats?period=week|month|quarter|year` — authenticated

## Backend Conventions (Go)

- MVC-like layers: handler → service → repository
- Handlers are thin: parse request, call service, write response
- Services hold all business logic
- Repositories hold all SQL (use `pgx/v5` or `sqlx`)
- Domain structs in `internal/domain/` have no framework dependencies
- Return typed errors from services; handlers map them to HTTP codes
- All routes except public ones go through `middleware/auth.go` (401 if unauthenticated)
- Admin routes additionally go through `middleware/RequireRole("admin")` → **403** for an authenticated non-admin (not 401)
- Access token 15 min, refresh token 7 days (rotated on use, stored in `refresh_tokens`)
- Passwords hashed with bcrypt (cost 10)
- JWT secret from env `JWT_SECRET`

## Security Requirements

- HTTPS in production (reverse proxy)
- Passwords: bcrypt hashed, never stored plaintext
- JWT: short-lived access token + refresh token rotation
- Rate limiting on `/auth/login` (max 3 attempts before lock)
- CORS configured for allowed origins only
- RGPD: no personal data exported outside EU, user can delete account
- Input validation on all endpoints
- SQL injections prevented via parameterized queries (never string concat)

## Frontend Notes

Stack: **React 18 + Vite + Tailwind v4 + react-router-dom**, packaged for mobile with **Capacitor**.

- Build context `./frontend`, target `dev`, port 5173, env `VITE_API_URL`
- Two areas: **Front-Office** (public + user, `AppShell` with bottom tab nav) and **Back-Office** (admin, `AdminHeader` with tab nav + `/admin` dashboard)
- API clients in `src/api/`, auth state in `src/contexts/AuthContext.tsx` (token + refresh in localStorage)
- Pages are flat in `src/pages/` (admin pages under `src/pages/admin/`)

### Design system (`src/index.css`, Tailwind `@theme`)

Strictly **3 colors** + transparency (frosted glass):
- Vert forêt `green-brand #1a5c32` (primary) · Jaune beurre `yellow-brand #f2e2a0` (accent) · Crème `cream #f7f5f0` (bg)
- Derived: `green-deep/mid/light`, `yellow-deep/dark`
- Secondary tones use forest-green transparency (`text-green-brand/45`…), not greys
- Glass surfaces: `bg-white/55 backdrop-blur-md border-white/60`; soft pastel tints per primary emotion in the tracker

## Mobile (Capacitor)

- `androidScheme: 'https'` + **CapacitorHttp enabled** (native HTTP) so `fetch()` reaches an HTTP dev backend without mixed-content/CORS issues
- **Back-office is web-only**: `AdminRoute` redirects to `/` when `Capacitor.isNativePlatform()` (real authz still enforced server-side)
- **Single source of truth for the dev backend IP = `MOBILE_HOST` in root `.env`.** `make mobile-build/sync` runs `scripts/sync-mobile-host.sh`, which **generates** `frontend/.env.mobile` (`VITE_API_URL`) and `android/.../network_security_config.xml` (cleartext) from it — do not edit those by hand
- Build: `make mobile-sync` then rebuild in Android Studio / Xcode

## Testing

- `make test` — Go unit/integration tests (`internal/{service,handler,middleware}/*_test.go`)
- `make recette` / `make TC-XXX` — functional end-to-end recette (`scripts/recette.sh`) replaying the cahier de tests against the live API, asserting exact HTTP codes (uses throwaway accounts)

## Liquibase Conventions

- One changeset per logical change
- changeset IDs: `YYYY-MM-DD-NNN-description` (e.g. `2024-01-01-001-create-users`)
- Author: your name
- Master file: `db/changelog/db.changelog-master.xml` includes all migration files
- Seed data in dedicated changesets: emotion referential, and the **default admin** (idempotent via a `preConditions onFail="MARK_RAN"` check)
- Tables include `refresh_tokens` and `password_reset_tokens` migrations alongside the core ones
