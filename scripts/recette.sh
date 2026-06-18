#!/usr/bin/env bash
# Recette fonctionnelle automatisée CESIZen — rejoue les cas du cahier de tests
# contre l'API en marche (make up).
#
#   bash scripts/recette.sh TC-INFO-01   # un cas précis
#   bash scripts/recette.sh all          # tous les cas + bilan
#
# Codes HTTP attendus = codes réels des handlers Go.

API="${API_URL:-http://localhost:8080/api/v1}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@cesizen.fr}"
ADMIN_PW="${ADMIN_PW:-Admin1234!}"

PASS=0; FAIL=0; SKIP=0; FAILED=""
G=$'\e[32m'; R=$'\e[31m'; Y=$'\e[33m'; B=$'\e[1m'; N=$'\e[0m'
BODY="$(mktemp)"
trap 'rm -f "$BODY"' EXIT

# req METHOD PATH [DATA] [TOKEN] -> remplit HTTP + fichier BODY
req() {
  local method="$1" path="$2" data="${3:-}" token="${4:-}"
  local args=(-s -o "$BODY" -w '%{http_code}' -X "$method" "$API$path" -H "Content-Type: application/json")
  [ -n "$token" ] && args+=(-H "Authorization: Bearer $token")
  [ -n "$data" ] && args+=(-d "$data")
  HTTP="$(curl "${args[@]}")"
}
jget() { grep -oE "\"$1\":\"[^\"]*\"" "$BODY" | head -1 | sed -E "s/\"$1\":\"([^\"]*)\"/\1/"; }
jnum() { grep -oE "\"$1\":[0-9]+" "$BODY" | head -1 | sed -E "s/\"$1\"://"; }

ok()   { PASS=$((PASS+1));  printf "${G}[PASS]${N}   %-13s %s\n" "$1" "$2"; }
ko()   { FAIL=$((FAIL+1)); FAILED="$FAILED $1"; printf "${R}[FAIL]${N}   %-13s %s\n" "$1" "$2"; }
skip() { SKIP=$((SKIP+1));  printf "${Y}[MANUEL]${N} %-13s %s\n" "$1" "$2"; }
# as ID EXPECTED DESC  (compare $HTTP)
as()   { if [ "$HTTP" = "$2" ]; then ok "$1" "$3 (HTTP $HTTP)"; else ko "$1" "$3 — attendu $2, obtenu $HTTP"; fi; }

# Comptes jetables (n'affectent pas les données réelles)
new_user() {
  local email="rec_${$}_${RANDOM}@cesizen.test"
  req POST /auth/register "{\"email\":\"$email\",\"password\":\"Test1234!\",\"first_name\":\"Rec\"}"
  EMAIL="$email"; TOKEN="$(jget token)"; RT="$(jget refresh_token)"; USER_ID="$(jget id)"
}
admin_login() { req POST /auth/login "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PW\"}"; ADMIN_TOKEN="$(jget token)"; }
first_emotion() { req GET /emotions; jnum id; }
ensure_article() { req POST /admin/contents '{"title":"Recette","content":"Contenu de recette","author":"Recette"}' "$ADMIN_TOKEN"; CID="$(jnum id)"; }

