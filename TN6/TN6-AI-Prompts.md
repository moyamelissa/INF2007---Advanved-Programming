# Prompts IA utilisés — TN6

Exemples des prompts utilisés avec l'assistant IA (GitHub Copilot) lors de la réalisation du TN6 (mini-projet).

**Note importante :** Chaque résultat généré par l'IA a été systématiquement révisé, testé et validé avant d'être intégré au projet. L'IA a été utilisée comme outil d'assistance, et non comme source de vérité. Toutes les suggestions ont été évaluées avec un esprit critique pour garantir leur exactitude, leur pertinence et leur conformité aux exigences du travail.

## Compréhension de l'énoncé

- « L'énoncé demande un robot d'exploration Web concurrent. Est-ce qu'on doit suivre les liens trouvés dans les pages (récursif) ou juste explorer une liste fixe d'URLs ? »
- « Est-ce que "compter les mots" signifie uniquement le texte visible, ou aussi le contenu des attributs HTML comme alt et title ? »

## Architecture et concurrence

- « Comment limiter le nombre de goroutines actives à N avec un canal buffered comme sémaphore ? Quelle est la différence avec un worker pool ? »
- « Est-ce qu'un sync.WaitGroup + canal est suffisant ou faut-il aussi un mutex pour agréger les résultats ? »
- « Pourquoi utiliser un canal de type `chan CrawlResult` plutôt qu'un map partagé avec un mutex pour collecter les résultats ? »

## Respect de robots.txt

- « Comment parser robots.txt en Go ? Est-ce qu'il existe une bibliothèque standard ou faut-il utiliser une dépendance externe ? »
- « Si robots.txt n'est pas accessible (erreur réseau, 404), est-ce qu'on autorise l'exploration par défaut ? Quel est le comportement standard des robots d'exploration ? »
- « Comment vérifier qu'un chemin spécifique est autorisé par les règles Disallow et Allow de robots.txt ? »

## Parsing HTML et comptage de mots

- « Comment ignorer le contenu des balises `<script>`, `<style>` et `<noscript>` lors du comptage de mots dans du HTML ? »
- « Est-ce que `golang.org/x/net/html` est la bonne approche pour un tokenizer HTML en Go, ou est-ce qu'un regex serait suffisant ? »
- « Comment gérer les entités HTML (&amp;, &nbsp;, etc.) dans le comptage de mots ? Est-ce que le tokenizer les décode automatiquement ? »

## Tests avec serveur local

- « Comment utiliser `httptest.NewServer` pour créer un serveur de test qui simule robots.txt et des pages HTML ? »
- « Est-ce qu'il faut tester le timeout HTTP séparément ou est-ce que c'est couvert par le test d'URL invalide ? »
- « Comment vérifier qu'une URL bloquée par robots.txt retourne bien une erreur dans `CrawlURLs` ? »

## Benchmarks et performance

- « Comment mesurer l'impact du nombre de goroutines (1, 2, 4, 8) sur le temps de crawl avec un serveur de test local ? »
- « Est-ce que le benchmark du parsing HTML est représentatif si la page de test fait ~1 900 mots ? Faut-il tester avec des pages plus grandes ? »
