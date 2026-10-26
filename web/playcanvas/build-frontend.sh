#!/usr/bin/env bash
set -euo pipefail

# This script is meant to be run in the CI/CD pipeline or locally before Go build.
echo "🚀 Starting PlayCanvas Frontend Build..."

cd "$(dirname "$0")"

echo "📦 Installing Node dependencies..."
npm ci

echo "🔨 Bundling TypeScript with esbuild..."
npm run build

echo "✅ Frontend build complete. Output is in internal/web/static/dist/main.js"
