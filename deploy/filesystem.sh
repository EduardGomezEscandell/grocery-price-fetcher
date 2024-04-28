#!/bin/bash
set -eux

ROOT=root
PKG=grocery-server

if [ $1 == "clean" ]; then
    rm -rf "${ROOT}"
    exit 0
elif [ $1 != "build" ]; then
    echo "Usage: $0 [build|clean]"
    exit 1
fi

mkdir -p "${ROOT}/usr/bin"
mkdir -p "${ROOT}/etc/${PKG}"
mkdir -p "${ROOT}/usr/share/${PKG}/frontend"
mkdir -p "${ROOT}/usr/share/${PKG}/storage"

# Config
cp "config.yaml"            "${ROOT}/etc/${PKG}/config.yaml"
cp "database.json"          "${ROOT}/usr/share/${PKG}/database.json"

# Build
cp "../bin/grocery-server"  "${ROOT}/usr/bin/${PKG}"
cp -R "../frontend/build/." "${ROOT}/usr/share/${PKG}/frontend/"

