# Prompts IA utilisés — TN4

Exemples des prompts utilisés avec l'assistant IA (GitHub Copilot) lors de la réalisation du TN4.

**Note importante :** Chaque résultat généré par l'IA a été systématiquement révisé, testé et validé avant d'être intégré au projet. L'IA a été utilisée comme outil d'assistance, et non comme source de vérité. Toutes les suggestions ont été évaluées avec un esprit critique pour garantir leur exactitude, leur pertinence et leur conformité aux exigences du travail.

## Compréhension de l'énoncé

- « L'énoncé demande de mesurer à 1%, 10%, 20%... 100% du tableau. Est-ce qu'on doit utiliser des sous-benchmarks avec `b.Run` ou des fonctions de benchmark séparées pour chaque pourcentage ? »
- « Est-ce qu'il faut pré-générer le tableau avant les benchmarks ou le régénérer à chaque itération ? »

## Architecture et choix de conception

- « Est-ce que séparer `computeSineSumInt` et `computeSineSumFloat` dans des fonctions distinctes est préférable à utiliser des génériques, sachant que les benchmarks doivent mesurer directement la boucle de calcul ? »
- « Pourquoi utiliser `rand.NewSource(42)` plutôt que la seed par défaut ou `crypto/rand` pour la génération du tableau ? »
- « Comment structurer le dispatch `computeSineSum` avec `interface{}` pour que les benchmarks puissent quand même appeler directement les fonctions typées ? »

## Benchmarks et performance

- « Comment interpréter la colonne ns/op de la sortie de `go test -bench` et la convertir en millisecondes pour un tableau ? »
- « Pourquoi les flottants sont plus rapides que les entiers alors que `math.Sin` retourne un `float64` dans les deux cas ? »
- « Est-ce que l'instruction CPU `CVTSI2SD` est bien ce qui explique le surcoût de la version Int, ou est-ce qu'il y a d'autres facteurs ? »

## Tests

- « Est-ce qu'une tolérance de 1e-9 est appropriée pour comparer des résultats de `math.Sin` sur 3 éléments ? »
- « Mon test sur le tableau vide vérifie que la somme est 0, mais est-ce qu'il faudrait aussi vérifier que `math.Sin` n'est jamais appelé ? »

## Rapport et présentation

- « Comment générer un graphique Mermaid `xychart-beta` dans un fichier Markdown pour comparer les courbes Int et Float ? »
- « Est-ce que je peux forcer un fond blanc sur un graphique Mermaid pour que ça rende bien en mode sombre sur GitHub ? »
- « Comment calculer la distance parcourue par la lumière pendant un appel à `math.Sin` et le nombre de sinus par frame à 120 fps ? »
