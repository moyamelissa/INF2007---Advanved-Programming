# INF2007 – TN1 – Melissa Moya

![Go Coverage Workflow](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/coverage.yml/badge.svg) ![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

## Justification des cas de test choisis

L'analyse de la fonction DaysUntilDeadline révèle quatre chemins de contrôle distincts : erreurs de parsing (currentDate ou deadline), erreur métier (échéance antérieure) et succès. Cette structure a dicté une stratégie de test visant une couverture d'instructions de 100 %. J'ai rédigé 24 tests organisés en six sections thématiques. La première valide le calcul nominal sur des plages allant jusqu'à dix ans pour confirmer la précision mathématique. La deuxième cible les erreurs de logique liées aux transitions calendaires et années bissextiles. Les sections trois et quatre saturent les erreurs de format (mois hors plage, séparateurs invalides, chaînes vides, format ISO, entrées aléatoires) pour garantir la fiabilité du parsing. La cinquième traite l'erreur métier et la sixième vérifie la prévisibilité de la fonction lorsque les deux entrées sont simultanément invalides.

## Comment les tests garantissent la correction de la fonction

Chacun des 24 tests vérifie systématiquement les deux sorties de la fonction, soit `days` et `err`, assurant un comportement cohérent dans tous les scénarios. Les tests positifs valident la justesse du calcul sur des plages allant de six jours à 3653 jours. Les tests négatifs confirment que chaque catégorie d'entrée invalide déclenche l'erreur attendue tout en retournant zéro jour. Les cas limites ciblent les transitions calendaires et les années bissextiles, des situations reconnues pour générer des anomalies subtiles. La couverture atteint 100 % des instructions, ce qui signifie que les quatre chemins d'exécution sont tous exercés. Cette couverture complète, combinée à la variété des cas testés, offre une validation rigoureuse de la correction de la fonction.
<img width="1036" height="90" alt="image" src="https://github.com/user-attachments/assets/cbe740e4-b55e-4b47-954e-84371cb1e6fd" />

## Défis rencontrés

L'ambiguïté des instances d'erreurs en Go (errors.New) a imposé une validation textuelle via err.Error() pour maintenir le déterminisme des tests. J'ai validé systématiquement le couple de sorties, m'assurant que days reste nul lors d'un échec pour renforcer la stabilité du composant. La gestion de la sévérité via t.Fatalf et t.Errorf a optimisé le diagnostic technique. Au-delà des exigences, l'implémentation du fuzz testing garantit l'absence de paniques face à des entrées corrompues. Enfin, j'ai exploré l'idiomaticité du langage via des tests tabulaires pour centraliser les cas limites et réduire la duplication. L'ajout de benchmarks complète cette approche en évaluant la performance, faisant de cette batterie de tests un investissement durable pour l'évolution du code.

### Liens

- GitHub Repo : https://github.com/moyamelissa/Advanved-Programming/tree/main/TN1
- Vidéo explicative : https://youtu.be/Tsw6rHtLz_k

### Fichiers TN1

- Implémentation principale : [deadline.go](TN1/deadline.go)
- Tests unitaires : [deadline_test.go](TN1/deadline_test.go)

### Bibliographie / Sources documentaires

- Manuel INF2007, chapitre 1 (tests unitaires).
- Documentation Go Testing : https://pkg.go.dev/testing
- Tutoriel Go sur les tests : https://go.dev/doc/tutorial/add-a-test
- A Tour of Go : https://tour.golang.org/
- Outil d'IA : GitHub Copilot — utilisé comme outil d'assistance avec une vérification systématique de chaque suggestion avant intégration. Voir la [liste des prompts utilisés](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN1/AI_Prompts.md).
