#Requires -Version 5.1
<#
.SYNOPSIS
    Run ApuesTec backend locally with environment from .env file.
.DESCRIPTION
    Reads backend/.env, overrides connection strings for local Docker infra,
    and runs go run ./cmd/api.
.NOTES
    Requires: Go, Docker infra (postgres, redis, nginx) running via docker compose.
    Run this from the backend/ directory.
#>

$envFile = Join-Path $PSScriptRoot ".env"
if (-not (Test-Path $envFile)) {
    Write-Error ".env file not found at $envFile"
    exit 1
}

# Read .env file and set environment variables
Get-Content $envFile | ForEach-Object {
    $line = $_.Trim()
    if ($line -and $line -notmatch '^\s*#') {
        $eqIndex = $line.IndexOf('=')
        if ($eqIndex -gt 0) {
            $key = $line.Substring(0, $eqIndex).Trim()
            $value = $line.Substring($eqIndex + 1).Trim()
            Set-Item -Path "env:$key" -Value $value
        }
    }
}

# Local overrides — Docker containers are accessible via localhost mapped ports
Set-Item -Path "env:DATABASE_URL" -Value "postgres://apuestec_user:apuestec_password@127.0.0.1:5433/apuestec?sslmode=disable"
Set-Item -Path "env:REDIS_URL" -Value "redis://127.0.0.1:6379/0"
Set-Item -Path "env:CORS_ALLOWED_ORIGINS" -Value "http://localhost:8081,http://localhost:3000"

Write-Host "Starting ApuesTec backend locally..." -ForegroundColor Cyan
Write-Host "  DATABASE_URL:   postgres://apuestec_user:***@localhost:5433/apuestec" -ForegroundColor Gray
Write-Host "  REDIS_URL:      redis://localhost:6379/0" -ForegroundColor Gray
Write-Host "  CORS_ORIGINS:   http://localhost:8081, http://localhost:3000" -ForegroundColor Gray
Write-Host ""

go run ./cmd/api
