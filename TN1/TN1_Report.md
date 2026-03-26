# INF2007 - TN1 - Melissa Moya

## Justification des cas de test choisis

J'ai choisi ces cas de test pour couvrir les comportements essentiels de `DaysUntilDeadline`, en situations normales et en erreur.

Les tests positifs (echeance future et meme jour) valident la fonctionnalite principale : calculer correctement les jours entre deux dates valides. Le cas futur reflete un usage reel, tandis que le cas "meme jour" verifie la limite ou le resultat attendu est 0.

Les tests negatifs (format invalide, separateur incorrect, echeance anterieure) verifient la robustesse. Ils confirment que la fonction rejette les entrees mal formees et retourne une erreur explicite quand la logique metier est violee. Dans ces cas, la valeur retournee doit rester 0, ce qui garantit un comportement coherent en echec.

J'ai aussi ajoute des cas limites (annee bissextile, date impossible, fin de mois, fin d'annee, chaine vide, espaces) pour reduire le risque de bogues subtils lies aux dates.

## Comment les tests garantissent la correction de la fonction

Les tests garantissent la correction en verifiant systematiquement les sorties attendues sur des entrees valides et invalides.

Les tests positifs confirment que la logique de calcul est correcte quand les dates sont bien formees : la fonction retourne le bon nombre de jours et aucune erreur. Cela valide le comportement nominal de `DaysUntilDeadline`.

Les tests negatifs valident la gestion des erreurs : pour des entrees invalides (format incorrect, separateur invalide, echeance anterieure), la fonction doit retourner une erreur appropriee et la valeur 0. Cette double verification (`days` et `err`) assure un echec controle et previsible.

Les cas limites (meme jour, annee bissextile, fin de mois, fin d'annee, date vide, espaces) renforcent la correction en ciblant des situations ou les erreurs de logique sont frequentes.

Enfin, l'execution automatisee avec `go test` et la couverture mesuree (100 %) montrent que toutes les branches actuelles du code sont exercees, sans confondre couverture elevee et qualite des cas de test.

## Defis rencontres

Le principal defi a ete la gestion des erreurs de facon coherente avec les attentes du devoir. En Go, comparer directement des erreurs creees avec `errors.New(...)` peut etre ambigu si une nouvelle erreur est recreee dans le test.

Pour rester aligne avec la fonction fournie, j'ai verifie la presence de l'erreur (`err != nil`), compare son message (`err.Error()`) et confirme que la valeur retournee reste 0 en cas d'echec.

Un deuxieme defi a ete le choix des cas limites pertinents. Le traitement des dates comporte des pieges, soit : annee bissextile, 29 fevrier invalide, fin de mois, fin d'annee, chaines vides ou avec espaces, separateurs incorrects. Il fallait couvrir ces risques sans rendre la suite redondante.

J'ai aussi du equilibrer couverture et pertinence : l'objectif n'est pas seulement d'atteindre un pourcentage eleve, mais d'avoir des tests lisibles, commentes et orientes sur des comportements metier clairs.

## Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Video Explicative : https://youtu.be/Tsw6rHtLz_k
