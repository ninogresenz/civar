#!/bin/sh
for os in linux darwin; do
  echo "compile ${os} binary..."
  env GOOS="${os}" GOARCH=amd64 CGO_ENABLED=0 go build -o "civar-${os}" --ldflags="-w -s"
done
echo "compilation done."