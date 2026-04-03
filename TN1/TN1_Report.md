# INF2007 – TN1 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/coverage.yml/badge.svg) ![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

## Justification des cas de test choisis

L'analyse de la fonction `DaysUntilDeadline` révèle quatre chemins d'exécution distincts, soit une erreur de parsing sur `currentDate`, une erreur de parsing sur `deadline`, une erreur métier lorsque la deadline précède la date actuelle, et le cas de succès. Cette structure a guidé ma stratégie de test.

J'ai rédigé 24 tests unitaires dans `deadline_test.go`, organisés en six sections thématiques. La première regroupe les cas positifs, notamment une échéance future, deux dates identiques et des plages d'un an et dix ans. La deuxième cible les transitions calendaires telles que la fin de mois, la fin d'année et les années bissextiles. Les troisième et quatrième couvrent les erreurs de format sur `currentDate` et `deadline` respectivement, incluant un mois hors plage, un séparateur incorrect, une chaîne vide, des espaces parasites, un jour inexistant, un mois ou jour à zéro, un format ISO avec heure, l'absence de zéro initial et une entrée aléatoire. La cinquième traite de l'erreur métier et la sixième vérifie le comportement lorsque les deux entrées sont invalides simultanément.

## Comment les tests garantissent la correction de la fonction

Chacun des 24 tests vérifie systématiquement les deux sorties de la fonction, soit `days` et `err`, assurant un comportement cohérent dans tous les scénarios.

Les tests positifs valident la justesse du calcul sur des plages allant de six jours à 3653 jours. Les tests négatifs confirment que chaque catégorie d'entrée invalide déclenche l'erreur attendue tout en retournant zéro jour. Les cas limites ciblent les transitions calendaires et les années bissextiles, des situations reconnues pour générer des anomalies subtiles.

La couverture atteint 100 % des instructions, ce qui signifie que les quatre chemins d'exécution sont tous exercés. Cette couverture complète, combinée à la variété des cas testés, offre une validation rigoureuse de la correction de la fonction.
<img width="1036" height="90" alt="image" src="https://github.com/user-attachments/assets/cbe740e4-b55e-4b47-954e-84371cb1e6fd" />

## Défis rencontrés

Le principal défi a été la gestion des erreurs. En Go, comparer des erreurs créées avec `errors.New(...)` peut s'avérer ambigu lorsqu'une nouvelle instance est recréée dans le test. J'ai donc adopté une approche en trois étapes pour chaque test négatif, soit vérifier la présence d'une erreur, comparer son message via `err.Error()`, puis confirmer que `days` demeure à zéro.

J'ai également porté attention au choix entre `t.Fatalf` et `t.Errorf`. Le premier interrompt l'exécution du test lorsqu'il serait trompeur de continuer, par exemple quand une erreur inattendue rend la vérification de `days` sans objet. Le second signale l'échec tout en laissant le test poursuivre afin de collecter davantage d'informations.

Enfin, la sélection des cas limites pertinents a demandé d'identifier les entrées qui exercent chacun des quatre chemins de manière significative, sans redondance. Le regroupement par catégorie préserve la lisibilité malgré le nombre élevé de cas. Afin de pousser l'expérimentation et d'appliquer l'ensemble des concepts vus au chapitre 1, j'ai également produit des fichiers complémentaires contenant des tests tabulaires, du fuzz testing et des benchmarks.

## Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Vidéo explicative : https://youtu.be/Tsw6rHtLz_k

## Fichiers utilisés

- Implémentation principale : [deadline.go](TN1/deadline.go)
- Tests unitaires (24 tests, organisés par section) : [deadline_test.go](TN1/deadline_test.go)
- Tests tabulaires (réplicat des 24 cas) : [deadline_table_test.go](TN1/deadline_table_test.go)
- Tests de fuzzing (entrées aléatoires) : [deadline_fuzz_test.go](TN1/deadline_fuzz_test.go)
- Benchmarks (6 scénarios de performance) : [deadline_benchmark_test.go](TN1/deadline_benchmark_test.go)
