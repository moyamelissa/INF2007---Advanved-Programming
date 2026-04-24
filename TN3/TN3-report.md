# TN3 — Rapport : Analyse des données binaires

## Approche pour extraire et valider les bits

La fonction `Analyse` traite chaque entrée 32 bits en appliquant des masques et décalages bit à bit :

- **Identifiant du capteur (bits 0–6)** : extraction avec le masque `0x7F` via l'opération `entry & 0x7F`, ce qui isole les 7 bits de poids faible.
- **Bit de validation (bit 7)** : vérification avec le masque `0x80` via `entry & 0x80 != 0`. Si ce bit est à 1, l'entrée est invalide.
- **Valeur du capteur (bits 8–31)** : extraction par décalage à droite de 8 positions (`entry >> 8`), puis comptage des bits à 1 avec `bits.OnesCount32`. Si plus d'un bit est actif, l'entrée est invalide. Si exactement un bit est à 1 et que l'identifiant correspond au capteur recherché, on incrémente le compteur de la position correspondante.

Cette approche garantit que chaque validation est réalisée avec des opérations O(1) sur les bits, sans conversion en chaîne ni allocation mémoire supplémentaire.

## Défis rencontrés

Le principal défi a été la gestion rigoureuse des erreurs : la fonction doit valider **toutes** les entrées du tableau, pas uniquement celles correspondant au capteur cible. Une entrée invalide (bit 7 à 1 ou plusieurs bits de mesure) pour un capteur différent doit quand même déclencher une erreur. Il a aussi fallu s'assurer que le type `uint8` du paramètre `capteur` permet bien de détecter les valeurs > 127, puisque `uint8` va jusqu'à 255.

Un autre point d'attention est la distinction entre « aucune mesure » (bits 8–31 tous à 0, qui est valide) et une mesure active (exactement un bit à 1). Il ne faut pas compter les entrées sans mesure dans le tableau de résultats.

## Importance des tests unitaires

Les tests unitaires sont essentiels pour ce type de problème car les manipulations bit à bit sont sujettes à des erreurs subtiles (décalage d'un bit, masque incorrect, oubli d'un cas limite). Les tests couvrent :

1. **Cas positif** : données valides avec comptage correct des mesures.
2. **Bit 7 invalide** : détection d'une erreur de validation.
3. **Plusieurs bits de mesure** : détection de données corrompues.
4. **Capteur hors limites** : validation du paramètre d'entrée.
5. **Tableau vide** et **exemple de l'énoncé** : couverture des cas limites.

Ces tests permettent de vérifier chaque branche de validation indépendamment et d'éviter les régressions lors de modifications futures.