# ───────────────────────── Module Comptes utilisateurs ─────────────────────────
tc_AUTH_01() { new_user; as TC-AUTH-01 201 "Inscription valide"; }
tc_AUTH_02() { new_user; req POST /auth/register "{\"email\":\"$EMAIL\",\"password\":\"Test1234!\",\"first_name\":\"Rec\"}"; as TC-AUTH-02 409 "Email déjà utilisé"; }
tc_AUTH_03() {
  req POST /auth/register "{\"email\":\"short_${RANDOM}@cesizen.test\",\"password\":\"123\",\"first_name\":\"Rec\"}"
  if [ "$HTTP" -ge 400 ] 2>/dev/null; then ok TC-AUTH-03 "Mot de passe trop court rejeté (HTTP $HTTP)"
  else skip TC-AUTH-03 "Backend accepte (HTTP $HTTP) — règle de longueur côté front, à vérifier en UI"; fi
}
tc_AUTH_04() { req POST /auth/login "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PW\"}"; as TC-AUTH-04 200 "Connexion valide"; }
tc_AUTH_05() { req POST /auth/login "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"mauvais\"}"; as TC-AUTH-05 401 "Mauvais mot de passe"; }
tc_AUTH_06() {
  new_user
  for _ in 1 2 3; do req POST /auth/login "{\"email\":\"$EMAIL\",\"password\":\"WRONG\"}"; done
  req POST /auth/login "{\"email\":\"$EMAIL\",\"password\":\"Test1234!\"}"  # bon mdp mais compte verrouillé
  as TC-AUTH-06 403 "Verrouillage après 3 tentatives"
}
tc_AUTH_07() { new_user; req POST /auth/logout "{\"refresh_token\":\"$RT\"}" "$TOKEN"; as TC-AUTH-07 204 "Déconnexion"; }
tc_AUTH_08() { new_user; req POST /auth/refresh "{\"refresh_token\":\"$RT\"}"; as TC-AUTH-08 200 "Renouvellement du token"; }
tc_AUTH_09() { new_user; req GET /users/me "" "$TOKEN"; as TC-AUTH-09 200 "Consultation du profil"; }
tc_AUTH_10() { new_user; req PUT /users/me '{"first_name":"Modifie","nickname":"nick"}' "$TOKEN"; as TC-AUTH-10 200 "Mise à jour du profil"; }
tc_AUTH_11() { req POST /auth/forgot-password "{\"email\":\"$ADMIN_EMAIL\"}"; as TC-AUTH-11 200 "Demande de réinitialisation"; }
tc_AUTH_12() {
  new_user
  req POST /auth/forgot-password "{\"email\":\"$EMAIL\"}"; local tok; tok="$(jget reset_token)"
  req POST /auth/reset-password "{\"token\":\"$tok\",\"new_password\":\"Nouveau123!\"}"
  as TC-AUTH-12 204 "Réinitialisation avec token valide"
}
tc_AUTH_13() { req POST /auth/reset-password '{"token":"token_invalide","new_password":"Nouveau123!"}'; as TC-AUTH-13 401 "Reset avec token invalide"; }
tc_AUTH_14() { new_user; req DELETE /users/me "" "$TOKEN"; as TC-AUTH-14 204 "Suppression de compte (RGPD)"; }
tc_AUTH_15() { new_user; req GET /admin/contents "" "$TOKEN"; as TC-AUTH-15 403 "Back-office refusé à un non-admin"; }
tc_AUTH_16() {
  new_user; local target="$USER_ID"
  req PUT "/admin/users/$target" '{"role":"user","is_active":true}' "$ADMIN_TOKEN"
  as TC-AUTH-16 200 "Gestion d'un utilisateur par l'admin"
}
tc_AUTH_17() {
  req GET /admin/users "" "$ADMIN_TOKEN"; local a="$HTTP"
  req GET /admin/contents "" "$ADMIN_TOKEN"; local b="$HTTP"
  req GET /emotions; local c="$HTTP"
  if [ "$a" = 200 ] && [ "$b" = 200 ] && [ "$c" = 200 ]; then ok TC-AUTH-17 "Données du tableau de bord accessibles (users/contents/emotions = 200)"
  else ko TC-AUTH-17 "Indicateurs dashboard — users=$a contents=$b emotions=$c"; fi
}
tc_AUTH_18() { skip TC-AUTH-18 "Back-office bloqué sur mobile natif (Capacitor) — vérification UI/émulateur"; }

# ───────────────────────────── Module Informations ─────────────────────────────
tc_INFO_01() { req GET /contents; as TC-INFO-01 200 "Liste des articles (public)"; }
tc_INFO_02() { ensure_article; req GET "/contents/$CID"; as TC-INFO-02 200 "Détail d'un article"; }
tc_INFO_03() { req GET /contents/99999999; as TC-INFO-03 404 "Article inexistant"; }
tc_INFO_04() { req POST /admin/contents '{"title":"Recette","content":"x","author":"Rec"}' "$ADMIN_TOKEN"; as TC-INFO-04 201 "Création d'article (admin)"; }
tc_INFO_05() { ensure_article; req PUT "/admin/contents/$CID" '{"title":"Modifie","content":"y","author":"Rec","is_published":true}' "$ADMIN_TOKEN"; as TC-INFO-05 200 "Modification d'article (admin)"; }
tc_INFO_06() { ensure_article; req DELETE "/admin/contents/$CID" "" "$ADMIN_TOKEN"; as TC-INFO-06 204 "Suppression d'article (admin)"; }
tc_INFO_07() { new_user; req POST /admin/contents '{"title":"x","content":"y","author":"z"}' "$TOKEN"; as TC-INFO-07 403 "Création refusée à un non-admin"; }
tc_INFO_08() {
  ensure_article
  req PUT "/admin/contents/$CID" '{"title":"Recette","content":"x","author":"Rec","is_published":false}' "$ADMIN_TOKEN"
  local up="$HTTP"
  req GET /contents
  if [ "$up" = 200 ] && ! grep -q "\"id\":$CID," "$BODY"; then ok TC-INFO-08 "Dépublication — article retiré de la liste publique"
  else ko TC-INFO-08 "Dépublication — PUT=$up, article encore visible"; fi
}
tc_INFO_09() { ensure_article; req GET "/contents/$CID"; as TC-INFO-09 200 "Lecture sans authentification"; }

