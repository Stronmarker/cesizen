#!/usr/bin/env bash
# Source UNIQUE de l'IP du backend mobile : MOBILE_HOST dans le .env racine.
#
# Ce script en dérive automatiquement :
#   - frontend/.env.mobile          (VITE_API_URL pour le build Vite)
#   - .../network_security_config.xml (cleartext Android)
#
# => Pour changer l'IP, on ne touche QUE .env (MOBILE_HOST). Lancé par make mobile-build.
set -e

# Racine du projet (le script est dans scripts/)
cd "$(dirname "$0")/.."

ROOT_ENV=".env"
MOBILE_ENV="frontend/.env.mobile"
XML="frontend/android/app/src/main/res/xml/network_security_config.xml"
PORT="8080"

# Lire MOBILE_HOST depuis .env (fallback émulateur si absent)
host=""
if [ -f "$ROOT_ENV" ]; then
  host="$(grep -E '^[[:space:]]*MOBILE_HOST=' "$ROOT_ENV" | head -1 | cut -d= -f2- | tr -d '[:space:]\r')"
fi
if [ -z "$host" ]; then
  host="10.0.2.2"
  echo "ℹ️  MOBILE_HOST absent de $ROOT_ENV → défaut émulateur ($host)"
fi

url="http://${host}:${PORT}/api/v1"

# 1) Générer frontend/.env.mobile (consommé par vite build --mode mobile)
cat > "$MOBILE_ENV" <<EOF
# FICHIER GÉNÉRÉ — ne pas éditer à la main.
# Dérivé de MOBILE_HOST (.env racine) par scripts/sync-mobile-host.sh.
VITE_API_URL=$url
EOF

# 2) Générer la config cleartext Android (toujours 10.0.2.2 + localhost + l'hôte)
cat > "$XML" <<EOF
<?xml version="1.0" encoding="utf-8"?>
<!-- FICHIER GÉNÉRÉ — ne pas éditer à la main.
     Dérivé de MOBILE_HOST (.env racine) par scripts/sync-mobile-host.sh.
     Hôte de dev autorisé en cleartext : $host -->
<network-security-config>
    <domain-config cleartextTrafficPermitted="true">
        <domain includeSubdomains="false">10.0.2.2</domain>
        <domain includeSubdomains="false">localhost</domain>
        <domain includeSubdomains="false">$host</domain>
    </domain-config>
</network-security-config>
EOF

echo "✅ Mobile configuré depuis .env (MOBILE_HOST=$host) → $url"
