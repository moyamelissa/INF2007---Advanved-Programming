# Prompts IA utilisés — TN3

Exemples des prompts utilisés avec l'assistant IA (GitHub Copilot) lors de la réalisation du TN3.

**Note importante :** Chaque résultat généré par l'IA a été systématiquement révisé, testé et validé avant d'être intégré au projet. L'IA a été utilisée comme outil d'assistance, et non comme source de vérité. Toutes les suggestions ont été évaluées avec un esprit critique pour garantir leur exactitude, leur pertinence et leur conformité aux exigences du travail.

## Compréhension de l'énoncé

- « Si le bit 7 est à 1 dans une entrée qui n'est pas celle du capteur demandé, est-ce que ça doit quand même retourner une erreur ? »
- « Est-ce qu'une entrée avec tous les bits 8 à 31 à 0 est considérée valide ou invalide selon l'énoncé ? »

## Choix techniques et opérations bit à bit

- « Est-ce que `x & (x-1)` est une approche correcte pour détecter si plus d'un bit est actif, et comment ça se comporte quand x vaut 0 ? »
- « Quelle est la différence entre utiliser `bits.OnesCount32` et l'identité `x & (x-1)` pour vérifier qu'un seul bit est à 1 ? »
- « Est-ce que `bits.TrailingZeros32` est plus efficace qu'une boucle pour trouver la position d'un bit actif ? »

## Validation et tests

- « Est-ce que mes six tests couvrent bien toutes les branches d'erreur demandées par l'énoncé ? »
- « Est-ce que le test `TestAnalyseBit7Invalide` devrait utiliser une entrée avec le même capteur ou un capteur différent pour mieux tester la validation globale ? »

## Style et commentaires

- « Comment écrire professionnellement une référence à un chapitre du cours dans mes commentaires ? »
- « Est-ce que mes commentaires dans le code ont l'air trop académiques ? Comment les rendre plus professionnels ? »

...
