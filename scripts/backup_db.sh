#!/usr/bin/env sh
set -eu

BACKUP_DIR=${BACKUP_DIR:-/backups}
RETENTION_DAYS=${RETENTION_DAYS:-7}
DB_HOST=${DB_HOST:-db}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-homenavi_marketplace}
DB_USER=${DB_USER:-postgres}

TS=$(date +"%Y%m%d_%H%M%S")
OUT_FILE="$BACKUP_DIR/marketplace_${TS}.sql"

mkdir -p "$BACKUP_DIR"

PGPASSWORD="$DB_PASSWORD" pg_dump \
  --host "$DB_HOST" \
  --port "$DB_PORT" \
  --username "$DB_USER" \
  --format plain \
  --file "$OUT_FILE" \
  "$DB_NAME"

find "$BACKUP_DIR" -type f -name "marketplace_*.sql" -mtime +"$RETENTION_DAYS" -delete

echo "Backup written to $OUT_FILE"
