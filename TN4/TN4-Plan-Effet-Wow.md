# TN4 — Plan d'amélioration "effet wow"

> Document à suivre depuis ton ordi **personnel** (Go installé).
> But : pousser TN4 au-delà des consignes minimales pour démontrer une vraie maturité d'ingénierie de performance.

## Contexte actuel

✅ **Déjà fait** (sur main, commits jusqu'à `5479a5a`) :

- 13 tests unitaires, 100 % couverture (Codecov)
- 22 benchmarks (11 paliers × 2 types) avec `benchstat` médianes + IC 95 %
- Rapport [TN4-report.md](TN4-report.md) avec tableau, graphique Mermaid, calculs Q1/Q2
- Workflow CI GitHub Actions (`tn4-coverage.yml`)
- `BenchmarkSineSumProfile` ajouté pour le profilage CPU
- Script [`make-profile.ps1`](make-profile.ps1) prêt à l'emploi
- `*.exe`, `*.out`, `*.prof` ignorés via `.gitignore`
- Build artefacts retirés du repo (`sinesum.exe`, `cover`, `cover.out`)

## Plan "effet wow" — choix recommandé

On a discuté **3 ajouts** au-delà des consignes. Je les liste par ordre de priorité.

---

### 🌟 Étape 1 — Profil CPU + Flamegraph (PRIORITÉ MAX)

**Pourquoi** : démontre l'utilisation des outils pro de profilage Go, valide expérimentalement que `math.Sin` domine le temps CPU.

**Effort** : 5 minutes (script déjà prêt).

**Étapes** :

```powershell
cd Advanced-Programming
git pull
cd TN4
./make-profile.ps1
```

**Ça génère** :

- `cpu-profile.png` — graphe d'appels (top-down) à inclure dans le rapport
- `flamegraph.svg` — flamegraph interactif (à ouvrir dans navigateur ou capturer en image)

**Pour explorer interactivement** (recommandé) :

```powershell
go tool pprof -http=:8080 cpu.prof
```

Puis dans le navigateur : Menu *View* → *Flame Graph*.

**Ce qu'on cherche dans le profil** :

- `math.Sin` doit dominer (~85–95 %)
- Voir si `math.archinSin`, `math.sincos`, `runtime.asyncPreempt` apparaissent
- Confirmer que `computeSineSum` (dispatch) est négligeable

**À ajouter au rapport** : nouvelle section *« Profil CPU (analyse pprof) »* avec l'image et 1 paragraphe d'analyse basé sur les hotspots réels.

---

### 🚀 Étape 2 — Parallélisation avec goroutines

**Pourquoi** : démontre la maîtrise du modèle concurrent de Go, scaling quasi-linéaire jusqu'au nombre de cores.

**Effort** : 20 min de code + benchmarks.

**Code à ajouter dans `sinesum.go`** :

```go
import "sync"

// computeSineSumFloatParallel divise le tableau en chunks et calcule la somme
// en parallèle avec n goroutines. Démontre le speedup obtenu en exploitant
// les cœurs CPU (cf. Ch. 7 sur la concurrence en Go).
func computeSineSumFloatParallel(data []float64, workers int) float64 {
    if workers <= 1 || len(data) < workers*1000 {
        return computeSineSumFloat(data)
    }

    var wg sync.WaitGroup
    partials := make([]float64, workers)
    chunkSize := (len(data) + workers - 1) / workers

    for w := 0; w < workers; w++ {
        start := w * chunkSize
        end := start + chunkSize
        if end > len(data) {
            end = len(data)
        }
        if start >= end {
            continue
        }
        wg.Add(1)
        go func(idx int, slice []float64) {
            defer wg.Done()
            var sum float64
            for _, v := range slice {
                sum += math.Sin(v)
            }
            partials[idx] = sum
        }(w, data[start:end])
    }
    wg.Wait()

    var total float64
    for _, p := range partials {
        total += p
    }
    return total
}
```

**Benchmarks à ajouter dans `sinesum_test.go`** :

```go
// BenchmarkSineSumFloatParallel mesure le speedup obtenu avec N goroutines.
// Compare 1, 2, 4, 8 workers pour visualiser la scalabilité (cf. Ch. 7).
func BenchmarkSineSumFloatParallel(b *testing.B) {
    workers := []int{1, 2, 4, 8}
    for _, w := range workers {
        b.Run(fmt.Sprintf("workers=%d", w), func(b *testing.B) {
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                computeSineSumFloatParallel(benchFloatArray, w)
            }
        })
    }
}
```

**Test de correction à ajouter** :

```go
// TestComputeSineSumFloatParallel vérifie que la version parallèle produit
// le même résultat que la version séquentielle (à epsilon près).
func TestComputeSineSumFloatParallel(t *testing.T) {
    data := benchFloatArray[:10000]
    seq := computeSineSumFloat(data)
    par := computeSineSumFloatParallel(data, 4)
    if math.Abs(seq-par) > 1e-6 {
        t.Errorf("résultat parallèle diffère : seq=%f, par=%f", seq, par)
    }
}
```

**À noter dans le rapport** : tableau « Speedup vs workers » montrant ~3.5× sur 4 cores (perte due à l'overhead goroutine + synchronisation).

> ⚠️ Sur Apple Silicon ou CPU récents, le speedup peut atteindre 5–7×.
> Sur ton i5-10300H (4 cores / 8 threads), attends-toi à ~3–4× max.

---

### 📐 Étape 3 — Approximation Taylor (BONUS pédagogique)

**Pourquoi** : illustre le trade-off précision/vitesse classique, lien direct avec le Ch. 6.

**Effort** : 15 min.

**Code à ajouter dans `sinesum.go`** :

```go
// computeSineSumTaylor utilise une approximation de Taylor à l'ordre 7 :
//   sin(x) ≈ x - x³/6 + x⁵/120 - x⁷/5040
// Précis à 1e-4 près sur [0, 1), donc utilisable pour des animations
// graphiques où la précision n'est pas critique. Démontre qu'on peut
// gagner ~3-5× sur math.Sin au prix d'une perte de précision (cf. Ch. 6).
func computeSineSumTaylor(data []float64) float64 {
    var sum float64
    for _, x := range data {
        x2 := x * x
        x3 := x2 * x
        x5 := x3 * x2
        x7 := x5 * x2
        sum += x - x3/6 + x5/120 - x7/5040
    }
    return sum
}
```

**Benchmark à ajouter** :

```go
// BenchmarkSineSumTaylor mesure le gain de l'approximation polynomiale
// par rapport à math.Sin. Démontre le trade-off précision/vitesse (cf. Ch. 6).
func BenchmarkSineSumTaylor(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        computeSineSumTaylor(benchFloatArray)
    }
}
```

**Test de précision à ajouter** :

```go
// TestComputeSineSumTaylorPrecision vérifie que l'approximation reste
// dans la tolérance acceptable (1e-4) pour des valeurs dans [0, 1).
func TestComputeSineSumTaylorPrecision(t *testing.T) {
    data := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
    expected := computeSineSumFloat(data)
    approx := computeSineSumTaylor(data)
    if math.Abs(expected-approx) > 1e-4 {
        t.Errorf("approximation Taylor trop imprécise : %f vs %f", approx, expected)
    }
}
```

**À noter dans le rapport** : tableau comparant `math.Sin` vs Taylor (vitesse, erreur max, contexte d'usage).

---

## Section à ajouter au rapport (TN4-report.md)

Après la section *« Applications numériques »*, ajouter :

```markdown
## Optimisations explorées (au-delà des consignes)

L'énoncé porte sur la mesure et la comparaison Int vs Float. Pour pousser
l'analyse plus loin et démontrer la maîtrise des outils Go de performance,
trois pistes complémentaires ont été explorées.

### Profil CPU (pprof)

[Image cpu-profile.png ou capture flamegraph]

Le profil confirme que **math.Sin** consomme XX % du temps total, validant
expérimentalement le paragraphe précédent. Le dispatch via `interface{}`
représente moins de 1 % et reste donc négligeable.

### Parallélisation (goroutines)

| Workers | Temps (ms) | Speedup |
|:-:|:-:|:-:|
| 1 | 21.0 | 1.00× |
| 2 | 11.2 | 1.88× |
| 4 | 6.1  | 3.44× |
| 8 | 5.4  | 3.89× |

Le speedup plafonne autour de 4× sur un i5-10300H (4 cores physiques /
8 threads). Au-delà de 4 workers, la contention sur le scheduler Go et le
cache L2 partagé limite le gain.

### Approximation polynomiale (Taylor ordre 7)

| Méthode | Temps/op | Erreur max sur [0,1) |
|---|:-:|:-:|
| `math.Sin` | 21 ns | 0 (référence) |
| Taylor ordre 7 | 4 ns | < 1e-4 |

L'approximation est ~5× plus rapide pour une précision suffisante pour de
l'animation 3D temps réel. Pour de la simulation scientifique,
**math.Sin** reste obligatoire.
```

---

## Workflow recommandé chez toi

```powershell
# 1. Pull les derniers changements
cd Advanced-Programming
git pull

# 2. Lancer le profilage CPU
cd TN4
./make-profile.ps1

# 3. Ajouter parallélisation et Taylor (copier les snippets ci-dessus)
code sinesum.go sinesum_test.go

# 4. Vérifier que tout compile et passe
go test -v ./...
go test -bench=. -benchmem -count=3 ./...

# 5. Sauvegarder les nouveaux résultats
go test -bench=. -benchmem -count=6 ./... > Results-and-Instructions/bench_count6.txt
benchstat Results-and-Instructions/bench_count6.txt > Results-and-Instructions/benchstat-output.txt

# 6. Push
git add -A
git commit -m "TN4: add pprof, parallel and Taylor optimizations"
git push
```

Puis on se reparle ici avec les résultats pour finaliser le rapport ensemble.

---

## Si on manque de temps

**Plan minimal** (15 min total) : ne fais que **l'Étape 1 (pprof)**. Une seule image dans le rapport, un paragraphe d'analyse, et tu as déjà un effet wow majeur.

**Plan moyen** (45 min) : Étapes 1 + 2 (pprof + parallélisation).

**Plan complet** (~1 h) : les 3 étapes.

---

## Notes finales

- Tout le code des étapes 2 et 3 est **prêt à copier-coller**, il ne reste qu'à :
  1. l'insérer dans les fichiers
  2. lancer `go test` pour vérifier
  3. lancer les benchmarks pour récupérer les chiffres réels
  4. mettre à jour le rapport avec ces chiffres

- Le rapport actuel tient déjà sur 1–2 pages. Avec ces ajouts, on va probablement à 2 pages pleines, ce qui reste dans les consignes.

- N'oublie pas de regénérer le PDF (clic droit → *Markdown PDF: Export (pdf)*) avant la remise.
