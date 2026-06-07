Write-Host "ApuesTec - Aplicacion manual de seeds"
Write-Host "Requisito: docker compose debe estar levantado y la migracion principal aplicada."
Write-Host "Ejecutando seed de roles globales..."

Get-Content -Raw .\database\seeds\001_seed_roles.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'

if ($LASTEXITCODE -ne 0) {
    Write-Error "La aplicacion del seed de roles fallo."
    exit $LASTEXITCODE
}

Write-Host "Ejecutando seed de equipos y partidos de prueba local..."

Get-Content -Raw .\database\seeds\002_seed_sample_teams_matches.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'

if ($LASTEXITCODE -ne 0) {
    Write-Error "La aplicacion del seed de equipos y partidos fallo."
    exit $LASTEXITCODE
}

Write-Host "Seeds aplicados correctamente."
