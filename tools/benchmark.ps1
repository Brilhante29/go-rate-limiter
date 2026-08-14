param(
  [string]$Image = "go-rate-limiter",
  [int]$DurationSeconds = 5,
  [int]$WarmupSeconds = 1,
  [int]$Repetitions = 3,
  [int]$Concurrency = 64,
  [double]$RatePerSecond = 1000,
  [int64]$Burst = 1000,
  [string]$HardwareClass = "local-docker",
  [string]$OutputPath = ""
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$resultPath = if ([string]::IsNullOrWhiteSpace($OutputPath)) {
  Join-Path $root "benchmarks/results/rate-limiter-baseline.json"
} elseif ([System.IO.Path]::IsPathRooted($OutputPath)) {
  $OutputPath
} else {
  Join-Path $root $OutputPath
}
$resultDirectory = Split-Path -Parent $resultPath
New-Item -ItemType Directory -Force -Path $resultDirectory | Out-Null
$resolvedResultDirectory = (Resolve-Path -LiteralPath $resultDirectory).Path
$resultFileName = Split-Path -Leaf $resultPath
$workloadPath = Join-Path $root "benchmarks/workload.json"
$goSumPath = Join-Path $root "go.sum"
$composeProject = "go-rate-limiter-bench-$PID"

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

$previousImage = $env:RATE_LIMITER_IMAGE
$previousNodeAPort = $env:NODE_A_PORT
$previousNodeBPort = $env:NODE_B_PORT
$previousRate = $env:RATE_PER_SECOND
$previousBurst = $env:BURST
$env:RATE_LIMITER_IMAGE = $Image
$env:NODE_A_PORT = "0"
$env:NODE_B_PORT = "0"
$env:RATE_PER_SECOND = [string]$RatePerSecond
$env:BURST = [string]$Burst

try {
  docker compose -p $composeProject up -d --wait redis node-a node-b
  if ($LASTEXITCODE -ne 0) { throw "Docker topology failed to become healthy" }

  docker run --rm `
    --network "${composeProject}_default" `
    --mount "type=bind,source=$resolvedResultDirectory,target=/results" `
    $Image `
    benchmark `
    --targets "http://node-a:8080,http://node-b:8080" `
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
    --output "/results/$resultFileName"
  if ($LASTEXITCODE -ne 0) { throw "Docker benchmark failed" }
} finally {
  docker compose -p $composeProject down --volumes --remove-orphans
  $env:RATE_LIMITER_IMAGE = $previousImage
  $env:NODE_A_PORT = $previousNodeAPort
  $env:NODE_B_PORT = $previousNodeBPort
  $env:RATE_PER_SECOND = $previousRate
  $env:BURST = $previousBurst
}
