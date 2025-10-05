#!/usr/bin/env bash
#
# Extract each branch from a source Git repository into its own bare Git repo.
#
# Usage:
#   ./update_mirror.sh <source_repo_url> <working_directory> <output_directory>
#
# Example:
#   ./update_mirror.sh https://github.com/archlinux/aur.git ./upstream ./repos
#
# Behavior:
#   - Clones the source repo into <working_directory> if not already cloned.
#   - Fetches updates if it exists.
#   - Exports each branch into its own bare repository under <output_directory>.
#   - Sets `http.receivepack=false` for each target repo.
#

set -euo pipefail

SRC_URL="${1:-}"
WORK_DIR="${2:-}"
DEST_DIR="${3:-}"

if [[ -z "$SRC_URL" || -z "$WORK_DIR" || -z "$DEST_DIR" ]]; then
  echo "Usage: $0 <source_repo_url> <working_directory> <output_directory>"
  exit 1
fi

mkdir -p "$WORK_DIR"
mkdir -p "$DEST_DIR"

REPO_NAME=$(basename "$SRC_URL" .git)
SRC_REPO="$WORK_DIR/$REPO_NAME"

# --- Clone or fetch the source repository ---
if [[ -d "$SRC_REPO/.git" ]]; then
  echo "📦 Updating existing repository at $SRC_REPO"
  cd "$SRC_REPO"
  git remote set-url origin "$SRC_URL" || true
  git fetch --all --prune
else
  echo "⬇️ Cloning source repository from $SRC_URL into $SRC_REPO"
  git clone --mirror "$SRC_URL" "$SRC_REPO"
  cd "$SRC_REPO"
  # --mirror gives us all refs (branches, tags, etc.)
fi

# Ensure we're in the repo
cd "$SRC_REPO"

# --- Collect all branches (both local + remote) ---
echo "🔍 Collecting branches..."
BRANCHES=$(git for-each-ref --format='%(refname:short)' refs/heads/ refs/remotes/ \
  | grep -vE 'HEAD$' \
  | sed 's#^origin/##' \
  | sort -u)

echo "Found branches:"
echo "$BRANCHES"
echo

# --- Export each branch into its own bare repo ---
for BRANCH in $BRANCHES; do
  echo "=== Processing branch: $BRANCH ==="

  TARGET_REPO="$DEST_DIR/${BRANCH}.git"

  if [[ ! -d "$TARGET_REPO" ]]; then
    echo "📁 Creating bare repo at $TARGET_REPO"
    git init --bare "$TARGET_REPO"
    git --git-dir="$TARGET_REPO" config http.receivepack false
  fi

  echo "🚀 Pushing branch '$BRANCH' to $TARGET_REPO"
  git push "$TARGET_REPO" "refs/remotes/origin/$BRANCH:refs/heads/$BRANCH" || \
  git push "$TARGET_REPO" "$BRANCH:$BRANCH"

  echo "✅ Done with branch: $BRANCH"
  echo
done

echo "🎉 All branches exported to $DEST_DIR"
