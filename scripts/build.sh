#!/bin/bash
echo "Building CLI application..."

# Build frontend
cd frontend && npm ci && npm run build
cd ..

# Copy frontend build to server static dir
rm -rf backend/cmd/server/static/*
cp -r frontend/build/* backend/cmd/server/static/

# Build Go binaries
cd backend
go build -o order-controller ./cmd/cli
go build -o order-server ./cmd/server

echo "Build completed"
