# ==============================
# Password-Danie - Test Suite PS
# ==============================

$ErrorActionPreference = 'Stop'
$BASE = "http://localhost:8080"

# variables de ámbito de script
$script:Headers = $null
$script:itemId  = $null

function It([string]$name, [scriptblock]$test) {
  try {
    & $test
    Write-Host "PASS - $name" -ForegroundColor Green
  } catch {
    Write-Host "FAIL - $name" -ForegroundColor Red
    if ($_.Exception.Response -and $_.Exception.Response.Content) {
      try {
        $body = $_.Exception.Response.Content | ConvertFrom-Json
        Write-Host ("Response: " + ($body | ConvertTo-Json -Compress)) -ForegroundColor DarkGray
      } catch {}
    }
    throw
  }
}

function InvokeJson([string]$method,[string]$url,$headers=$null,$bodyObj=$null) {
  $params = @{ Method = $method; Uri = $url }
  if ($headers) { $params['Headers'] = $headers }
  if ($bodyObj -ne $null) {
    $params['ContentType'] = 'application/json'
    $params['Body']        = ($bodyObj | ConvertTo-Json -Depth 8)
  }
  return Invoke-RestMethod @params
}

function ExpectStatus([scriptblock]$call, [int[]]$expected) {
  try {
    $res = & $call
    if (-not ($expected -contains 200) -and -not ($expected -contains 201) -and -not ($expected -contains 204)) {
      throw "Esperaba status en $expected pero fue 2xx"
    }
    return $res
  } catch {
    $we = $_.Exception
    if ($we.Response -and $we.Response.StatusCode.value__) {
      $code = $we.Response.StatusCode.value__
      if ($expected -contains $code) { return $null }
      throw "HTTP $code no está en esperados: $($expected -join ', ')"
    }
    throw
  }
}

# ---------- 1) Salud ----------
It "Health devuelve ok" {
  $h = InvokeJson GET "$BASE/health"
  if ($h.status -ne 'ok') { throw "status=$($h.status)" }
}

# ---------- 2) Auth ----------
$email = "dani+" + ([guid]::NewGuid().ToString("N").Substring(0,6)) + "@example.com"
$pwd   = "secret123"

It "Register 201 usuario nuevo" {
  ExpectStatus { InvokeJson POST "$BASE/auth/register" $null @{ email=$email; password=$pwd } } 201 | Out-Null
}
It "Register 409 duplicado" {
  ExpectStatus { InvokeJson POST "$BASE/auth/register" $null @{ email=$email; password=$pwd } } 409 | Out-Null
}

It "Login 200 ok devuelve token" {
  $resp = ExpectStatus { InvokeJson POST "$BASE/auth/login" $null @{ email=$email; password=$pwd } } 200
  if (-not $resp.token) { throw "no token" }
  $script:Headers = @{ Authorization = "Bearer " + $resp.token }
}
It "Login 401 password mala" {
  ExpectStatus { InvokeJson POST "$BASE/auth/login" $null @{ email=$email; password="xx" } } 401 | Out-Null
}

# ---------- 3) Protección ----------
It "401 sin token" {
  ExpectStatus { InvokeJson GET "$BASE/api/vault" $null $null } 401 | Out-Null
}
It "401 token inválido" {
  $h = @{ Authorization = "Bearer abc.def.ghi" }
  ExpectStatus { InvokeJson GET "$BASE/api/vault" $h $null } 401 | Out-Null
}

# ---------- 4) Vault CRUD + filtros ----------
It "400 crear sin title" {
  ExpectStatus {
    InvokeJson POST "$BASE/api/vault" $script:Headers @{ username="dani"; password="p"; url="https://github.com" }
  } 400 | Out-Null
}
It "201 crear item válido" {
  $res = ExpectStatus {
    InvokeJson POST "$BASE/api/vault" $script:Headers @{
      title="GitHub"; username="dani"; password="super-secret"
      url="https://github.com"; notes="Repo personal"
      icon="https://github.githubassets.com/favicons/favicon.png"
    }
  } 201
  if (-not $res.id) { throw "no id" }
  $script:itemId = [int]$res.id
  if ($res.password -or $res.notes) { throw "no deberían exponerse password/notes en list item" }
}
It "200 listar items sin campos sensibles" {
  $list = ExpectStatus { InvokeJson GET "$BASE/api/vault" $script:Headers $null } 200
  if (@($list).Count -lt 1) { throw "lista vacía" }
  foreach($i in @($list)) {
    if ($i.password -or $i.notes) { throw "expuestos password/notes en list" }
  }
}

# ⬇️ Cambiado: el detalle NO debe traer password; solo metadatos + notes descifradas.
It "200 get detalle (sin password) con notes descifradas" {
  $det = ExpectStatus { InvokeJson GET "$BASE/api/vault/$($script:itemId)" $script:Headers $null } 200


  if ($det.PSObject.Properties.Name -contains 'password' -and $det.password) {
    throw "el detalle no debe incluir password"
  }


  if ($det.PSObject.Properties.Name -contains 'notes') {
    if (-not ($det.notes -is [string]) -or ([string]::IsNullOrWhiteSpace($det.notes))) {
      throw "notes descifrado vacío o no es string"
    }
  }
}


