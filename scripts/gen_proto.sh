#!/usr/bin/env bash
set -euo pipefail

# Generate Go protobuf code from proto/*.proto into internal/proto/gen
# NOTE: This batch does not ship .proto files yet; this script is ready for when they land.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTO_DIR="${ROOT_DIR}/proto"
OUT_DIR="${ROOT_DIR}/internal/proto/gen"

command -v protoc >/dev/null 2>&1 || {
  echo "ERROR: protoc not found on PATH."
  echo "Install protoc, then re-run."
  exit 1
}

command -v protoc-gen-go >/dev/null 2>&1 || {
  echo "ERROR: protoc-gen-go not found on PATH."
  echo "Install with:"
  echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  exit 1
}

mkdir -p "${OUT_DIR}"

shopt -s nullglob
PROTO_FILES=("${PROTO_DIR}"/*.proto)
shopt -u nullglob

if [ ${#PROTO_FILES[@]} -eq 0 ]; then
  echo "No .proto files found in ${PROTO_DIR}."
  echo "This is expected in Batch 00."
  exit 0
fi

protoc   -I "${PROTO_DIR}"   --go_out="${OUT_DIR}"   --go_opt=paths=source_relative   "${PROTO_FILES[@]}"

echo "Generated Go protobufs into: ${OUT_DIR}"
