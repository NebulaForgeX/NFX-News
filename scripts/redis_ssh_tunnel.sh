#!/usr/bin/env bash
# Interactive SSH local-forward to ElastiCache via a VPC jump EC2 (PulsoLoop-owned).
# Used by `task run` before docker compose up. Independent of PulsoLink's tunnel
# (default local port 16380; Link uses 16379).
#
# Home / outside VPC → answer y (kills any stale tunnel on the local port, then opens a fresh one).
# Already on VPC EC2 / private network → answer n (no tunnel; use REDIS_* as-is).
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

ENV_FILE=".env"
SIDECAR_ENV=".env.redis-tunnel"
PID_FILE=".redis-ssh-tunnel.pid"
DEFAULT_LOCAL_PORT="16380"

log() { printf '%s\n' "$*"; }
log_err() { printf '%s\n' "$*" >&2; }

# Read one KEY=VALUE from dotenv without eval (values may contain <empty>, spaces, quotes).
env_get_from() {
  local file="$1"
  local key="$2"
  local line
  [[ -f "$file" ]] || { printf ''; return 0; }
  line="$(grep -E "^${key}=" "$file" | tail -n1 | sed 's/\r$//' || true)"
  [[ -n "$line" ]] || { printf ''; return 0; }
  local val="${line#*=}"
  if [[ "$val" =~ ^\"(.*)\"$ ]]; then
    val="${BASH_REMATCH[1]}"
  elif [[ "$val" =~ ^\'(.*)\'$ ]]; then
    val="${BASH_REMATCH[1]}"
  fi
  printf '%s' "$val"
}

env_get() { env_get_from "$ENV_FILE" "$1"; }

resolve_local_port() {
  local port
  port="$(env_get REDIS_SSH_LOCAL_PORT)"
  if [[ -z "$port" && -f "$SIDECAR_ENV" ]]; then
    port="$(env_get_from "$SIDECAR_ENV" REDIS_DATABASE_PORT)"
  fi
  printf '%s' "${port:-$DEFAULT_LOCAL_PORT}"
}

# Kill recorded ssh pid and anything still bound to the local forward port.
stop_tunnel() {
  local local_port="$1"
  log "==> redis_ssh_tunnel --down (port :${local_port})"

  if [[ -f "$PID_FILE" ]]; then
    local pid
    pid="$(cat "$PID_FILE" 2>/dev/null || true)"
    log "  pid file: $PID_FILE → ${pid:-empty}"
    if [[ -n "${pid:-}" ]] && kill -0 "$pid" 2>/dev/null; then
      log "  killing ssh tunnel pid $pid ..."
      kill "$pid" 2>/dev/null || true
      sleep 0.2
      if kill -0 "$pid" 2>/dev/null; then
        log "  pid $pid still alive; kill -9 ..."
        kill -9 "$pid" 2>/dev/null || true
      fi
      log "  stopped pid $pid"
    else
      log "  no live process for recorded pid"
    fi
    rm -f "$PID_FILE"
  else
    log "  no pid file ($PID_FILE)"
  fi

  if command -v fuser >/dev/null 2>&1; then
    if fuser "${local_port}/tcp" >/dev/null 2>&1; then
      log "  fuser -k ${local_port}/tcp ..."
      fuser -k "${local_port}/tcp" >/dev/null 2>&1 || true
      sleep 0.2
    else
      log "  nothing listening on :${local_port}"
    fi
  elif command -v ss >/dev/null 2>&1; then
    local pids
    pids="$(ss -ltnp "( sport = :${local_port} )" 2>/dev/null | sed -n 's/.*pid=\([0-9]\+\).*/\1/p' | sort -u || true)"
    if [[ -n "${pids:-}" ]]; then
      log "  killing listeners on :${local_port}: $pids"
      # shellcheck disable=SC2086
      kill $pids 2>/dev/null || true
      sleep 0.2
      # shellcheck disable=SC2086
      kill -9 $pids 2>/dev/null || true
    fi
  fi

  if [[ -f "$SIDECAR_ENV" ]]; then
    log "  removing $SIDECAR_ENV"
    rm -f "$SIDECAR_ENV"
  fi
  log "==> down done"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --env-file)
      ENV_FILE="${2:-}"
      shift 2
      ;;
    --down)
      stop_tunnel "$(resolve_local_port)"
      exit 0
      ;;
    *)
      log_err "Unknown arg: $1"
      exit 2
      ;;
  esac
done

if [[ ! -f "$ENV_FILE" ]]; then
  log_err "ERROR: env file not found: $ENV_FILE"
  exit 1
fi

log "==> redis_ssh_tunnel"
log "  cwd:      $ROOT_DIR"
log "  env file: $ENV_FILE"

display_or_empty() {
  if [[ -n "${1:-}" ]]; then
    printf '%s' "$1"
  else
    printf '(empty)'
  fi
}

LOCAL_PORT="$(resolve_local_port)"
REMOTE_HOST="$(env_get REDIS_DATABASE_HOST)"
REMOTE_PORT="$(env_get REDIS_DATABASE_PORT)"
REMOTE_PORT="${REMOTE_PORT:-6379}"
JUMP_HOST="$(env_get REDIS_SSH_JUMP_HOST)"
JUMP_USER="$(env_get REDIS_SSH_JUMP_USER)"
JUMP_USER="${JUMP_USER:-ubuntu}"
IDENTITY="$(env_get REDIS_SSH_IDENTITY_FILE)"
REDIS_TLS="$(env_get REDIS_TLS)"
REDIS_TLS_SERVER_NAME="$(env_get REDIS_TLS_SERVER_NAME)"

log "  loaded from $ENV_FILE:"
log "    REDIS_DATABASE_HOST     = $(display_or_empty "$REMOTE_HOST")"
log "    REDIS_DATABASE_PORT     = $(display_or_empty "$REMOTE_PORT")"
log "    REDIS_TLS               = $(display_or_empty "$REDIS_TLS")"
log "    REDIS_TLS_SERVER_NAME   = $(display_or_empty "$REDIS_TLS_SERVER_NAME")"
log "    REDIS_SSH_JUMP_HOST     = $(display_or_empty "$JUMP_HOST")"
log "    REDIS_SSH_JUMP_USER     = $(display_or_empty "$JUMP_USER")"
log "    REDIS_SSH_IDENTITY_FILE = $(display_or_empty "$IDENTITY")"
log "    REDIS_SSH_LOCAL_PORT    = $(display_or_empty "$LOCAL_PORT")"

port_listening() {
  if command -v ss >/dev/null 2>&1; then
    log "  checking listen :${LOCAL_PORT} via ss ..."
    if ss -ltn "( sport = :${LOCAL_PORT} )" | grep -q ":${LOCAL_PORT}"; then
      log "  port :${LOCAL_PORT} IS listening"
      return 0
    fi
    log "  port :${LOCAL_PORT} is NOT listening"
    return 1
  fi
  log "  ss not found; probing 127.0.0.1:${LOCAL_PORT} via /dev/tcp ..."
  if (echo >/dev/tcp/127.0.0.1/"$LOCAL_PORT") 2>/dev/null; then
    log "  port :${LOCAL_PORT} IS open"
    return 0
  fi
  log "  port :${LOCAL_PORT} is NOT open"
  return 1
}

write_sidecar() {
  local tls_name="$1"
  cat >"$SIDECAR_ENV" <<EOF
# Generated by scripts/redis_ssh_tunnel.sh — do not commit.
REDIS_DATABASE_HOST=host.docker.internal
REDIS_DATABASE_PORT=${LOCAL_PORT}
REDIS_TLS=true
REDIS_TLS_SERVER_NAME=${tls_name}
EOF
  log "Wrote $SIDECAR_ENV:"
  sed 's/^/  /' "$SIDECAR_ENV"
  log "  containers will dial host.docker.internal:${LOCAL_PORT} with TLS SNI=${tls_name}"
}

print_fill_help() {
  log_err ""
  log_err "To open a tunnel from home, set these in $ENV_FILE then re-run:"
  log_err "  REDIS_SSH_JUMP_HOST=<jump EC2 public IP or DNS>"
  log_err "  REDIS_SSH_JUMP_USER=ubuntu"
  log_err "  REDIS_SSH_IDENTITY_FILE=/absolute/path/to/your.pem"
  log_err "  REDIS_SSH_LOCAL_PORT=${DEFAULT_LOCAL_PORT}   # optional"
  log_err ""
  log_err "Example:"
  log_err "  REDIS_SSH_JUMP_HOST=ec2-12-34-56-78.compute-1.amazonaws.com"
  log_err "  REDIS_SSH_IDENTITY_FILE=/home/lucas/Keys/NebulaForgeX-Lucas-Key.pem"
  log_err ""
  log_err "If you are already on the jump EC2 / inside the VPC, answer n next time."
}

log ""
log "Redis / ElastiCache access"
log "  ElastiCache has no public IP. From home you need:"
log "    laptop → SSH tunnel → jump EC2 (same VPC) → ElastiCache"
log "  Security groups cannot make ElastiCache reachable from the public internet."
log ""
log "  y = kill any stale local tunnel, then open a fresh SSH forward"
log "  n = skip (already on jump EC2 or inside the VPC)"
log ""
printf "Open Redis SSH tunnel? [y/N] "
read -r answer || answer="n"
answer="$(printf '%s' "$answer" | tr '[:upper:]' '[:lower:]')"
log "  answer: ${answer:-empty}"

if [[ "$answer" != "y" && "$answer" != "yes" ]]; then
  if [[ -f "$SIDECAR_ENV" ]]; then
    log "Removing stale $SIDECAR_ENV (tunnel skipped)."
    rm -f "$SIDECAR_ENV"
  fi
  log "Skipped Redis SSH tunnel; compose will use REDIS_* from $ENV_FILE as-is."
  log "==> done (skip)"
  exit 0
fi

log "Validating tunnel prerequisites ..."
ok=1
if [[ -z "$REMOTE_HOST" ]]; then
  log_err "ERROR: REDIS_DATABASE_HOST is empty in $ENV_FILE"
  ok=0
else
  log "  OK REDIS_DATABASE_HOST=$REMOTE_HOST"
fi
if [[ -z "$JUMP_HOST" ]]; then
  log_err "ERROR: REDIS_SSH_JUMP_HOST is empty in $ENV_FILE (EC2 public IP/DNS required)"
  ok=0
else
  log "  OK REDIS_SSH_JUMP_HOST=$JUMP_HOST"
fi
if [[ -z "$IDENTITY" ]]; then
  log_err "ERROR: REDIS_SSH_IDENTITY_FILE is empty in $ENV_FILE (path to .pem required)"
  ok=0
else
  IDENTITY="${IDENTITY/#\~/$HOME}"
  log "  OK REDIS_SSH_IDENTITY_FILE=$IDENTITY"
  if [[ ! -f "$IDENTITY" ]]; then
    log_err "ERROR: identity file not found: $IDENTITY"
    ok=0
  else
    log "  OK identity file exists ($(ls -la "$IDENTITY" | awk '{print $1, $3, $4, $5, $9}'))"
  fi
fi
if [[ "$ok" -ne 1 ]]; then
  print_fill_help
  exit 1
fi

# Always clear first so wedged ssh (LISTEN but dead forward) cannot be reused.
stop_tunnel "$LOCAL_PORT"

log "Starting SSH tunnel:"
log "  local bind:  0.0.0.0:${LOCAL_PORT}"
log "  remote:      ${REMOTE_HOST}:${REMOTE_PORT}"
log "  via:         ${JUMP_USER}@${JUMP_HOST}"
log "  identity:    ${IDENTITY}"
log "  ssh cmd:     ssh -v -f -N -o ExitOnForwardFailure=yes -o ServerAliveInterval=30 -o ServerAliveCountMax=3 -i ${IDENTITY} -L 0.0.0.0:${LOCAL_PORT}:${REMOTE_HOST}:${REMOTE_PORT} ${JUMP_USER}@${JUMP_HOST}"
log ""

set +e
ssh -v -f -N \
  -o ExitOnForwardFailure=yes \
  -o ServerAliveInterval=30 \
  -o ServerAliveCountMax=3 \
  -i "$IDENTITY" \
  -L "0.0.0.0:${LOCAL_PORT}:${REMOTE_HOST}:${REMOTE_PORT}" \
  "${JUMP_USER}@${JUMP_HOST}"
ssh_rc=$?
set -e
log "ssh exit code: $ssh_rc"
if [[ "$ssh_rc" -ne 0 ]]; then
  log_err "ERROR: ssh failed with exit code $ssh_rc (see -v output above)"
  exit "$ssh_rc"
fi

log "Waiting for local port :${LOCAL_PORT} ..."
sleep 0.5
if ! port_listening; then
  log_err "ERROR: tunnel did not open on :${LOCAL_PORT} after ssh returned 0"
  log_err "  Check jump SG (22 from your IP), key permissions, and that ElastiCache SG allows the jump."
  exit 1
fi

if command -v fuser >/dev/null 2>&1; then
  log "Recording listener pid via fuser ${LOCAL_PORT}/tcp ..."
  fuser "${LOCAL_PORT}/tcp" | awk '{print $1}' | head -1 >"$PID_FILE" || true
fi
if [[ ! -s "$PID_FILE" ]]; then
  log "fuser empty; falling back to pgrep ..."
  pgrep -af "ssh .*${LOCAL_PORT}:${REMOTE_HOST}:${REMOTE_PORT}" || true
  pgrep -f "ssh .*${LOCAL_PORT}:${REMOTE_HOST}:${REMOTE_PORT}" | head -1 >"$PID_FILE" || true
fi
if [[ -s "$PID_FILE" ]]; then
  log "Recorded tunnel pid $(cat "$PID_FILE") → $PID_FILE"
else
  log_err "warn: could not record tunnel pid (tunnel may still be up; use redis:tunnel:down / kill manually)"
fi

log "Redis SSH tunnel is up (local :${LOCAL_PORT})."
write_sidecar "$REMOTE_HOST"
log "==> done (tunnel started)"
