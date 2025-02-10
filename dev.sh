export GODEBUG=smtp=2
python scripts/update_version.py
APP_ENV=development go run .
