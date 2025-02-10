set -o allexport
source .env
set +o allexport


python scripts/update_version.py
APP_ENV=development go run .
