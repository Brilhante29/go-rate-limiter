param(
  [string]$Image = "go-rate-limiter",
  [int]$DurationSeconds = 10,
  [int]$Concurrency = 64,
  [double]$RatePerSecond = 1000,
  [int64]$Burst = 1000
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$results = Join-Path $root "benchmarks/results"
New-Item -ItemType Directory -Force -Path $results | Out-Null
$resolvedResults = (Resolve-Path -LiteralPath $results).Path

docker build -t $Image $root
if ($LASTEXITCODE -ne 0) { throw "Docker build failed" }

docker run --rm `
  --mount "type=bind,source=$resolvedResults,target=/results" `
  $Image `
  demo `
  --duration "${DurationSeconds}s" `
  --concurrency $Concurrency `
  --rate $RatePerSecond `
  --burst $Burst `
  --output /results/rate-limiter-baseline.json
if ($LASTEXITCODE -ne 0) { throw "Docker benchmark failed" }
