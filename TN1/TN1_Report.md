# INF2007 – TN1 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanced-Programming/actions/workflows/coverage.yml/badge.svg) [![codecov](https://codecov.io/gh/moyamelissa/Advanced-Programming/branch/main/graph/badge.svg?flag=tn1)](https://codecov.io/gh/moyamelissa/Advanced-Programming)

## Justification des cas de test choisis

L'analyse de la fonction DaysUntilDeadline révèle quatre chemins de contrôle distincts : erreur de parsing de currentDate, erreur de parsing de deadline, erreur métier (échéance antérieure) et succès.. Cette structure a dicté une stratégie de test visant une couverture d'instructions de 100%. J'ai rédigé 24 tests organisés en six sections thématiques. La première valide le calcul nominal sur des plages allant jusqu'à dix ans pour confirmer la précision mathématique. La deuxième cible les erreurs de logique liées aux transitions calendaires et années bissextiles. Les sections trois et quatre saturent les erreurs de format (mois hors plage, séparateurs invalides, chaînes vides, format ISO, entrées aléatoires) pour garantir la fiabilité du parsing. La cinquième traite l'erreur métier et la sixième vérifie la prévisibilité de la fonction lorsque les deux entrées sont simultanément invalides.

## Comment les tests garantissent la correction de la fonction

Chaque test valide systématiquement le couple de sorties (days, err), garantissant un comportement cohérent dans tous les scénarios. Les tests positifs confirment la justesse du calcul sur des plages étendues. Les tests négatifs explorent les catégories d'entrées invalides pour assurer le déclenchement des erreurs spécifiques tout en vérifiant le retour de zéro jour. L'inclusion de cas limites (bissextiles, transitions calendaires) neutralise les anomalies subtiles inhérentes aux dates. Enfin, la couverture de 100 % des instructions certifie que les quatre chemins de contrôle identifiés sont exercés. Cette approche exhaustive, alliant couverture structurelle et variété fonctionnelle, offre une validation rigoureuse de la correction logicielle.
<img width="1036" height="90" alt="image" src="https://github.com/user-attachments/assets/cbe740e4-b55e-4b47-954e-84371cb1e6fd" />

## Défis rencontrés

Go ne permet pas de comparer directement les instances d'erreurs créées avec errors.New. J'ai donc utilisé err.Error() pour valider le message exact retourné, ce qui maintient le déterminisme des assertions. Le choix entre t.Fatalf et t.Errorf a été guidé par la dépendance logique : Fatalf arrête le test lorsqu'un résultat subséquent dépend d'une valeur critique, tandis que Errorf permet de capturer plusieurs échecs dans un même test. J'ai aussi implémenté des tests tabulaires (table-driven) pour réduire la verbosité du code et centraliser les cas limites en un seul point de maintenance. Enfin, j'ai ajouté du fuzz testing pour vérifier l'absence de paniques face à des entrées aléatoires, ainsi que des benchmarks pour mesurer la performance du parsing et du calcul.

### Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Vidéo explicative : https://youtu.be/jC_luJYEZWw

### Fichiers TN1

- Implémentation principale : [deadline.go](TN1/deadline.go)
- Tests unitaires : [deadline_test.go](TN1/deadline_test.go)

### Bibliographie / Sources documentaires

- Manuel INF2007, chapitre 1 (tests unitaires).
- Documentation Go Testing : https://pkg.go.dev/testing
- Tutoriel Go sur les tests : https://go.dev/doc/tutorial/add-a-test
- A Tour of Go : https://tour.golang.org/
- Outil d'IA : GitHub Copilot — utilisé comme outil d'assistance avec une vérification systématique de chaque suggestion avant intégration. Voir la [liste des prompts utilisés](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN1/AI_Prompts.md).
