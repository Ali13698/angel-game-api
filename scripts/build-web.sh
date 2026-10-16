#!/bin/sh
set -e
if ! command -v esbuild >/dev/null 2>&1; then
  echo "[build-web] esbuild not found, skipping"
  exit 0
fi
mkdir -p internal/web/static
cp web/profile/index.html internal/web/static/
cp web/profile/style.css internal/web/static/
esbuild web/profile/app.ts --bundle --minify --target=es2020 --outfile=internal/web/static/app.js
echo "[build-web] OK"
