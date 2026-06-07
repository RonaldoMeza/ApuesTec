Write-Host "ApuesTec - Aplicacion manual de migraciones"
Write-Host "Requisito: docker compose debe estar levantado antes de ejecutar este script."
Write-Host "Ejecutando migracion principal sobre la base configurada en el contenedor postgres..."

$existingRelations = docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\\dt"'

if ($LASTEXITCODE -ne 0) {
    Write-Error "No se pudo consultar el estado actual de la base."
    exit $LASTEXITCODE
}

$relationsOutput = $existingRelations | Out-String

if ($relationsOutput -notmatch "Did not find any relations") {
    Write-Host "La base ya contiene tablas publicas. No se reejecuta la migracion para evitar modificar datos existentes."
    Write-Host "Si necesitas recrear la base local de desarrollo, usa .\database\scripts\reset-database.ps1."
    exit 0
}

Get-Content -Raw .\database\migrations\001_init_schema.up.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'

if ($LASTEXITCODE -ne 0) {
    Write-Error "La aplicacion de migraciones fallo."
    exit $LASTEXITCODE
}

Write-Host "Migracion principal aplicada correctamente."
