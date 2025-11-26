# scripts/bump-version.ps1 - Bump version across all files

param(
    [Parameter(Mandatory=$true)]
    [string]$NewVersion
)

Write-Host "🔄 Bumping version to $NewVersion..." -ForegroundColor Green

# Update frontend/package.json
$packageJson = Get-Content "frontend/package.json" -Raw | ConvertFrom-Json
$packageJson.version = $NewVersion
$packageJson | ConvertTo-Json -Depth 10 | Set-Content "frontend/package.json"
Write-Host "✅ Updated frontend/package.json" -ForegroundColor Green

# Update wails.json
$wailsJson = Get-Content "wails.json" -Raw | ConvertFrom-Json
$wailsJson.version = $NewVersion
$wailsJson.info.productVersion = $NewVersion
$wailsJson | ConvertTo-Json -Depth 10 | Set-Content "wails.json"
Write-Host "✅ Updated wails.json" -ForegroundColor Green

# Update internal/version/version.go
$versionFile = Get-Content "internal/version/version.go"
$versionFile = $versionFile -replace 'const Version = "[^"]*"', "const Version = `"$NewVersion`""
$versionFile | Set-Content "internal/version/version.go"
Write-Host "✅ Updated internal/version/version.go" -ForegroundColor Green

# Update frontend/src/components/modals/settings/AboutTab.vue
$aboutTabFile = Get-Content "frontend/src/components/modals/settings/AboutTab.vue"
$aboutTabFile = $aboutTabFile -replace 'appVersion: "[^"]*"', "appVersion: `"$NewVersion`""
$aboutTabFile | Set-Content "frontend/src/components/modals/settings/AboutTab.vue"
Write-Host "✅ Updated AboutTab.vue" -ForegroundColor Green

# Update README.md badges
$readmeContent = Get-Content "README.md" -Raw
$readmeContent = $readmeContent -replace 'version-[0-9]+\.[0-9]+\.[0-9]+', "version-$NewVersion"
$readmeContent | Set-Content "README.md"
Write-Host "✅ Updated README.md version badge" -ForegroundColor Green

# Update README_zh.md badges
$readmeZhContent = Get-Content "README_zh.md" -Raw
$readmeZhContent = $readmeZhContent -replace 'version-[0-9]+\.[0-9]+\.[0-9]+', "version-$NewVersion"
$readmeZhContent | Set-Content "README_zh.md"
Write-Host "✅ Updated README_zh.md version badge" -ForegroundColor Green

Write-Host "🎉 Version bumped to $NewVersion!" -ForegroundColor Green
Write-Host "📝 Remember to update CHANGELOG.md" -ForegroundColor Yellow