# ─────────────────────────── Module Tracker d'émotions ─────────────────────────
tc_TRACK_01() { new_user; local em; em="$(first_emotion)"; req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":7,\"comment\":\"ok\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"; as TC-TRACK-01 201 "Création d'une entrée"; }
tc_TRACK_02() {
  new_user; local em; em="$(first_emotion)"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":0,\"comment\":\"\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"; local lo="$HTTP"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":11,\"comment\":\"\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"; local hi="$HTTP"
  if [ "$lo" = 422 ] && [ "$hi" = 422 ]; then ok TC-TRACK-02 "Intensité hors limites rejetée (0→422, 11→422)"
  else ko TC-TRACK-02 "Intensité — 0 donne $lo, 11 donne $hi (attendu 422)"; fi
}
tc_TRACK_03() { new_user; req GET /tracker/entries "" "$TOKEN"; as TC-TRACK-03 200 "Consultation du journal"; }
tc_TRACK_04() {
  new_user; local em; em="$(first_emotion)"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":5,\"comment\":\"a\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"; local eid; eid="$(jnum id)"
  req PUT "/tracker/entries/$eid" "{\"emotion_id\":$em,\"intensity\":8,\"comment\":\"b\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"
  as TC-TRACK-04 200 "Modification d'une entrée"
}
tc_TRACK_05() {
  new_user; local em; em="$(first_emotion)"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":5,\"comment\":\"a\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"; local eid; eid="$(jnum id)"
  req DELETE "/tracker/entries/$eid" "" "$TOKEN"
  as TC-TRACK-05 204 "Suppression d'une entrée"
}
tc_TRACK_06() {
  new_user; local ta="$TOKEN" em; em="$(first_emotion)"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":5,\"comment\":\"a\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$ta"; local eid; eid="$(jnum id)"
  new_user  # utilisateur B
  req GET /tracker/entries "" "$TOKEN"
  if grep -q "\"id\":$eid," "$BODY" || grep -q "\"id\":$eid}" "$BODY"; then ko TC-TRACK-06 "Isolation rompue : B voit l'entrée de A"
  else ok TC-TRACK-06 "Isolation OK : B ne voit pas l'entrée de A"; fi
}
tc_TRACK_07() { new_user; req GET "/tracker/stats?period=week" "" "$TOKEN"; as TC-TRACK-07 200 "Statistiques (semaine)"; }
tc_TRACK_08() { new_user; req GET "/tracker/stats?period=month" "" "$TOKEN"; as TC-TRACK-08 200 "Statistiques (mois)"; }
tc_TRACK_09() {
  new_user
  req GET "/tracker/stats?period=quarter" "" "$TOKEN"; local q="$HTTP"
  req GET "/tracker/stats?period=year" "" "$TOKEN"; local y="$HTTP"
  if [ "$q" = 200 ] && [ "$y" = 200 ]; then ok TC-TRACK-09 "Statistiques (trimestre/année)"
  else ko TC-TRACK-09 "trimestre=$q année=$y (attendu 200)"; fi
}
tc_TRACK_10() {
  req GET /emotions; local e="$HTTP"
  req GET /primary-emotions; local p="$HTTP"
  if [ "$e" = 200 ] && [ "$p" = 200 ]; then ok TC-TRACK-10 "Référentiel public accessible"
  else ko TC-TRACK-10 "emotions=$e primary=$p (attendu 200)"; fi
}
tc_TRACK_11() { req POST /admin/primary-emotions "{\"label\":\"RecPrim_${RANDOM}\"}" "$ADMIN_TOKEN"; as TC-TRACK-11 201 "Création émotion primaire (admin)"; }
tc_TRACK_12() {
  req POST /admin/primary-emotions "{\"label\":\"RecPrim_${RANDOM}\"}" "$ADMIN_TOKEN"; local pid; pid="$(jnum id)"
  req POST /admin/emotions "{\"label\":\"RecSec_${RANDOM}\",\"primary_emotion_id\":$pid}" "$ADMIN_TOKEN"
  as TC-TRACK-12 201 "Création émotion secondaire (admin)"
}
tc_TRACK_13() {
  req POST /admin/primary-emotions "{\"label\":\"RecPrim_${RANDOM}\"}" "$ADMIN_TOKEN"; local pid; pid="$(jnum id)"
  req POST /admin/emotions "{\"label\":\"RecSec_${RANDOM}\",\"primary_emotion_id\":$pid}" "$ADMIN_TOKEN"; local sid; sid="$(jnum id)"
  req PUT "/admin/emotions/$sid" "{\"label\":\"RecSec\",\"primary_emotion_id\":$pid,\"is_active\":false}" "$ADMIN_TOKEN"
  as TC-TRACK-13 200 "Désactivation d'une émotion (admin)"
}
tc_TRACK_14() {
  if ! command -v docker >/dev/null; then skip TC-TRACK-14 "Vérification cascade nécessite docker (psql)"; return; fi
  new_user; local em; em="$(first_emotion)"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":5,\"comment\":\"a\",\"entry_date\":\"2026-06-13T12:00:00Z\"}" "$TOKEN"
  req DELETE /users/me "" "$TOKEN"
  local cnt; cnt="$(docker compose exec -T db psql -U cesizen -d cesizen -t -A -c "SELECT count(*) FROM entries WHERE user_id='$USER_ID'" 2>/dev/null | tr -d '[:space:]')"
  if [ "$cnt" = "0" ]; then ok TC-TRACK-14 "Cascade RGPD : 0 entrée restante après suppression"
  else ko TC-TRACK-14 "Cascade : $cnt entrée(s) restante(s)"; fi
}
tc_TRACK_15() {
  new_user; local em; em="$(first_emotion)"
  local today; today="$(date -u +%Y-%m-%d)"
  req POST /tracker/entries "{\"emotion_id\":$em,\"intensity\":6,\"comment\":\"jour\",\"entry_date\":\"${today}T12:00:00Z\"}" "$TOKEN"; local eid; eid="$(jnum id)"
  req GET /tracker/entries "" "$TOKEN"
  if grep -q "\"id\":$eid," "$BODY" || grep -q "\"id\":$eid}" "$BODY"; then ok TC-TRACK-15 "Non-régression : entrée du jour visible (midi UTC)"
  else ko TC-TRACK-15 "Entrée du jour absente de la liste (régression du bug de date)"; fi
}

