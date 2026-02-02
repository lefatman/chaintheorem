# Generate Go protobuf code from proto\*.proto into internal\proto\gen
# NOTE: This batch does not ship .proto files yet; this script is ready for when they land.

$ErrorActionPreference = "Stop"

$RootDir  = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$ProtoDir = Join-Path $RootDir "proto"
$OutDir   = Join-Path $RootDir "internal\proto\gen"

function Require-Cmd($name) {
  $cmd = Get-Command $name -ErrorAction SilentlyContinue
  if (-not $cmd) {
    throw "ERROR: $name not found on PATH."
  }
}

Require-Cmd "protoc"
Require-Cmd "protoc-gen-go"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$protoFiles = Get-ChildItem -Path $ProtoDir -Filter "*.proto" -File -ErrorAction SilentlyContinue
if (-not $protoFiles -or $protoFiles.Count -eq 0) {
  Write-Host "No .proto files found in $ProtoDir."
  Write-Host "This is expected in Batch 00."
  exit 0
}

$protoPaths = $protoFiles | ForEach-Object { $_.FullName }

& protoc `
  -I $ProtoDir `
  --go_out=$OutDir `
  --go_opt=paths=source_relative `
  $protoPaths

Write-Host "Generated Go protobufs into: $OutDir"
