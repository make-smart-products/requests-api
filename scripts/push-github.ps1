$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot\..

Write-Host "==> Git status"
git status -sb

Write-Host "==> Stage and commit"
git add -A
git commit -m "feat: Go backend for client cabinet with auth, applications and notifications" 2>$null

Write-Host "==> Configure remote"
$remote = "https://github.com/make-smart-products/requests-api.git"
if (git remote | Select-String -Pattern "^origin$" -Quiet) {
    git remote set-url origin $remote
} else {
    git remote add origin $remote
}

Write-Host "==> Push to GitHub"
git branch -M main
git push -u origin main

Write-Host "==> Done: https://github.com/make-smart-products/requests-api"
