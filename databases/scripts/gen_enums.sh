#!/usr/bin/env bash
set -euo pipefail
export NO_COLOR=1 CLICOLOR=0 FORCE_COLOR=0 TERM=dumb

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
ATLAS_DIR="${REPO_ROOT}/databases"
GEN_DIR="${ATLAS_DIR}/gen/enums"
DEST_DIR="${REPO_ROOT}/enums"

# Source Go environment if available
if [ -f /volume1/use-menv.sh ]; then
  source /volume1/use-menv.sh >/dev/null 2>&1 || true
fi

# Ensure goimports is in PATH (it's installed in $HOME/go/bin)
# Try multiple possible HOME locations (root user might have different HOME)
for home_dir in "$HOME" "/home/LucasAsustor" "/root"; do
  if [ -d "$home_dir/go/bin" ] && [ -f "$home_dir/go/bin/goimports" ]; then
    export PATH="$PATH:$home_dir/go/bin"
    break
  fi
done

if [[ -z "${POSTGRES_USER}" ]]; then echo "Error: POSTGRES_USER environment variable is required"; exit 1; fi
if [[ -z "${POSTGRES_PASSWORD}" ]]; then echo "Error: POSTGRES_PASSWORD environment variable is required"; exit 1; fi
if [[ -z "${POSTGRES_HOST}" ]]; then echo "Error: POSTGRES_HOST environment variable is required"; exit 1; fi
if [[ -z "${POSTGRES_PORT}" ]]; then echo "Error: POSTGRES_PORT environment variable is required"; exit 1; fi
if [[ -z "${POSTGRES_DB_DEV}" ]]; then echo "Error: POSTGRES_DB_DEV environment variable is required"; exit 1; fi
if [[ -z "${POSTGRES_DB_PROD}" ]]; then echo "Error: POSTGRES_DB_PROD environment variable is required"; exit 1; fi
if [[ -z "${POSTGRES_DB_SHADOW}" ]]; then echo "Error: POSTGRES_DB_SHADOW environment variable is required"; exit 1; fi
if [[ -z "${ATLAS_ENV}" ]]; then echo "Error: ATLAS_ENV environment variable is required"; exit 1; fi
# POSTGRES_CONTAINER_NAME is optional, defaults to NFX-Stack-PostgreSQL
POSTGRES_CONTAINER_NAME="${POSTGRES_CONTAINER_NAME:-NFX-Stack-PostgreSQL}"

POSTGRES_USER="${POSTGRES_USER}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD}"
POSTGRES_HOST="${POSTGRES_HOST}"
POSTGRES_PORT="${POSTGRES_PORT}"
POSTGRES_DB_DEV="${POSTGRES_DB_DEV}"
POSTGRES_DB_PROD="${POSTGRES_DB_PROD}"
POSTGRES_DB_SHADOW="${POSTGRES_DB_SHADOW}"
ATLAS_ENV="${ATLAS_ENV}"

if [[ "${ATLAS_ENV}" == "prod" ]]; then
  POSTGRES_DB_DEV="${POSTGRES_DB_PROD}"
fi

# --- pre-check ---
# Note: gofmt and goimports are optional - script will continue if they're missing

# --- paths and temporary files ---
if [ -d "$GEN_DIR" ]; then
  rm -rf "${GEN_DIR:?}/"*
else
  mkdir -p "$GEN_DIR"
fi

# Clear all generated enum files before generating (avoid stale files when enums are removed)
if [ -d "$DEST_DIR" ]; then
  rm -f "${DEST_DIR}"/*_enum_dbgen.go
fi

ATLAS_ENV_ARGS=(
  -e "POSTGRES_USER=${POSTGRES_USER}"
  -e "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}"
  -e "POSTGRES_HOST=${POSTGRES_HOST}:${POSTGRES_PORT}"
  -e "POSTGRES_DB_DEV=${POSTGRES_DB_DEV}"
  -e "POSTGRES_DB_PROD=${POSTGRES_DB_PROD}"
  -e "POSTGRES_DB_SHADOW=${POSTGRES_DB_SHADOW}"
)

# -------------------------------
# Ensure shadow DB exists
# -------------------------------
docker exec "${POSTGRES_CONTAINER_NAME}" psql -U "${POSTGRES_USER}" -d postgres -c "CREATE DATABASE ${POSTGRES_DB_SHADOW};" >/dev/null 2>&1 || true

# -------------------------------
# Run Atlas inspect + generate
# Use local atlas command directly (same as pipeline script)
# -------------------------------
cd "${ATLAS_DIR}" || exit 1
if ! atlas schema inspect --env gen-enums; then
  echo "Error: Atlas schema inspect failed" >&2
  echo "Please check:" >&2
  echo "  1. Database connection settings (POSTGRES_PASSWORD, etc.)" >&2
  echo "  2. Database is running and accessible" >&2
  exit 1
fi

# -------------------------------
# Move generated files to enums directory
# Format: {schema}__{bucket}.go -> enums/{schema}_{bucket}_enum_dbgen.go
# -------------------------------
mkdir -p "${DEST_DIR}"

shopt -s nullglob
file_count=0
for src in "${GEN_DIR}"/*.go; do
  base="$(basename "${src}")"
  # Extract schema and bucket from filename: {schema}__{bucket}.go
  if [[ "$base" =~ ^([^_]+)__(.+)$ ]]; then
    schema="${BASH_REMATCH[1]}"
    rest="${BASH_REMATCH[2]}"
    bucket="${rest%.go}"
    dest_name="${schema}_${bucket}_enum_dbgen.go"
    
    # Clean old file for this specific enum
    rm -f "${DEST_DIR}/${dest_name}"
    
    mv "${src}" "${DEST_DIR}/${dest_name}"
    echo "Moved ${base} -> ${DEST_DIR}/${dest_name}"
    file_count=$((file_count + 1))
  else
    echo "Warning: Unexpected filename format: ${base}" >&2
  fi
done
shopt -u nullglob

if [ $file_count -eq 0 ]; then
  echo "Info: No enum files generated in ${GEN_DIR}" >&2
  echo "This might indicate that:"
  echo "  1. No enums are defined in the schema (this is OK)"
  echo "  2. Template generation failed"
  # Don't exit with error - it's OK to have no enums
fi

# -------------------------------
# Format generated Go files
# -------------------------------
set +e  # Temporarily disable exit on error for formatting
format_success=true
if command -v goimports >/dev/null 2>&1; then
  if ! goimports -w "${DEST_DIR}" 2>&1; then
    echo "Warning: goimports failed (non-fatal)" >&2
    format_success=false
  fi
else
  echo "Info: goimports not found, skipping formatting" >&2
fi
if command -v gofmt >/dev/null 2>&1; then
  if ! gofmt -s -w "${DEST_DIR}" 2>&1; then
    echo "Warning: gofmt failed (non-fatal)" >&2
    format_success=false
  fi
else
  echo "Info: gofmt not found, skipping formatting" >&2
fi
# Reset exit status to success before re-enabling set -e
true
set -e  # Re-enable exit on error

# Always exit successfully if we got here
echo "Enums generated successfully."
exit 0

