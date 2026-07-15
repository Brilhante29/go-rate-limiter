param(
  [switch]$SkipDocker
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$failures = New-Object System.Collections.Generic.List[string]

function Add-Failure {
  param([string]$Message)
  $script:failures.Add($Message)
}

function Require-File {
  param([string]$RelativePath)
  $path = Join-Path $root $RelativePath
  if (-not (Test-Path -LiteralPath $path -PathType Leaf)) {
    Add-Failure "Missing file: $RelativePath"
  }
}

function Invoke-Checked {
  param(
    [string]$Label,
    [scriptblock]$Command
  )
  & $Command
  $exitCode = $LASTEXITCODE
  if ($exitCode -ne 0) {
    Add-Failure "$Label failed with exit code $exitCode"
  }
  $global:LASTEXITCODE = 0
}

$requiredFiles = @(
  "README.md",
  "project.yaml",
  "REFERENCES.md",
  "AGENTS.md",
  "sdd/spec.md",
  "sdd/benchmark-plan.md",
  "sdd/architecture-decision.md",
  "sdd/technical-decision.md",
  "sdd/agent-handoff.md",
  "sdd/reuse-improvement-review.md"
)
foreach ($file in $requiredFiles) { Require-File $file }

$reuseReviewPath = Join-Path $root "sdd/reuse-improvement-review.md"
if (Test-Path -LiteralPath $reuseReviewPath -PathType Leaf) {
  $reuseReview = Get-Content -Raw -LiteralPath $reuseReviewPath
  if ($reuseReview -match "<id>|<project-name>") {
    Add-Failure "Reuse improvement review still contains template placeholders"
  }
  if ($reuseReview.Contains('|  | `patch_now|backlog|reject` |')) {
    Add-Failure "Reuse improvement review still contains the blank template finding row"
  }
  $requiredFinalGatePatterns = @(
    "(?m)^- \[x\] Reusable improvements were patched or recorded\.\r?$",
    "(?m)^- \[x\] Project-specific implementation was not moved into the kit\.\r?$",
    "(?m)^- \[x\] Validation reflects .+\.\r?$"
  )
  foreach ($pattern in $requiredFinalGatePatterns) {
    if ($reuseReview -notmatch $pattern) {
      Add-Failure "Reuse improvement review final gate is incomplete: $pattern"
    }
  }
}

$benchmarkFiles = @()
$benchmarkDir = Join-Path $root "benchmarks/results"
if (Test-Path -LiteralPath $benchmarkDir -PathType Container) {
  $benchmarkFiles = @(Get-ChildItem -LiteralPath $benchmarkDir -Filter *.json -File)
}
if ($benchmarkFiles.Count -eq 0) {
  Add-Failure "Missing benchmark JSON under benchmarks/results"
}

Push-Location -LiteralPath $root
try {
  foreach ($file in $benchmarkFiles) {
    Invoke-Checked "benchmark JSON validation: $($file.Name)" { python -m json.tool $file.FullName | Out-Null }
  }

  if (Test-Path -LiteralPath (Join-Path $root "src") -PathType Container) {
    $previousPythonPath = $env:PYTHONPATH
    $srcPath = Join-Path $root "src"
    if ($previousPythonPath) {
      $env:PYTHONPATH = $srcPath + [System.IO.Path]::PathSeparator + $previousPythonPath
    } else {
      $env:PYTHONPATH = $srcPath
    }
    Invoke-Checked "python compile src" { python -m compileall -q (Join-Path $root "src") }
    if (Test-Path -LiteralPath (Join-Path $root "tests") -PathType Container) {
      Invoke-Checked "python compile tests" { python -m compileall -q (Join-Path $root "tests") }
      Invoke-Checked "python unittest" { python -m unittest discover -s (Join-Path $root "tests") -v }
    }
    $env:PYTHONPATH = $previousPythonPath
  }
} finally {
  Pop-Location
}

$legacy = ("ro" + "che" + "do")
$patterns = @($legacy, ($legacy.Substring(0,1).ToUpper() + $legacy.Substring(1)))
$searchFiles = Get-ChildItem -Path $root -Recurse -File | Where-Object {
  $normalized = $_.FullName -replace "\\", "/"
  $normalized -notmatch "/.git/" -and
  $normalized -notmatch "/data/runtime/" -and
  $_.Extension -in @(".md", ".yaml", ".yml", ".json", ".ps1", ".py", ".js", ".ts", ".tsx", ".go", ".kt", ".java")
}
$forbidden = Select-String -Path $searchFiles.FullName -Pattern $patterns -SimpleMatch -ErrorAction SilentlyContinue
if ($forbidden) {
  Add-Failure "Forbidden legacy project nickname found"
}

if (-not $SkipDocker -and (Test-Path -LiteralPath (Join-Path $root "Dockerfile") -PathType Leaf)) {
  $imageName = (Split-Path -Leaf $root).ToLowerInvariant()
  Invoke-Checked "docker build" { docker build -t $imageName $root | Out-Null }
}

if ($failures.Count -gt 0) {
  $failures | ForEach-Object { Write-Error $_ }
  exit 1
}

Write-Host "portfolio project validation passed"