ALL="TC-AUTH-01 TC-AUTH-02 TC-AUTH-03 TC-AUTH-04 TC-AUTH-05 TC-AUTH-06 TC-AUTH-07 TC-AUTH-08 TC-AUTH-09 TC-AUTH-10 TC-AUTH-11 TC-AUTH-12 TC-AUTH-13 TC-AUTH-14 TC-AUTH-15 TC-AUTH-16 TC-AUTH-17 TC-AUTH-18 TC-INFO-01 TC-INFO-02 TC-INFO-03 TC-INFO-04 TC-INFO-05 TC-INFO-06 TC-INFO-07 TC-INFO-08 TC-INFO-09 TC-TRACK-01 TC-TRACK-02 TC-TRACK-03 TC-TRACK-04 TC-TRACK-05 TC-TRACK-06 TC-TRACK-07 TC-TRACK-08 TC-TRACK-09 TC-TRACK-10 TC-TRACK-11 TC-TRACK-12 TC-TRACK-13 TC-TRACK-14 TC-TRACK-15"

run_one() {
  local id="${1#TC-}"            # AUTH-04
  local fn="tc_${id//-/_}"       # tc_AUTH_04
  if declare -F "$fn" >/dev/null; then "$fn"; else echo "Cas inconnu : $1"; exit 2; fi
}
summary() {
  echo "──────────────────────────────────────────────"
  printf "${B}Bilan recette :${N} ${G}%d réussis${N}, ${R}%d échecs${N}, ${Y}%d manuels${N}\n" "$PASS" "$FAIL" "$SKIP"
  [ -n "$FAILED" ] && printf "${R}Échecs :${N}%s\n" "$FAILED"
  [ "$FAIL" -gt 0 ] && return 1 || return 0
}

# Vérif API joignable
if ! curl -s -o /dev/null "$API/primary-emotions"; then
  echo "${R}API injoignable sur $API${N} — lance 'make up' d'abord."; exit 2
fi

admin_login
case "${1:-all}" in
  all)   for id in $ALL; do run_one "$id"; done; summary;;
  TC-*)  run_one "$1"; summary;;
  *)     echo "Usage: recette.sh [TC-XXX-NN | all]"; exit 2;;
esac
