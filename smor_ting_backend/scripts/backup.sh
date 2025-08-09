#!/bin/bash

# Automated Backup Script for Smor-Ting (MongoDB Atlas or self-hosted)
# Creates compressed, timestamped backups and prunes old backups

set -euo pipefail

# Defaults (can be overridden by env vars)
BACKUP_DIR=${BACKUP_DIR:-"/var/backups/smor-ting"}
RETENTION_DAYS=${RETENTION_DAYS:-7}
MONGODB_URI=${MONGODB_URI:-""}
DB_NAME=${DB_NAME:-"smor_ting"}
S3_BUCKET=${S3_BUCKET:-""} # optional: s3://bucket/path

timestamp() { date +"%Y%m%d-%H%M%S"; }
log() { echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*"; }
err() { echo "[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $*" >&2; }

ensure_tools() {
  for bin in mongodump tar; do
    if ! command -v "$bin" >/dev/null 2>&1; then
      err "Required tool not found: $bin"
      exit 1
    fi
  done
}

prepare_dirs() {
  mkdir -p "$BACKUP_DIR"
}

perform_backup() {
  local ts=$(timestamp)
  local dump_dir="$BACKUP_DIR/dump-$ts"
  local archive="$BACKUP_DIR/mongo-$DB_NAME-$ts.tar.gz"

  log "Starting mongodump..."
  if [ -n "$MONGODB_URI" ]; then
    mongodump --uri="$MONGODB_URI" --db "$DB_NAME" --out "$dump_dir"
  else
    mongodump --db "$DB_NAME" --out "$dump_dir"
  fi

  log "Compressing backup..."
  tar -czf "$archive" -C "$dump_dir" .
  rm -rf "$dump_dir"
  log "Backup created at $archive"

  if [ -n "$S3_BUCKET" ]; then
    if command -v aws >/dev/null 2>&1; then
      log "Uploading to S3: $S3_BUCKET"
      aws s3 cp "$archive" "$S3_BUCKET/"
    else
      err "aws CLI not found; skipping S3 upload"
    fi
  fi
}

prune_old() {
  log "Pruning backups older than $RETENTION_DAYS days in $BACKUP_DIR"
  find "$BACKUP_DIR" -type f -name "mongo-$DB_NAME-*.tar.gz" -mtime +"$RETENTION_DAYS" -delete || true
}

case "${1:-run}" in
  run)
    ensure_tools
    prepare_dirs
    perform_backup
    prune_old
    ;;
  test)
    # dry run: verify we can create directories and list tools
    ensure_tools
    prepare_dirs
    log "Dry run complete (no dump executed)"
    ;;
  *)
    echo "Usage: $0 [run|test]";
    exit 2
    ;;
esac


