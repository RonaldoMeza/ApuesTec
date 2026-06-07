Write-Host "ApuesTec - Reset manual de base de datos local"
Write-Host "ADVERTENCIA: este script elimina tablas y datos locales."
Write-Host "Usar solo en desarrollo local. No se ejecuta automaticamente desde Docker Compose."
Write-Host "Requisito: docker compose debe estar levantado."

$confirmation = Read-Host "Escribe RESET para eliminar y recrear la base local"

if ($confirmation -ne "RESET") {
    Write-Host "Operacion cancelada. No se realizaron cambios."
    exit 0
}

Write-Host "Ejecutando migracion down: elimina tablas y datos locales..."

Get-Content -Raw .\database\migrations\001_init_schema.down.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'

if ($LASTEXITCODE -ne 0) {
    Write-Error "La migracion down fallo."
    exit $LASTEXITCODE
}

Write-Host "Ejecutando migracion up: recrea tablas..."

Get-Content -Raw .\database\migrations\001_init_schema.up.sql | docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1'

if ($LASTEXITCODE -ne 0) {
    Write-Error "La migracion up fallo."
    exit $LASTEXITCODE
}

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

Write-Host "Reset de base local finalizado correctamente."