It "200 patch actualiza notes" {
  $esperado = "Repo personal (actualizado)"
  $patch = @{ notes = $esperado }
  $null  = ExpectStatus { InvokeJson PATCH "$BASE/api/vault/$($script:itemId)" $script:Headers $patch } 200

  $det   = InvokeJson GET "$BASE/api/vault/$($script:itemId)" $script:Headers
  if ($det.PSObject.Properties.Name -contains 'notes') {
    if ($det.notes -ne $esperado) { throw "no actualizó notes" }
  }
}

It "GET /api/vault/:id/reveal devuelve solo password" {
  $reveal = ExpectStatus { InvokeJson GET "$BASE/api/vault/$($script:itemId)/reveal" $script:Headers } 200
  if ($reveal.password -ne "super-secret") { throw "password incorrecto en reveal" }
  if ($reveal.PSObject.Properties.Name -contains 'notes') { throw "reveal no debe incluir notes" }
}
It "Filtros q y domain" {
  $q = InvokeJson GET "$BASE/api/vault?q=git" $script:Headers
  if (@($q).Count -lt 1) { throw "q=git no devolvió resultados" }
  $d = InvokeJson GET "$BASE/api/vault?domain=github.com" $script:Headers
  if (@($d).Count -lt 1) { throw "domain=github.com no devolvió resultados" }
}
It "Paginación limit/offset y sort" {
  1..2 | ForEach-Object {
    InvokeJson POST "$BASE/api/vault" $script:Headers @{
      title="Extra$_"; username="u$_"; password="p$_"; url="https://site$_.com"; notes="n$_"
    } | Out-Null
  }

  $all = InvokeJson GET "$BASE/api/vault" $script:Headers
  $p1  = InvokeJson GET "$BASE/api/vault?limit=2&offset=0" $script:Headers
  $p2  = InvokeJson GET "$BASE/api/vault?limit=2&offset=2" $script:Headers

  if (@($p1).Count -ne 2) { throw "p1 debe tener 2" }

  $expectedP2 = [math]::Min(2, [math]::Max(0, @($all).Count - 2))
  if (@($p2).Count -ne $expectedP2) { throw "p2 esperado=$expectedP2 pero fue $(@($p2).Count)" }

  $sorted = InvokeJson GET "$BASE/api/vault?sort=updated_desc" $script:Headers
  if (@($sorted).Count -lt 1) { throw "sort vacío" }
}


It "200 delete y 404 al volver a pedir" {
  ExpectStatus { InvokeJson DELETE "$BASE/api/vault/$($script:itemId)" $script:Headers $null } 200 | Out-Null
  ExpectStatus { InvokeJson GET    "$BASE/api/vault/$($script:itemId)" $script:Headers $null } 404 | Out-Null
}

# ---------- 5) Aislamiento ----------
It "Aislamiento entre usuarios" {
  $email2 = "dani+" + ([guid]::NewGuid().ToString("N").Substring(0,6)) + "@example.com"
  InvokeJson POST "$BASE/auth/register" $null @{ email=$email2; password=$pwd } | Out-Null
  $r2 = InvokeJson POST "$BASE/auth/login" $null @{ email=$email2; password=$pwd }
  $H2 = @{ Authorization = "Bearer " + $r2.token }
  $list2 = InvokeJson GET "$BASE/api/vault" $H2
  if (@($list2).Count -gt 0) { throw "user2 ve items del user1" }
  ExpectStatus { InvokeJson GET "$BASE/api/vault/999999" $H2 $null } 404 | Out-Null
}

# ---------- 6) Reset de contraseña ----------
It "Forgot devuelve reset_token y Reset funciona" {
  $fg = InvokeJson POST "$BASE/auth/forgot" $null @{ email=$email }
  if (-not $fg.reset_token) { throw "no reset_token (en dev debe venir en JSON)" }
  ExpectStatus { InvokeJson POST "$BASE/auth/reset" $null @{ token=$fg.reset_token; new_password="nueva12345" } } 200 | Out-Null
  ExpectStatus { InvokeJson POST "$BASE/auth/reset" $null @{ token=$fg.reset_token; new_password="otra123"     } } 400 | Out-Null
  ExpectStatus { InvokeJson POST "$BASE/auth/login" $null @{ email=$email; password="nueva12345" } } 200 | Out-Null
  ExpectStatus { InvokeJson POST "$BASE/auth/login" $null @{ email=$email; password=$pwd } } 401 | Out-Null
}

# ---------- 7) CORS ----------
It "Preflight CORS OPTIONS /api/vault → 204" {
  $h = @{
    "Origin" = "http://localhost:3000"
    "Access-Control-Request-Method" = "POST"
  }
  $resp = Invoke-WebRequest -Method OPTIONS "$BASE/api/vault" -Headers $h
  if ($resp.StatusCode -ne 204) { throw "status=$($resp.StatusCode)" }
  if (-not $resp.Headers["Access-Control-Allow-Origin"])  { throw "falta ACAO" }
  if (-not $resp.Headers["Access-Control-Allow-Methods"]) { throw "falta ACAM" }
  if (-not $resp.Headers["Access-Control-Allow-Headers"]) { throw "falta ACAH" }
}

Write-Host "`n--- FIN DE PRUEBAS ---`n" -ForegroundColor Cyan
