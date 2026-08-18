.PHONY: dev prod up down build prod-build restart logs ps clean \
        prod-down prod-logs \
        db-shell db-migrate db-rollback db-status backend-shell frontend-shell \
        backend-logs frontend-logs db-logs \
        test test-watch recette \
        web-build web-preview \
        mobile-host mobile-build mobile-sync mobile-android mobile-ios \
        sonar-up sonar-down sonar-logs sonar-scan \
        zap-baseline zap-api-scan zap-full-scan

# Fichiers Compose : base + override d'environnement.
COMPOSE := docker compose
DEV     := docker compose -f docker-compose.yml -f docker-compose.dev.yml
PROD    := docker compose -f docker-compose.yml -f docker-compose.prod.yml

# ── Environnements ──────────────────────────────────────────────────────────────
# dev  : hot-reload (air + Vite), bind mounts.
# prod : binaire Go compilé + front statique servi par nginx (rebuild à chaque up).
dev:
	$(DEV) up -d
	@echo "DEV → front http://localhost:5173  ·  api http://localhost:8080"

prod:
	$(PROD) up -d --build
	@echo "PROD → front http://localhost:8081  ·  api http://localhost:8080"

# Alias historique.
up: dev

# Arrête l'app (dev OU prod : `down` agit sur tout le projet).
down:
	$(DEV) down

prod-down:
	$(PROD) down

build:
	$(DEV) build

prod-build:
	$(PROD) build

restart: down dev

logs:
	$(DEV) logs -f

prod-logs:
	$(PROD) logs -f

ps:
	$(DEV) ps

clean:
	$(DEV) down -v --remove-orphans

# ── Selective logs ────────────────────────────────────────────────────────────
backend-logs:
	$(DEV) logs -f backend

frontend-logs:
	$(DEV) logs -f frontend

db-logs:
	$(DEV) logs -f db

# ── Shells ────────────────────────────────────────────────────────────────────
backend-shell:
	$(DEV) exec backend sh

frontend-shell:
	$(DEV) exec frontend sh

db-shell:
	$(DEV) exec db psql -U $${POSTGRES_USER:-cesizen} -d $${POSTGRES_DB:-cesizen}

# ── Tests ─────────────────────────────────────────────────────────────────────
# Toujours en dev : seule l'image dev embarque la toolchain Go.
test:
	$(DEV) exec backend go test ./... -v

# Recette fonctionnelle (API) — rejoue les cas du cahier de tests
recette:
	@bash scripts/recette.sh all

# Lancer un cas précis : make TC-INFO-01
TC-%:
	@bash scripts/recette.sh TC-$*

test-watch:
	$(DEV) exec backend sh -c 'find . -name "*.go" | entr -c go test ./... -v'

# ── Migrations ────────────────────────────────────────────────────────────────
# liquibase est dans la base commune (indépendant de l'environnement).
db-migrate:
	$(COMPOSE) run --rm liquibase \
		--url=jdbc:postgresql://db:5432/$${POSTGRES_DB:-cesizen} \
		--username=$${POSTGRES_USER:-cesizen} \
		--password=$${POSTGRES_PASSWORD:-cesizen} \
		--changeLogFile=changelog/db.changelog-master.xml \
		update

db-rollback:
	$(COMPOSE) run --rm liquibase \
		--url=jdbc:postgresql://db:5432/$${POSTGRES_DB:-cesizen} \
		--username=$${POSTGRES_USER:-cesizen} \
		--password=$${POSTGRES_PASSWORD:-cesizen} \
		--changeLogFile=changelog/db.changelog-master.xml \
		rollbackCount 1

db-status:
	$(COMPOSE) run --rm liquibase \
		--url=jdbc:postgresql://db:5432/$${POSTGRES_DB:-cesizen} \
		--username=$${POSTGRES_USER:-cesizen} \
		--password=$${POSTGRES_PASSWORD:-cesizen} \
		--changeLogFile=changelog/db.changelog-master.xml \
		status

# ── Web (build de production) ─────────────────────────────────────────────────
# Dev web : `make dev` (Vite sur :5173, hot-reload).
# Prod conteneurisée : `make prod`. Ci-dessous : build local hors Docker.

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

# ── Qualité de code (SonarQube) ───────────────────────────────────────────────
# Services derrière le profile "sonar" (ne démarrent pas avec `make dev/prod`).
# UI sur http://localhost:9000 (login initial admin/admin).
# Prérequis hôte Linux : sudo sysctl -w vm.max_map_count=262144
# NB : on cible les services par leur nom (up/rm) pour ne PAS toucher à l'app.

sonar-up:
	$(COMPOSE) --profile sonar up -d sonarqube
	@echo "SonarQube démarre sur http://localhost:9000 (patiente ~1 min au 1er lancement)."

sonar-down:
	$(COMPOSE) --profile sonar rm -sfv sonarqube sonar-db

sonar-logs:
	$(COMPOSE) --profile sonar logs -f sonarqube

# Analyse backend + frontend. Nécessite SONAR_TOKEN dans .env (créé dans l'UI).
sonar-scan:
	$(COMPOSE) --profile sonar-scan run --rm sonar-scanner

# ── Sécurité DAST (OWASP ZAP) ─────────────────────────────────────────────────
# Scanne l'app EN COURS D'EXÉCUTION → lancer `make dev` d'abord.
# Rapports écrits dans security/zap/. Cibles atteintes par nom de service Docker.
# ZAP_TARGET surchargeable : make zap-baseline ZAP_TARGET=http://backend:8080
ZAP_TARGET ?= http://frontend:5173

# Scan passif rapide (spider + règles passives). Le plus courant en CI.
zap-baseline:
	$(DEV) --profile zap run --rm zap \
		zap-baseline.py -t $(ZAP_TARGET) -r zap-baseline-report.html -I
	@echo "Rapport : security/zap/zap-baseline-report.html"

# Scan d'API REST à partir de l'OpenAPI (défaut : format openapi).
# make zap-api-scan ZAP_TARGET=http://backend:8080/openapi.json
zap-api-scan:
	$(DEV) --profile zap run --rm zap \
		zap-api-scan.py -t $(ZAP_TARGET) -f openapi -r zap-api-report.html -I
	@echo "Rapport : security/zap/zap-api-report.html"

# Scan actif complet (attaques réelles). Long. À réserver à un env de test.
zap-full-scan:
	$(DEV) --profile zap run --rm zap \
		zap-full-scan.py -t $(ZAP_TARGET) -r zap-full-report.html -I
	@echo "Rapport : security/zap/zap-full-report.html"
