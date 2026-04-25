# Travail 3 : Analyse des données binaires

## Informations générales

**Date de remise** : Fin de la semaine 8

**Objectif** : Appliquer les concepts de manipulation de bits du chapitre 4 pour analyser des données binaires de capteurs IoT, en validant les identifiants des capteurs et en détectant des anomalies à l'aide d'opérations bit à bit en Go.

## Contexte pratique

Vous travaillez pour une entreprise qui surveille des capteurs déployés dans une usine intelligente. Chaque capteur envoie des données sous forme d'entiers 32 bits, où :

- **Bits 0 à 6** (7 bits les moins significatifs) : Identifiant du capteur (valeur entre 0 et 127).
- **Bit 7** : Bit de validation (doit être 0, sinon il y a une erreur).
- **Bits 8 à 31** (24 bits les plus significatifs) : Valeur du capteur. Exactement un de ces bits doit être à 1 (indiquant une mesure spécifique), ou tous peuvent être à 0 (indiquant aucune mesure).

Votre tâche est d'écrire un programme Go avec une fonction `Analyse(data []uint32, capteur uint8) ([24]int, error)` qui analyse ces données pour un capteur donné, compte les occurrences de chaque mesure, et valide les données.

## Tâche à réaliser

Implémentez la fonction `Analyse(data []uint32, capteur uint8) ([24]int, error)` qui :

- Prend un tableau d'entiers 32 bits (`data`) et un identifiant de capteur (`capteur`, valeur entre 0 et 127).
- Retourne un tableau de 24 entiers représentant le nombre de fois où chaque bit (8 à 31) est à 1 pour les données valides du capteur spécifié.
- Valide les données et retourne une erreur dans les cas suivants :
  - L'identifiant `capteur` est supérieur à 127.
  - Pour une entrée, le bit 7 est à 1.
  - Pour une entrée, plus d'un bit parmi les bits 8 à 31 est à 1.

### Exemple

```go
data := []uint32{
    0x00000105, // ID=5, bit 7=0, bit 8=1 (valide)
    0x00000205, // ID=5, bit 7=0, bit 9=1 (valide)
    0x00000080, // ID=0, bit 7=1 (invalide)
    0x00000305, // ID=5, bit 7=0, bits 8 et 9=1 (invalide)
}
counts, err := Analyse(data, 5)
// Résultat attendu : err != nil (à cause des entrées invalides)
// Si aucune erreur : counts = [1, 1, 0, 0, ..., 0] (bit 8: 1 fois, bit 9: 1 fois)
```

## Exigences techniques

- **Opérations bit à bit** : Utilisez des opérations comme `&`, `>>`, et masques pour extraire et valider les bits.
- **Tests unitaires** : Fournissez un fichier `analyse_test.go` avec au moins 4 tests :
  - Test positif : Données valides pour un capteur.
  - Test négatif : Bit 7 à 1.
  - Test négatif : Plusieurs bits à 1 parmi les bits 8 à 31.
  - Test négatif : Identifiant de capteur invalide (>127).
- **Documentation** : Documentez la fonction `Analyse` avec des commentaires expliquant son objectif, ses paramètres, et sa valeur de retour.
- **Rapport** : Soumettez un rapport (1 page, PDF) décrivant :
  - Votre approche pour extraire et valider les bits.
  - Les défis rencontrés (par exemple, gestion des erreurs).
  - L'importance des tests unitaires pour ce problème.

## Critères d'évaluation

| Critère | Points | Description |
|---------|--------|-------------|
| Correction | 40 | La fonction gère correctement les cas valides et invalides. |
| Tests unitaires | 25 | Tests complets, clarté des cas de test. |
| Efficacité | 15 | Utilisation optimale des opérations bit à bit, code clair. |
| Documentation et rapport | 20 | Documentation précise, rapport clair et réfléchi. |

## Ressources recommandées

- Manuel, chapitre 4
- [Documentation Go](https://go.dev/doc/)
- [A Tour of Go](https://tour.golang.org/)

## Conseils

- Utilisez un masque comme `0x7F` pour extraire l'ID du capteur et `0x80` pour vérifier le bit 7.
- Pour compter les bits à 1 dans les bits 8 à 31, utilisez une fonction comme `bits.OnesCount32` ou une boucle sur les bits.
- Vous pouvez remettre une vidéo avec votre travail.

Bonne chance et bon codage !
