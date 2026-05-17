# Prompts IA utilisés — TN5

Exemples des prompts utilisés avec l'assistant IA (GitHub Copilot) lors de la réalisation du TN5.

**Note importante :** Chaque résultat généré par l'IA a été systématiquement révisé, testé et validé avant d'être intégré au projet. L'IA a été utilisée comme outil d'assistance, et non comme source de vérité. Toutes les suggestions ont été évaluées avec un esprit critique pour garantir leur exactitude, leur pertinence et leur conformité aux exigences du travail.

## Compréhension de l'énoncé

- « Si la coupure tombe au milieu d'un mot, faut-il reculer au mot précédent ou avancer au prochain espace ? »
- « Est-ce que `strings.Fields` considère les tabulations et retours à la ligne comme séparateurs de mots ? »

## Architecture et concurrence

- « Comment structurer un fan-out / fan-in avec des goroutines et un canal en mémoire tampon ? »
- « Un canal de taille `len(segments)` suffit-il pour éviter les interblocages sans goroutine de collecte séparée ? »
- « Pourquoi préférer un canal plutôt qu'un mutex avec un compteur partagé ? »

## Gestion du split et des mots coupés

- « Dans `splitIntoSegments`, vaut-il mieux reculer au dernier espace ou avancer au prochain pour ne pas couper un mot ? »
- « Comment gérer un segment qui ne contient que des espaces après la coupure ? »

## Benchmarks et performance

- « Pourquoi un segment de 10 caractères (environ 100 000 goroutines) est-il plus lent que le séquentiel ? »
- « Comment expliquer que l'accélération plafonne à 1.67× au lieu des 4× théoriques sur 4 cœurs ? La loi d'Amdahl s'applique-t-elle ici ? »
- « La performance devrait-elle croître linéairement avec le nombre de goroutines ? »

## Tests

- « Mon test `TestCountWordsConsistency` couvre 7 tailles de segment. Est-ce suffisant pour garantir que le split ne perd aucun mot ? »
- « Faut-il tester un fichier contenant uniquement des espaces et des retours à la ligne ? »

## Rapport

- « Comment expliquer le plateau d'accélération de façon concise pour un rapport d'une page ? »
- « La corrélation entre le nombre d'allocations et la performance est-elle pertinente à mentionner dans le rapport ? »
