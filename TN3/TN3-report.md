# INF2007 – TN3 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/tn3-coverage.yml/badge.svg) ![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

## Approche pour extraire et valider les bits

La structure binaire de chaque entrée 32 bits se divise en trois zones. Les bits 0 à 6 contiennent l'identifiant du capteur, le bit 7 sert de bit de validation et les bits 8 à 31 représentent la valeur mesurée. Pour extraire l'identifiant, j'ai construit un masque avec `(1<<7)-1 = 0x7F` qui isole les 7 bits de poids faible via un AND. Le bit de validation est testé par `entry & (1<<7)`, une opération qui cible un seul bit sans affecter les autres. Pour la valeur, le décalage `entry >> 8` ramène les bits 8 à 31 en positions 0 à 23, ce qui équivaut à une division par 2⁸. La détection de plusieurs bits actifs exploite l'identité `x & (x-1)` qui efface le bit le plus bas à 1. Si le résultat est non nul, il existait au moins deux bits actifs, ce qui constitue une violation de la spécification. Cette approche remplace un appel à `bits.OnesCount32` par une seule soustraction suivie d'un AND. Le masque des bits 8 à 31 est construit avec l'opérateur AND NOT (`&^`) appliqué à `^uint32(0)`, ce qui efface explicitement les bits 0 à 7 plutôt que de coder la valeur `0xFFFFFF00` en dur. Enfin, la position du bit actif est obtenue par `bits.TrailingZeros32` qui se traduit par l'instruction CPU `TZCNT` en un seul cycle, remplaçant une boucle de 24 itérations par une opération O(1).

## Défis rencontrés

Le premier défi a été de déterminer si la validation devait être globale ou locale. L'exemple de l'énoncé montre que `0x00000080` (ID=0, bit 7=1) déclenche une erreur alors qu'on cherche le capteur 5. Cela m'a confirmé que toute entrée invalide doit provoquer un rejet immédiat, peu importe le capteur recherché. Le second défi concerne le type `uint8` du paramètre `capteur`. Sa plage va de 0 à 255, mais les identifiants valides n'utilisent que 7 bits, donc 0 à 127. Sans vérification explicite, un appel avec `capteur=200` serait accepté silencieusement. Aucune entrée ne correspondrait et la fonction retournerait un tableau de zéros sans erreur, un comportement techniquement correct, mais sémantiquement trompeur. La connaissance de la plage des entiers non signés m'a permis d'identifier ce piège. Enfin, le cas `valeur = 0` (aucune mesure) interagit avec l'identité `x & (x-1)` de façon subtile. `uint32(0) - 1` produit `0xFFFFFFFF` par wrap-around modulo 2³² et `0 & 0xFFFFFFFF = 0`, donc le cas est correctement traité comme valide sans condition supplémentaire.

## Importance des tests unitaires

Les manipulations bit à bit sont particulièrement vulnérables aux erreurs off-by-one. Un masque `0x7E` au lieu de `0x7F` ou un décalage de 7 au lieu de 8 compilent sans erreur, mais corrompent silencieusement les résultats. Les six tests couvrent les quatre branches obligatoires de l'énoncé et deux cas limites. `TestAnalyseDonneesValides` est le seul test positif et il vérifie le contenu numérique exact du tableau `[24]int` avec des entrées variées incluant plusieurs bits, un autre capteur ignoré et une valeur nulle. Les trois tests négatifs ciblent chacun une branche d'erreur spécifique. `TestAnalyseCapteurInvalide` passe `capteur=200` et vérifie que la fonction rejette l'appel avant même d'itérer sur les données, car `uint8` accepte 200 mais seuls 0 à 127 sont valides. `TestAnalyseBit7Invalide` soumet une entrée dont le bit 7 est à 1 avec un ID différent du capteur demandé, confirmant que la validation est globale et non limitée au capteur cible. `TestAnalysePlusieursBitsValeur` est critique car il valide directement l'identité `x & (x-1)`. Si la formule était mal implémentée, `0x00000305` (bits 8 et 9 simultanément actifs) ne serait pas détecté. `TestAnalyseExempleEnonce` reproduit exactement l'exemple du professeur, ce qui élimine tout risque d'interprétation erronée de la spécification. Cette stratégie de tests, combinant validation fonctionnelle et couverture des branches, m'a permis de détecter une erreur initiale où la validation ne portait que sur les entrées du capteur demandé.

### Liens

- Dépôt GitHub : [github.com/moyamelissa/Advanved-Programming/tree/main/TN3](https://github.com/moyamelissa/Advanved-Programming/tree/main/TN3)

### Fichiers TN3

- Implémentation principale : [analyse.go](analyse.go)
- Tests unitaires : [analyse_test.go](analyse_test.go)

### Bibliographie / Sources documentaires

- Manuel INF2007, chapitres 3 et 4.
- Documentation Go `math/bits` : https://pkg.go.dev/math/bits
- Documentation Go Testing : https://pkg.go.dev/testing
- A Tour of Go : https://tour.golang.org/
- Outil d'IA : GitHub Copilot — utilisé comme outil d'assistance avec une vérification systématique de chaque suggestion avant intégration.
