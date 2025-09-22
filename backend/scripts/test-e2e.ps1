# test-e2e.ps1 — corre el test de integración e2e sin Docker (modernc.org/sqlite)
Set-Location -Path (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location ..  # subir a backend/
go test ./internal/integration -run Test_FullAPI_HappyPath -v
