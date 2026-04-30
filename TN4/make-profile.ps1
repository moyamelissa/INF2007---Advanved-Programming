# make-profile.ps1
# Génère le profil CPU et le flamegraph pour le rapport TN4.
# Prérequis : Go installé et accessible via PATH.
#
# Usage (depuis le dossier TN4) :
#   ./make-profile.ps1
#
# Sortie :
#   - cpu.prof              : profil binaire pprof
#   - cpu-profile.png       : graphe d'appels (top-down)
#   - flamegraph.svg        : flamegraph interactif (ouvre dans navigateur)

$ErrorActionPreference = "Stop"

Write-Host "==> Compilation des tests..." -ForegroundColor Cyan
go test -c -o sinesum.test.exe

Write-Host "==> Lancement du benchmark de profilage (10s)..." -ForegroundColor Cyan
./sinesum.test.exe -test.bench=BenchmarkSineSumProfile -test.run=^$ -test.cpuprofile=cpu.prof -test.benchtime=10s

Write-Host "==> Génération du graphe d'appels PNG..." -ForegroundColor Cyan
go tool pprof -png -output=cpu-profile.png ./sinesum.test.exe cpu.prof

Write-Host "==> Génération du flamegraph SVG..." -ForegroundColor Cyan
# -web crée un SVG dans %TEMP% et l'ouvre. On force la sortie locale via -output.
go tool pprof -svg -output=flamegraph.svg ./sinesum.test.exe cpu.prof

# Cleanup binaire intermédiaire
Remove-Item sinesum.test.exe -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "✅ Terminé !" -ForegroundColor Green
Write-Host "  - cpu.prof           : profil brut"
Write-Host "  - cpu-profile.png    : graphe d'appels (à inclure dans le rapport)"
Write-Host "  - flamegraph.svg     : flamegraph (ouvrir dans navigateur)"
Write-Host ""
Write-Host "Pour explorer interactivement :" -ForegroundColor Yellow
Write-Host "  go tool pprof -http=:8080 cpu.prof"
Write-Host "  -> Menu View -> Flame Graph"
