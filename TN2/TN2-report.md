# INF2007 - TN2 - Melissa Moya

## Workflow Git 
Le dépôt a été initialisé conformément au protocole de laboratoire, la branche main étant utilisée comme branche de référence du projet. Un premier commit a permis d’y figer une version fonctionnelle minimale contenant uniquement la fonction countLines, matérialisant le principe du tronc commun stable propre aux systèmes de contrôle de version distribués. À partir de main, deux branches fonctionnelles ont été créées, soit count-words et count-chars. Chacune était dédiée à l’implémentation d’une fonctionnalité distincte, soit le comptage des mots et le comptage des caractères. Cette organisation applique le concept de développement parallèle, permettant l’évolution simultanée de plusieurs fonctionnalités sans modification directe de la branche principale. Les modifications sur chaque branche ont été validées par des commits atomiques, limités à une seule responsabilité fonctionnelle. Cette pratique favorise la lisibilité de l’historique et la traçabilité des changements. Conformément aux consignes, la branche count-words a été fusionnée en premier dans main, entraînant une fusion de type fast-forward, illustrant une intégration sans divergence préalable dans le graphe des commits.

## Résolution du conflit de fusion
Un conflit de fusion a été provoqué lors de la fusion subséquente de count-chars dans main, à la suite de modifications concurrentes apportées à la fonction main() dans le fichier main.go. Chaque branche introduisait un affichage différent, empêchant toute fusion automatique. Lors de l’exécution de git merge, Git a détecté la divergence et inséré des marqueurs de conflit standards. La résolution a été effectuée manuellement par l’analyse de ces marqueurs, permettant d’identifier les sections issues de chaque branche. Après suppression complète des marqueurs, une résolution sémantique a été appliquée en combinant les deux intentions de développement. La fonction main() conserve ainsi l’affichage du nombre de mots et du nombre de caractères, préservant toutes les fonctionnalités développées en parallèle. Cette étape illustre une caractéristique des systèmes distribués où Git opère localement sur des différences textuelles, contrairement à un système centralisé, tandis que la cohérence fonctionnelle relève de l’intervention humaine. La fusion a été finalisée par un commit documentant explicitement la résolution.

## Défis rencontrés
Le principal défi technique rencontré a concerné l’atteinte d’une couverture de tests complète. La fonction main(), bien que simple, repose sur des entrées et sorties standard, ce qui complique son intégration dans des tests automatisés. Afin de valider son exécution sans perturber la sortie standard, une redirection explicite de stdout vers os.DevNull a été mise en place. Cette adaptation a permis d’atteindre une couverture de 100 %, illustrant le lien entre validation logicielle et intégration du code dans un contexte de développement contrôlé. L’expérimentation met ainsi en évidence que la qualité de l’historique Git et l’intégration des fonctionnalités doivent être accompagnées d’une stratégie de tests adaptée afin de garantir la cohérence globale du projet.

---

### Liens

- Dépôt GitHub : [github.com/moyamelissa/Advanced-Programming/tree/main/TN2](https://github.com/moyamelissa/Advanced-Programming/tree/main/TN2)
- Guide d'expérimentation étape par étape : [TN2-Experimentation-Guide.md](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN2/TN2-Experimentation-Guide.md)

### Fichiers TN2

- Code source principal : [main.go](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN2/word-stats/main.go)
- Tests unitaires : [main_test.go](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN2/word-stats/main_test.go)
- Historique des commits : [history.txt](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN2/word-stats/history.txt)
- Documentation du projet : [README.md](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN2/word-stats/README.md)

### Bibliographie / Sources documentaires

- Manuel INF2007, chapitre 2
- Documentation Go Testing : https://pkg.go.dev/testing
- Tutoriel Go sur les tests : https://go.dev/doc/tutorial/add-a-test
- A Tour of Go : https://tour.golang.org/
- Outil d'IA : GitHub Copilot — utilisé comme outil d'assistance avec une vérification systématique de chaque suggestion avant intégration. Voir la [liste des prompts utilisés](https://github.com/moyamelissa/Advanced-Programming/blob/main/TN2/TN2-AI-Prompts.md).
