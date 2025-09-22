# backend/scripts/test-e2e.ps1
# Corre el test de integración e2e sin Docker (modernc.org/sqlite)
# Funciona se lance desde donde se lance.

# 1) Localiza la carpeta backend/ a partir de la ruta del script
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendDir = Split-Path -Parent $ScriptDir  # backend/

if (-not (Test-Path (Join-Path $BackendDir 'go.mod'))) {
  Write-Error "No encuentro go.mod en $BackendDir. Asegúrate de que este script está en backend/scripts/"
  exit 1
}

# 2) Entra en backend/
Set-Location -Path $BackendDir

# 3) Asegura dependencias (por si en local faltan)
go get github.com/gin-contrib/cors@v1.5.0
go mod tidy

# 4) Ejecuta el test E2E
go test ./internal/integration -run Test_FullAPI_HappyPath -v
