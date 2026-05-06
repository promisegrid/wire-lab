#!/bin/bash
# bin/ppx-bootstrap.sh
#
# Idempotent session-logging bootstrap for the wire-lab Perplexity Computer
# environment. See AGENTS-ppx.md "Session Logging" section for context.
#
# Steps (each is a no-op if already satisfied):
#   1. Identity gate: only run when git config user.email is Steve's bot.
#   2. Verify or create ~/.creds/ directory (chmod 700).
#   3. Verify ~/.creds/session-logs.pat exists; if missing, prompt.
#   4. Verify or add `private` remote pointing at the session-logs repo.
#   5. Verify or wire the credential helper.
#   6. Verify or create the worktree at /home/user/workspace/wire-lab-logs.
#   7. Fetch private/wire-lab.
#   8. Probe credentials with a dry-run push.
#
# Exit codes:
#   0  - bootstrap successful, or skipped because identity gate failed
#   1  - PAT missing and not interactively provided
#   2  - other failure

set -euo pipefail

REPO_DIR="/home/user/workspace/wire-lab"
LOGS_DIR="/home/user/workspace/wire-lab-logs"
CREDS_DIR="$HOME/.creds"
PAT_FILE="$CREDS_DIR/session-logs.pat"
PRIVATE_URL="https://github.com/stevegt/session-logs.git"
HELPER_PATH="$REPO_DIR/bin/git-cred-private"
EXPECTED_EMAILS=("stevegt+ppx@t7a.org" "stevegt@t7a.org")

cd "$REPO_DIR"

# Step 1: identity gate
CURRENT_EMAIL="$(git config user.email 2>/dev/null || true)"
GATE_OK=0
for e in "${EXPECTED_EMAILS[@]}"; do
  if [ "$CURRENT_EMAIL" = "$e" ]; then
    GATE_OK=1
    break
  fi
done

if [ "$GATE_OK" -eq 0 ]; then
  echo "ppx-bootstrap: identity gate not matched (user.email=$CURRENT_EMAIL); session logging inactive."
  exit 0
fi

echo "ppx-bootstrap: identity gate ok ($CURRENT_EMAIL)"

# Step 2: ~/.creds/
if [ ! -d "$CREDS_DIR" ]; then
  mkdir -p "$CREDS_DIR"
  chmod 700 "$CREDS_DIR"
  echo "ppx-bootstrap: created $CREDS_DIR"
fi

# Step 3: PAT file
if [ ! -f "$PAT_FILE" ]; then
  if [ -t 0 ]; then
    echo "ppx-bootstrap: $PAT_FILE not found."
    echo "Paste the github_pat_... PAT for the stevegt/session-logs repo, then press Enter:"
    read -r PAT_INPUT
    if [ -z "$PAT_INPUT" ]; then
      echo "ppx-bootstrap: empty PAT; aborting." >&2
      exit 1
    fi
    printf '%s\n' "$PAT_INPUT" > "$PAT_FILE"
    chmod 600 "$PAT_FILE"
    echo "ppx-bootstrap: wrote $PAT_FILE"
  else
    echo "ppx-bootstrap: $PAT_FILE missing and stdin is not a tty; cannot prompt." >&2
    echo "ppx-bootstrap: write the PAT to $PAT_FILE (chmod 600) and re-run." >&2
    exit 1
  fi
fi

# Step 4: private remote
if git remote get-url private >/dev/null 2>&1; then
  CURRENT_URL="$(git remote get-url private)"
  if [ "$CURRENT_URL" != "$PRIVATE_URL" ]; then
    echo "ppx-bootstrap: rewriting private remote URL ($CURRENT_URL -> $PRIVATE_URL)"
    git remote set-url private "$PRIVATE_URL"
  fi
else
  git remote add private "$PRIVATE_URL"
  echo "ppx-bootstrap: added private remote"
fi

# Step 5: credential helper
HELPER_KEY="credential.${PRIVATE_URL}.helper"
if ! git config --global --get "$HELPER_KEY" >/dev/null 2>&1; then
  git config --global "$HELPER_KEY" "$HELPER_PATH"
  echo "ppx-bootstrap: wired credential helper"
fi

# Step 6: worktree
if [ -d "$LOGS_DIR/.git" ] || [ -f "$LOGS_DIR/.git" ]; then
  echo "ppx-bootstrap: worktree at $LOGS_DIR already exists"
else
  # Need to fetch first so private/wire-lab is known locally
  git fetch private 2>&1 | grep -v "^From " || true
  if git rev-parse --verify --quiet private/wire-lab >/dev/null; then
    git worktree add "$LOGS_DIR" -B wire-lab private/wire-lab
    echo "ppx-bootstrap: created worktree at $LOGS_DIR on wire-lab branch"
  else
    echo "ppx-bootstrap: private/wire-lab branch not found; cannot create worktree." >&2
    exit 2
  fi
fi

# Step 7: fetch
git fetch private 2>&1 | grep -v "^From " || true

# Step 8: dry-run push probe
if git -C "$LOGS_DIR" push --dry-run private wire-lab >/dev/null 2>&1; then
  echo "ppx-bootstrap: push channel ok"
else
  echo "ppx-bootstrap: push probe failed; check PAT scope and repo access." >&2
  exit 2
fi

echo "ppx-bootstrap: complete"
