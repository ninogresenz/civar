#!/bin/bash
docker run \
  --rm \
  -v "$(pwd):/app" \
  -w /app \
  golang:1.19-alpine \
  /bin/sh -c "bin/compile.sh && chown $(id -u):$(id -u) civar-*"