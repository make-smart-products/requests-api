#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> Fetch remote"
git fetch origin

echo "==> Merge remote main"
if ! git pull origin main --rebase; then
  echo "==> Rebase failed, trying unrelated histories merge"
  git rebase --abort 2>/dev/null || true
  git pull origin main --allow-unrelated-histories -m "merge remote main"
fi

if git diff --name-only --diff-filter=U | grep -q '^README.md$'; then
  echo "==> Resolve README conflict (keep local)"
  git checkout --ours README.md
  git add README.md
  if [ -d .git/rebase-merge ] || [ -d .git/rebase-apply ]; then
    git rebase --continue
  else
    git commit --no-edit || git commit -m "merge remote main"
  fi
fi

echo "==> Push"
git push -u origin main

echo "==> Done"
git status -sb
