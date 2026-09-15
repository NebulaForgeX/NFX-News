#!/usr/bin/env bash
set -euo pipefail
export NO_COLOR=1 CLICOLOR=0 FORCE_COLOR=0 TERM=dumb

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
ATLAS_DIR="${REPO_ROOT}/databases"
GEN_DIR="${ATLAS_DIR}/gen/models"

# Source Go environment if available
if [ -f /volume1/use-menv.sh ]; then
  source /volume1/use-menv.sh >/dev/null 2>&1 || true
fi

# Ensure goimports is in PATH (it's installed in $HOME/go/bin)
# This must be done before the 'need' check below
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
if [[ -z "${RESOURCES_DOCKER_COMPOSE}" ]]; then echo "Error: RESOURCES_DOCKER_COMPOSE environment variable is required"; exit 1; fi

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
need() { command -v "$1" >/dev/null 2>&1 || { echo "Missing dependency: $1" >&2; exit 1; }; }

# Double-check PATH for goimports (in case source didn't work)
# Try multiple possible HOME locations
if ! command -v goimports >/dev/null 2>&1; then
  for home_dir in "$HOME" "/home/LucasAsustor" "/root"; do
    if [ -f "$home_dir/go/bin/goimports" ]; then
      export PATH="$PATH:$home_dir/go/bin"
      break
    fi
  done
fi

need goimports
need gofmt

# --- paths and temporary files ---
if [ -d "$GEN_DIR" ]; then
  rm -rf "${GEN_DIR:?}/"*
else
  mkdir -p "$GEN_DIR"
fi

# Clear all generated model files in module dirs before generating (avoid stale files when tables are removed)
for module_models in "${REPO_ROOT}/modules"/*/infrastructure/rdb/models; do
  if [ -d "$module_models" ]; then
    rm -f "${module_models}"/*_dbgen.go
  fi
done

# set local module prefix, for goimports grouping
MODPATH="$(go list -m 2>/dev/null || echo "")"
[[ -n "$MODPATH" ]] && export GOIMPORTSLOCAL="$MODPATH"

# Use direct docker exec with container name (same as pipeline script)
# Get container name from environment or use default
POSTGRES_CONTAINER_NAME="${POSTGRES_CONTAINER_NAME:-NFX-Stack-PostgreSQL}"

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
if ! atlas schema inspect --env gen-models; then
  echo "Error: Atlas schema inspect failed" >&2
  echo "Please check:" >&2
  echo "  1. Database connection settings (POSTGRES_PASSWORD, etc.)" >&2
  echo "  2. Database is running and accessible" >&2
  exit 1
fi

# -------------------------------
# Distribute generated files to module directories based on schema name
# Format: {schema}__{table}.go -> modules/{schema}/infrastructure/rdb/models/{table}_dbgen.go
# -------------------------------
shopt -s nullglob
file_count=0
for src in "${GEN_DIR}"/*.go; do
  base="$(basename "${src}")"
  
  # Extract schema and table from filename: {schema}__{table}.go
  if [[ "$base" =~ ^([^_]+)__(.+)$ ]]; then
    schema_name="${BASH_REMATCH[1]}"
    table_name="${BASH_REMATCH[2]}"
    
    # Remove .go extension from table_name if present
    table_name="${table_name%.go}"
    
    # Determine destination directory based on schema
    DEST_DIR="${REPO_ROOT}/modules/${schema_name}/infrastructure/rdb/models"
    mkdir -p "${DEST_DIR}"
    
    # Clean old files for this specific table
    rm -f "${DEST_DIR}/${table_name}_dbgen.go"
    
    # Rename to {table}_dbgen.go format
    dest_file="${DEST_DIR}/${table_name}_dbgen.go"
    mv "${src}" "${dest_file}"
    
    echo "Moved ${base} -> ${dest_file}"
    file_count=$((file_count + 1))
  else
    echo "Warning: Unexpected filename format: ${base}" >&2
  fi
done
shopt -u nullglob

if [ $file_count -eq 0 ]; then
  echo "Warning: No model files generated in ${GEN_DIR}" >&2
  echo "This might indicate that:"
  echo "  1. Tables are in 'public' schema (which is excluded)"
  echo "  2. Schema apply did not run successfully"
  echo "  3. Template generation failed"
  exit 1
fi

# -------------------------------
# Format generated Go files in all module directories
# -------------------------------
set +e  # Temporarily disable exit on error for formatting
for module_dir in "${REPO_ROOT}/modules"/*/infrastructure/rdb/models; do
  if [ -d "$module_dir" ]; then
    if command -v goimports >/dev/null 2>&1; then
      goimports -w "$module_dir" 2>&1 || echo "Warning: goimports failed for $module_dir (non-fatal)" >&2
    fi
    if command -v gofmt >/dev/null 2>&1; then
      gofmt -s -w "$module_dir" 2>&1 || echo "Warning: gofmt failed for $module_dir (non-fatal)" >&2
    fi
  fi
done
# Reset exit status to success before re-enabling set -e
true
set -e  # Re-enable exit on error

echo "Models generated successfully."
exit 0
