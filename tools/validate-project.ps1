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
  "openspec/config.yaml",
  "openspec/artifacts/intent.md",
  "openspec/artifacts/portfolio-impact.md",
  "openspec/artifacts/architecture-record.md",
  "openspec/artifacts/component-pack.md",
  "openspec/artifacts/reuse-delta.md",
  "openspec/artifacts/benchmark-proof.md",
  "openspec/artifacts/tasks.md",
  "openspec/artifacts/verification.md",
  "openspec/artifacts/article-draft.md",
  "openspec/artifacts/voice-check.md",
  "sdd/spec.md",
  "sdd/benchmark-plan.md",
  "sdd/architecture-decision.md",
  "sdd/technical-decision.md",
  "sdd/agent-handoff.md",
  "sdd/reuse-improvement-review.md"
)
foreach ($file in $requiredFiles) { Require-File $file }

$manifestPath = Join-Path $root "project.yaml"
$manifestPrimaryMetric = ""
$manifestResultPath = ""
if (Test-Path -LiteralPath $manifestPath -PathType Leaf) {
  $manifestText = Get-Content -Raw -LiteralPath $manifestPath
  if ($manifestText.Contains("`t")) {
    Add-Failure "project.yaml must not contain tab indentation"
  }

  $requiredTopLevelKeys = @("id", "name", "program", "status", "stack", "benchmark", "release")
  foreach ($key in $requiredTopLevelKeys) {
    if ($manifestText -notmatch "(?m)^$([regex]::Escape($key)):\s*") {
      Add-Failure "project.yaml is missing top-level key: $key"
    }
  }

  foreach ($line in ($manifestText -split "`r?`n")) {
    if ($line -match "^( +)\S" -and ($Matches[1].Length % 2) -ne 0) {
      Add-Failure "project.yaml indentation must use multiples of two spaces"
      break
    }
  }

  $primaryMetricMatch = [regex]::Match($manifestText, "(?m)^\s+primary_metric:\s*([^\r\n#]+)")
  if ($primaryMetricMatch.Success) {
    $manifestPrimaryMetric = $primaryMetricMatch.Groups[1].Value.Trim().Trim('"').Trim("'")
  } else {
    Add-Failure "project.yaml is missing benchmark.primary_metric"
  }

  $resultPathMatch = [regex]::Match($manifestText, "(?m)^\s+result_path:\s*([^\r\n#]+)")
  if ($resultPathMatch.Success) {
    $manifestResultPath = $resultPathMatch.Groups[1].Value.Trim().Trim('"').Trim("'")
  } else {
    Add-Failure "project.yaml is missing benchmark.result_path"
  }

  $pythonCommand = Get-Command python -ErrorAction SilentlyContinue
  if ($pythonCommand) {
    & python -c "import yaml" 2>$null
    $yamlAvailable = $LASTEXITCODE -eq 0
    $global:LASTEXITCODE = 0
    if ($yamlAvailable) {
      Invoke-Checked "project YAML parsing" { python -c "import pathlib, sys, yaml; data = yaml.safe_load(pathlib.Path(sys.argv[1]).read_text(encoding='utf-8')); assert isinstance(data, dict)" $manifestPath }
    } else {
      Write-Host "yaml_parser=not_found; structural manifest validation applied"
    }
  }
}
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

if ($manifestResultPath -ne "") {
  $primaryResultPath = $root
  foreach ($part in ($manifestResultPath -split "[/\\]")) {
    if ($part -ne "") {
      $primaryResultPath = Join-Path $primaryResultPath $part
    }
  }

  if (-not (Test-Path -LiteralPath $primaryResultPath -PathType Leaf)) {
    Add-Failure "Manifest benchmark result does not exist: $manifestResultPath"
  } else {
    try {
      $primaryResult = Get-Content -Raw -LiteralPath $primaryResultPath | ConvertFrom-Json
      $resultMetric = ""
      $resultValue = $null
      if ($primaryResult.PSObject.Properties.Name -contains "metric") {
        $resultMetric = [string]$primaryResult.metric
        $resultValue = $primaryResult.value
      } elseif ($primaryResult.PSObject.Properties.Name -contains "primary_metric") {
        $resultMetric = [string]$primaryResult.primary_metric
        $metricProperty = $primaryResult.PSObject.Properties[$resultMetric]
        if ($null -ne $metricProperty) {
          $resultValue = $metricProperty.Value
        }
      }

      if ($resultMetric -eq "" -or $null -eq $resultValue) {
        Add-Failure "Benchmark result must expose metric/value or primary_metric with its value"
      } else {
        if ($manifestPrimaryMetric -ne "" -and $resultMetric -ne $manifestPrimaryMetric) {
          Add-Failure "Benchmark metric mismatch: project.yaml=$manifestPrimaryMetric result=$resultMetric"
        }
        $readmePath = Join-Path $root "README.md"
        if (Test-Path -LiteralPath $readmePath -PathType Leaf) {
          $readmeOpening = ((Get-Content -LiteralPath $readmePath -TotalCount 8) -join "`n").Replace(",", "")
          $valueText = [Convert]::ToString($resultValue, [System.Globalization.CultureInfo]::InvariantCulture)
          if (-not $readmeOpening.Contains($valueText)) {
            Add-Failure "README opening does not include primary benchmark value: $valueText"
          }
        }
      }
    } catch {
      Add-Failure "Cannot read primary benchmark result: $($_.Exception.Message)"
    }
  }
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

  $goModPath = Join-Path $root "go.mod"
  if (Test-Path -LiteralPath $goModPath -PathType Leaf) {
    $goCommand = Get-Command go -ErrorAction SilentlyContinue
    if ($goCommand) {
      $goFiles = @(Get-ChildItem -Path $root -Recurse -Filter *.go -File | Where-Object {
        $normalized = $_.FullName -replace "\\", "/"
        $normalized -notmatch "/.git/" -and
        $normalized -notmatch "/.portfolio/" -and
        $normalized -notmatch "/vendor/"
      })
      if ($goFiles.Count -gt 0) {
        $unformatted = @(& gofmt -l $goFiles.FullName)
        if ($LASTEXITCODE -ne 0) {
          Add-Failure "gofmt failed with exit code $LASTEXITCODE"
          $global:LASTEXITCODE = 0
        } elseif ($unformatted.Count -gt 0) {
          Add-Failure "Go files require gofmt: $($unformatted -join ', ')"
        }
      }
      Invoke-Checked "go test" { go test ./... }
      Invoke-Checked "go vet" { go vet ./... }
    } elseif ($SkipDocker) {
      Add-Failure "Go toolchain is required to validate go.mod projects when Docker validation is skipped"
    } else {
      Write-Host "go_toolchain=not_found; relying on Docker build for Go validation"
    }
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
