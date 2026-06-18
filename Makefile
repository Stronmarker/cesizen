.PHONY: up down build restart logs ps clean \
        db-shell db-migrate backend-shell frontend-shell \
        backend-logs frontend-logs db-logs \
        test test-watch recette \
        web-build web-preview \
        mobile-host mobile-build mobile-sync mobile-android mobile-ios

# ── Infrastructure ────────────────────────────────────────────────────────────
up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

restart: down up

logs:
	docker compose logs -f

ps:
	docker compose ps

clean:
	docker compose down -v --remove-orphans

# ── Selective logs ────────────────────────────────────────────────────────────
backend-logs:
	docker compose logs -f backend

frontend-logs:
	docker compose logs -f frontend

db-logs:
	docker compose logs -f db

# ── Shells ────────────────────────────────────────────────────────────────────
backend-shell:
	docker compose exec backend sh

frontend-shell:
	docker compose exec frontend sh

db-shell:
	docker compose exec db psql -U $${POSTGRES_USER:-cesizen} -d $${POSTGRES_DB:-cesizen}

# ── Tests ─────────────────────────────────────────────────────────────────────
test:
	docker compose exec backend go test ./... -v

# Recette fonctionnelle (API) — rejoue les cas du cahier de tests
recette:
	@bash scripts/recette.sh all

# Lancer un cas précis : make TC-INFO-01
TC-%:
	@bash scripts/recette.sh TC-$*

test-watch:
	docker compose exec backend sh -c 'find . -name "*.go" | entr -c go test ./... -v'

# ── Migrations ────────────────────────────────────────────────────────────────
db-migrate:
	docker compose run --rm liquibase \
		--url=jdbc:postgresql://db:5432/$${POSTGRES_DB:-cesizen} \
		--username=$${POSTGRES_USER:-cesizen} \
		--password=$${POSTGRES_PASSWORD:-cesizen} \
		--changeLogFile=changelog/db.changelog-master.xml \
		update

db-rollback:
	docker compose run --rm liquibase \
		--url=jdbc:postgresql://db:5432/$${POSTGRES_DB:-cesizen} \
		--username=$${POSTGRES_USER:-cesizen} \
		--password=$${POSTGRES_PASSWORD:-cesizen} \
		--changeLogFile=changelog/db.changelog-master.xml \
		rollbackCount 1

db-status:
	docker compose run --rm liquibase \
		--url=jdbc:postgresql://db:5432/$${POSTGRES_DB:-cesizen} \
		--username=$${POSTGRES_USER:-cesizen} \
		--password=$${POSTGRES_PASSWORD:-cesizen} \
		--changeLogFile=changelog/db.changelog-master.xml \
		status

# ── Web (build de production) ─────────────────────────────────────────────────
# Dev web : `make up` (Vite sur :5173, hot-reload).
# Prod : copier .env.production.example → .env.production et ajuster VITE_API_URL.

web-build:
	cd frontend && npm run build

web-preview: web-build
	cd frontend && npm run preview

# ── Mobile (Capacitor) ────────────────────────────────────────────────────────
# Source unique de l'IP du backend : VITE_API_URL dans frontend/.env.mobile.
# La config cleartext Android en est dérivée automatiquement (mobile-host).

# Régénère network_security_config.xml depuis frontend/.env.mobile
mobile-host:
	@bash scripts/sync-mobile-host.sh

mobile-build: mobile-host
	cd frontend && npm run build:mobile

mobile-sync: mobile-build
	cd frontend && npx cap sync

mobile-android: mobile-sync
	cd frontend && npx cap open android

mobile-ios: mobile-sync
	cd frontend && npx cap open ios
