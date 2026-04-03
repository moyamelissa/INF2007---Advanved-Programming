# INF2007 – TN1 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/coverage.yml/badge.svg) ![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

## Justification des cas de test choisis

La fonction `DaysUntilDeadline` possède exactement 4 chemins d'exécution : erreur de parsing de `currentDate`, erreur de parsing de `deadline`, erreur métier (deadline antérieure) et succès. J'ai écrit 24 tests unitaires dans `deadline_test.go` pour couvrir ces 4 chemins de façon exhaustive, organisés en sections : cas positifs (échéance future, même jour, 1 an, 10 ans), cas limites de transition (fin de mois, fin d'année, année bissextile), erreurs de format sur `currentDate` (mois invalide, séparateur `/`, vide, espaces, jour hors plage, mois/jour zéro, heure ISO, sans zéro initial, garbage), erreurs de format sur `deadline` (mois invalide, vide, espaces, jour hors plage), erreur métier (deadline antérieure) et les deux entrées invalides simultanément.

En complément, j'ai poussé l'expérimentation en appliquant tous les concepts du chapitre 1 : tests tabulaires, fuzz testing et benchmarks, dans des fichiers séparés.

<img width="975" height="94" alt="image" src="https://github.com/user-attachments/assets/6f0cd66e-6999-4351-a728-bb7ee81442ae" />

## Comment les tests garantissent la correction de la fonction

Les 24 tests vérifient systématiquement les deux sorties (`days` et `err`) sur chacun des 4 chemins d'exécution. Les tests positifs confirment le calcul correct sur des plages variées (6 jours, 0 jour, 365 jours, 3653 jours). Les tests négatifs valident que chaque type d'entrée invalide retourne l'erreur appropriée et `days == 0`, assurant un échec contrôlé. Les cas limites (transitions de mois/année, bissextile) ciblent les situations où les erreurs de logique sont fréquentes.

Avec 100 % de couverture et les 4 chemins d'exécution tous exercés, la correction est validée de façon rigoureuse.

## Défis rencontrés

Le principal défi a été la gestion des erreurs de façon cohérente avec les attentes du devoir. En Go, comparer directement des erreurs créées avec `errors.New(...)` peut être ambigu si une nouvelle erreur est recréée dans le test.

Pour rester aligné avec la fonction fournie, j'ai vérifié la présence de l'erreur (`err != nil`), comparé son message (`err.Error()`) et confirmé que la valeur retournée reste 0 en cas d'échec. J'ai aussi utilisé `t.Fatalf` et `t.Errorf` de façon intentionnelle : `t.Fatalf` arrête immédiatement le test quand continuer n'a pas de sens (par exemple, si une erreur inattendue survient et que vérifier `days` après serait trompeur), tandis que `t.Errorf` signale l'échec mais laisse le test continuer pour collecter toutes les informations utiles au débogage.

Un deuxième défi a été le choix des cas limites pertinents. La fonction n'ayant que 4 chemins d'exécution, il fallait identifier les entrées qui exercent chacun de façon significative (dates hors plage, mois/jour zéro, format ISO avec heure, garbage) sans rendre la suite redondante. J'ai regroupé les tests par catégorie (format currentDate, format deadline, erreur métier) pour maintenir la lisibilité malgré les 24 cas.

Enfin, pour pousser l'expérimentation, j'ai aussi créé des tests tabulaires, de fuzzing et des benchmarks, ce qui m'a permis d'appliquer l'ensemble des concepts vus au chapitre 1.


## Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Vidéo explicative : https://youtu.be/Tsw6rHtLz_k

## Fichiers utilisés

- Implémentation principale : [deadline.go](TN1/deadline.go)
- Tests unitaires (24 tests, organisés par section) : [deadline_test.go](TN1/deadline_test.go)
- Tests tabulaires (réplicat des 24 cas) : [deadline_table_test.go](TN1/deadline_table_test.go)
- Tests de fuzzing (entrées aléatoires) : [deadline_fuzz_test.go](TN1/deadline_fuzz_test.go)
- Benchmarks (6 scénarios de performance) : [deadline_benchmark_test.go](TN1/deadline_benchmark_test.go)
