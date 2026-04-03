# INF2007 – TN1 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/coverage.yml/badge.svg) ![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

## Justification des cas de test choisis

L'analyse de la fonction `DaysUntilDeadline` révèle exactement quatre chemins d'exécution distincts, soit une erreur de parsing sur `currentDate`, une erreur de parsing sur `deadline`, une erreur métier lorsque la deadline précède la date actuelle, et enfin le cas de succès où le nombre de jours est calculé normalement. Cette structure a guidé l'ensemble de ma stratégie de test.

J'ai rédigé 24 tests unitaires dans `deadline_test.go`, organisés en six sections thématiques. La première section regroupe les cas positifs, notamment le calcul d'une échéance future, le cas où les deux dates sont identiques, ainsi que des plages plus larges d'un an et de dix ans. La deuxième section cible les cas limites liés aux transitions calendaires, comme le passage de fin de mois, la transition de fin d'année et le traitement d'une année bissextile. Les troisième et quatrième sections couvrent respectivement les erreurs de format sur `currentDate` et sur `deadline`, incluant des scénarios variés tels qu'un mois hors plage, un séparateur incorrect, une chaîne vide, des espaces parasites, un jour inexistant, un mois ou jour à zéro, un format ISO comportant une heure, l'absence de zéro initial et une entrée complètement aléatoire. La cinquième section traite de l'erreur métier lorsque la deadline est antérieure, et la sixième vérifie le comportement lorsque les deux entrées sont invalides simultanément.

Afin de pousser l'expérimentation et d'appliquer l'ensemble des concepts vus au chapitre 1, j'ai également produit des fichiers complémentaires contenant des tests tabulaires, du fuzz testing et des benchmarks.

<img width="670" height="70" alt="image" src="https://github.com/user-attachments/assets/0fc87747-37a3-4f12-ba89-3ecfd86f4445" />

## Comment les tests garantissent la correction de la fonction

Chacun des 24 tests vérifie systématiquement les deux sorties de la fonction, soit la valeur entière `days` et l'erreur `err`. Cette double vérification permet de s'assurer que le comportement est cohérent dans tous les scénarios.

Les tests positifs valident la justesse du calcul sur des plages de durées variées, allant de six jours jusqu'à 3653 jours sur une décennie complète. Les tests négatifs confirment que chaque catégorie d'entrée invalide déclenche l'erreur attendue tout en retournant zéro jour, garantissant ainsi un échec contrôlé et prévisible. Les cas limites, quant à eux, ciblent les transitions calendaires et les années bissextiles, des situations reconnues pour générer des anomalies subtiles dans le traitement des dates.

La couverture atteint 100 % des instructions, ce qui signifie que les quatre chemins d'exécution sont tous exercés. Cette couverture complète, combinée à la variété des cas testés, offre une validation rigoureuse de la correction de la fonction.

## Défis rencontrés

Le principal défi a été la gestion des erreurs de façon cohérente avec les attentes du devoir. En Go, comparer directement des erreurs créées avec `errors.New(...)` peut s'avérer ambigu lorsqu'une nouvelle instance est recréée dans le test. Pour rester aligné avec la fonction fournie, j'ai adopté une approche en trois étapes pour chaque test négatif, soit vérifier la présence d'une erreur avec `err != nil`, comparer son message via `err.Error()`, puis confirmer que la valeur retournée demeure à zéro.

J'ai également porté une attention particulière au choix entre `t.Fatalf` et `t.Errorf`. Le premier interrompt immédiatement l'exécution du test lorsqu'il serait trompeur de continuer, par exemple quand une erreur inattendue rend la vérification subséquente de `days` sans objet. Le second signale l'échec tout en laissant le test poursuivre, ce qui permet de collecter davantage d'informations utiles au débogage.

Un deuxième défi a été la sélection des cas limites pertinents. La fonction ne comportant que quatre chemins d'exécution, il était essentiel d'identifier les entrées qui exercent chacun d'entre eux de manière significative, sans introduire de redondance. J'ai regroupé les tests par catégorie, soit format de `currentDate`, format de `deadline` et erreur métier, afin de préserver la lisibilité malgré le nombre élevé de cas.


## Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Vidéo explicative : https://youtu.be/Tsw6rHtLz_k

## Fichiers utilisés

- Implémentation principale : [deadline.go](TN1/deadline.go)
- Tests unitaires (24 tests, organisés par section) : [deadline_test.go](TN1/deadline_test.go)
- Tests tabulaires (réplicat des 24 cas) : [deadline_table_test.go](TN1/deadline_table_test.go)
- Tests de fuzzing (entrées aléatoires) : [deadline_fuzz_test.go](TN1/deadline_fuzz_test.go)
- Benchmarks (6 scénarios de performance) : [deadline_benchmark_test.go](TN1/deadline_benchmark_test.go)
