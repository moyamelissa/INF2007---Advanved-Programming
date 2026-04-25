# Prompts IA utilisés — TN5

Exemples des prompts utilisés avec l'assistant IA (GitHub Copilot) lors de la réalisation du TN5.

**Note importante :** Chaque résultat généré par l'IA a été systématiquement révisé, testé et validé avant d'être intégré au projet. L'IA a été utilisée comme outil d'assistance, et non comme source de vérité. Toutes les suggestions ont été évaluées avec un esprit critique pour garantir leur exactitude, leur pertinence et leur conformité aux exigences du travail.

## Compréhension de l'énoncé

- « L'énoncé dit de diviser le contenu en segments de N caractères, mais qu'un segment peut se terminer au milieu d'un mot. Est-ce qu'on doit tronquer au mot précédent ou étendre au mot suivant ? »
- « Est-ce que "un mot est une séquence séparée par des espaces" inclut les tabulations et retours à la ligne, ou seulement les espaces simples ? »

## Architecture et concurrence

- « Comment structurer le programme en fan-out / fan-in avec des goroutines et un canal buffered ? Est-ce qu'un `sync.WaitGroup` est nécessaire en plus du canal ? »
- « Est-ce qu'un canal buffered de taille `len(segments)` est suffisant pour éviter les deadlocks, ou faut-il un canal non buffered avec une goroutine de collecte ? »
- « Pourquoi ne pas utiliser un mutex et un compteur partagé au lieu d'un canal pour sommer les résultats ? »

## Gestion du split et des mots coupés

- « Dans `splitIntoSegments`, si la position de coupure tombe au milieu d'un mot, est-ce qu'il vaut mieux reculer au dernier espace ou avancer au prochain espace ? »
- « Comment gérer le cas où le segment contient uniquement des espaces après la coupure ? »
- « Est-ce que `strings.Fields` gère correctement les tabulations, retours chariot et espaces multiples ? »

## Benchmarks et performance

- « Pourquoi est-ce que 100 000 goroutines (segment de 10 caractères) est plus lent que le séquentiel ? Est-ce le coût de création des goroutines ou le scheduling ? »
- « Comment expliquer que le speedup plafonne à 3× sur 8 threads logiques au lieu de 4× ou 8× ? Est-ce un problème de bande passante mémoire ? »
- « Est-ce que la performance devrait croître linéairement avec le nombre de goroutines ? L'énoncé demande d'essayer. »

## Tests

- « Mon test `TestCountWordsConsistency` vérifie le résultat pour 7 tailles de segment différentes. Est-ce que c'est suffisant pour garantir que le split ne perd pas de mots ? »
- « Est-ce qu'il faut tester le cas d'un fichier qui ne contient que des espaces et retours à la ligne ? »

## Rapport

- « Comment expliquer la courbe en U inversé du speedup en fonction du nombre de goroutines de manière concise pour un rapport d'une page ? »
- « Est-ce que la corrélation allocations/performance est un bon indicateur à mentionner dans le rapport ? »
