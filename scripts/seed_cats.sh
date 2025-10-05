#!/usr/bin/env bash
set -euo pipefail

APP_URL="${APP_URL:-http://app:8080}"

CAT_API="${CAT_API:-https://api.thecatapi.com/v1/breeds}"
LIMIT="${LIMIT:-25}"
API_KEY="${THECATAPI_KEY:-}"

echo "cats from TheCatAPI..."

HDR=()
[ -n "$API_KEY" ] && HDR+=(-H "x-api-key: $API_KEY")
JSON="$(curl -sS "${HDR[@]}" "$CAT_API")"

echo "Seeding up to $LIMIT cats into $APP_URL ..."

echo "$JSON" | jq -c ".[:$LIMIT][] | {name: .name, breed: .name}" | while read -r item; do
  years=$((RANDOM % 16))
  salary=$(( (RANDOM % 4000) + 1000 ))
  name=$(echo "$item"  | jq -r .name)
  breed=$(echo "$item" | jq -r .breed)

  body=$(jq -n --arg name "$name" --arg breed "$breed" --argjson years "$years" --argjson salary "$salary" \
    '{name:$name, years_experience:$years, breed:$breed, salary:$salary}')

  status=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$APP_URL/cats/create" \
    -H "Content-Type: application/json" \
    -d "$body")

  if [[ "$status" == "200" || "$status" == "201" ]]; then
    echo "  200! $name (exp=$years, salary=$salary)"
  else
    echo " 404! $name (HTTP $status)"
  fi
done

echo "Done."
