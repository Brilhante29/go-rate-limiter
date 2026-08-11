param(
  [string]$Image = "go-rate-limiter",
  [int]$DurationSeconds = 5,
  [int]$WarmupSeconds = 1,
  [int]$Repetitions = 3,
  [int]$Concurrency = 64,
  [double]$RatePerSecond = 1000,
  [int64]$Burst = 1000,
  [string]$HardwareClass = "local-docker"
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$results = Join-Path $root "benchmarks/results"
New-Item -ItemType Directory -Force -Path $results | Out-Null
$resolvedResults = (Resolve-Path -LiteralPath $results).Path
$workloadPath = Join-Path $root "benchmarks/workload.json"
$goSumPath = Join-Path $root "go.sum"

if ($Repetitions -lt 3) { throw "Repetitions must be at least 3" }
if ($WarmupSeconds -lt 1) { throw "WarmupSeconds must be at least 1" }

$treeState = @(git -C $root status --porcelain --untracked-files=normal)
if ($LASTEXITCODE -ne 0) { throw "Cannot inspect Git tree" }
if ($treeState.Count -ne 0) { throw "Benchmark requires a clean Git tree so provenance is exact" }
$sourceCommit = (git -C $root rev-parse HEAD).Trim()
if ($LASTEXITCODE -ne 0 -or $sourceCommit -notmatch '^[0-9a-f]{40}$') { throw "Cannot resolve exact source commit" }

$fixtureDigest = "sha256:" + (Get-FileHash -Algorithm SHA256 -LiteralPath $workloadPath).Hash.ToLowerInvariant()
$dependencyLockDigest = "sha256:" + (Get-FileHash -Algorithm SHA256 -LiteralPath $goSumPath).Hash.ToLowerInvariant()

docker build -t $Image $root
if ($LASTEXITCODE -ne 0) { throw "Docker build failed" }

$imageDigest = (docker image inspect $Image --format '{{.Id}}').Trim()
if ($LASTEXITCODE -ne 0 -or $imageDigest -notmatch '^sha256:[0-9a-f]{64}$') { throw "Cannot resolve image digest" }
$artifactOutput = (docker run --rm --entrypoint sha256sum $Image /usr/local/bin/go-rate-limiter).Trim()
if ($LASTEXITCODE -ne 0 -or $artifactOutput -notmatch '^([0-9a-f]{64})\s+') { throw "Cannot resolve binary artifact digest" }
$artifactDigest = "sha256:" + $Matches[1]

docker run --rm `
  --mount "type=bind,source=$resolvedResults,target=/results" `
  $Image `
  demo `
  --duration "${DurationSeconds}s" `
  --warmup-duration "${WarmupSeconds}s" `
  --warmup-iterations 1 `
  --measured-iterations $Repetitions `
  --concurrency $Concurrency `
  --rate $RatePerSecond `
  --burst $Burst `
  --fixture-digest $fixtureDigest `
  --source-commit $sourceCommit `
  --clean-tree=true `
  --image-ref $Image `
  --image-digest $imageDigest `
  --dependency-lock-digest $dependencyLockDigest `
  --producer local `
  --artifact-digest $artifactDigest `
  --hardware-class $HardwareClass `
  --redis-version 8.8.0 `
  --output /results/rate-limiter-baseline.json
if ($LASTEXITCODE -ne 0) { throw "Docker benchmark failed" }
