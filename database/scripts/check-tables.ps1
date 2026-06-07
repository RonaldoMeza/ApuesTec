Write-Host "ApuesTec - Validacion manual de tablas y roles"
Write-Host "Requisito: docker compose debe estar levantado."
Write-Host "Listando tablas de la base configurada en el contenedor postgres..."

docker compose exec postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "\\dt"'

if ($LASTEXITCODE -ne 0) {
    Write-Error "La validacion de tablas fallo."
    exit $LASTEXITCODE
}

Write-Host "Listando roles iniciales..."

docker compose exec postgres sh -c "psql -U `"`$POSTGRES_USER`" -d `"`$POSTGRES_DB`" -c 'SELECT name FROM roles ORDER BY name;'"

if ($LASTEXITCODE -ne 0) {
    Write-Error "La validacion de roles fallo."
    exit $LASTEXITCODE
}

Write-Host "Validacion finalizada correctamente."
