# INF2007 – TN1 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/coverage.yml/badge.svg) ![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

## Justification des cas de test choisis

J’ai choisi ces cas de test pour couvrir les comportements essentiels de `DaysUntilDeadline`, en situations normales et en erreur. Les tests positifs (échéance future et même jour) valident la fonctionnalité principale : calculer correctement les jours entre deux dates valides. Le cas futur reflète un usage réel, tandis que le cas « même jour » vérifie la limite où le résultat attendu est 0.

Les tests négatifs (format invalide, séparateur incorrect, échéance antérieure) vérifient la robustesse. Ils confirment que la fonction rejette les entrées mal formées et retourne une erreur explicite quand la logique métier est violée. Dans ces cas, la valeur retournée doit rester 0, ce qui garantit un comportement cohérent en échec.

J’ai aussi ajouté des cas limites (année bissextile, date impossible, fin de mois, fin d’année, chaîne vide, espaces) pour réduire le risque de bogues subtils liés aux dates.

<img width="975" height="94" alt="image" src="https://github.com/user-attachments/assets/6f0cd66e-6999-4351-a728-bb7ee81442ae" />

## Comment les tests garantissent la correction de la fonction

Les tests garantissent la correction en vérifiant systématiquement les sorties attendues sur des entrées valides et invalides. Les tests positifs confirment que la logique de calcul est correcte quand les dates sont bien formées : la fonction retourne le bon nombre de jours et aucune erreur. Cela valide le comportement nominal de `DaysUntilDeadline`.

Les tests négatifs valident la gestion des erreurs : pour des entrées invalides (format incorrect, séparateur invalide, échéance antérieure), la fonction doit retourner une erreur appropriée et la valeur 0. Cette double vérification (`days` et `err`) assure un échec contrôlé et prévisible.

Les cas limites (même jour, année bissextile, fin de mois, fin d’année, date vide, espaces) renforcent la correction en ciblant des situations où les erreurs de logique sont fréquentes. Enfin, l’exécution automatisée avec `go test` et la couverture mesuree (100 %) montrent que toutes les branches actuelles du code sont exercées, sans confondre couverture élevée et qualité des cas de test.

## Défis rencontrés

Le principal défi a été la gestion des erreurs de façon cohérente avec les attentes du devoir. En Go, comparer directement des erreurs créées avec `errors.New(...)` peut être ambigu si une nouvelle erreur est recréée dans le test.

Pour rester aligné avec la fonction fournie, j’ai vérifié la présence de l’erreur (`err != nil`), comparé son message (`err.Error()`) et confirmé que la valeur retournée reste 0 en cas d’échec.

Un deuxième défi a été le choix des cas limites pertinents. Le traitement des dates comporte des pièges, soit : année bissextile, 29 février invalide, fin de mois, fin d’année, chaînes vides ou avec espaces, séparateurs incorrects. Il fallait couvrir ces risques sans rendre la suite redondante.

J’ai aussi dû équilibrer couverture et pertinence : l’objectif n’est pas seulement d’atteindre un pourcentage élevé, mais d’avoir des tests lisibles, commentés et orientés sur des comportements métier clairs.


## Fichiers utilisés

- Implémentation principale : https://github.com/moyamelissa/Advanved-Programming/blob/main/TN1/deadline.go
- Tests unitaires : https://github.com/moyamelissa/Advanved-Programming/blob/main/TN1/deadline_test.go

## Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Vidéo explicative : https://youtu.be/Tsw6rHtLz_k

## Fichiers utilisés

- Implémentation principale : [deadline.go](TN1/deadline.go)
- Tests unitaires : [deadline_test.go](TN1/deadline_test.go)
