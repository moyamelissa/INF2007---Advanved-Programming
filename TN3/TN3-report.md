# TN3 — Rapport : Analyse des données binaires

**Cours** : INF2007 — Programmation avancée

---

## 1. Concepts des chapitres 3 et 4 appliqués

Ce travail met en pratique les concepts de représentation des données (Ch. 3) et de manipulation de bits (Ch. 4). Voici un tableau récapitulatif des concepts utilisés dans le code, suivi d'explications détaillées :

| Concept (Chapitre) | Où dans le code | Pourquoi |
|---------------------|----------------|----------|
| Mots de 32 bits, représentation hexadécimale (Ch. 3) | `data []uint32`, `0x7F`, `0x80` | Les données capteur sont des mots de 32 bits, manipulés en hexadécimal pour lisibilité |
| Entiers non signés, plage 0 à 2ⁿ−1 (Ch. 3) | `uint32`, `uint8` | `uint8` va de 0 à 255 (2⁸−1) ; on vérifie explicitement `> 127` (2⁷−1) |
| Shift = division par puissance de 2 (Ch. 3) | `entry >> 8` | Décaler de 8 = diviser par 2⁸ = ramener les bits 8–31 en positions 0–23 |
| Tableaux de taille fixe (Ch. 3) | `[24]int` | 24 compteurs pour 24 bits de mesure possibles |
| AND masking `&` (Ch. 4) | `entry & maskID` | Isoler les bits 0–6 de l'identifiant |
| Bit range mask `(1<<n)-1` (Ch. 4) | `const maskID = (1 << 7) - 1` | Construire un masque de n bits consécutifs |
| AND NOT `&^` pour effacer des bits (Ch. 4) | `^uint32(0) &^ ((1<<8)-1)` | Construire le masque des bits 8–31 en effaçant les bits 0–7 |
| Test d'un bit unique `1 << k` (Ch. 4) | `entry & maskBit7` | Vérifier si le bit 7 est à 1 |
| **Bit identity `x & (x-1)`** (Ch. 4) | `valeur & (valeur-1) != 0` | Détecter > 1 bit à 1 **sans boucle ni `OnesCount`** |
| **`bits.TrailingZeros` — branchless** (Ch. 4) | `bits.TrailingZeros32(valeur)` | Trouver la position du bit à 1 en **O(1)** sans boucle |

## 2. Approche pour extraire et valider les bits

### Construction des masques (Ch. 4 — bit ranges)

Plutôt que d'écrire les masques en dur (`0x7F`), j'ai utilisé la technique du chapitre 4 pour les construire à partir de la formule `(1<<n)-1` :

```go
const maskID    = (1 << 7) - 1        // 0x7F : bits 0–6
const maskBit7  = 1 << 7              // 0x80 : bit 7 seul
const maskValeur = ^uint32(0) &^ ((1 << 8) - 1)  // bits 8–31
```

`maskValeur` utilise l'opérateur **AND NOT** (`&^`) : on part de tous les bits à 1 (`^uint32(0) = 0xFFFFFFFF`), puis on efface les bits 0–7 avec `&^ ((1<<8)-1)`. Cette technique est plus explicite que d'écrire `0xFFFFFF00` en dur, car elle montre clairement qu'on conserve 32−8 = 24 bits.

### Extraction de l'identifiant (Ch. 4 — AND masking)

```go
id := entry & maskID    // AND : isole bits 0–6
```

Pour `entry = 0x00000105` (binaire `...0001 0000 0101`), le résultat est `id = 5`.

### Vérification du bit 7 (Ch. 4 — test d'un bit)

```go
if entry&maskBit7 != 0 {
    return counts, errors.New("bit de validation (bit 7) est à 1 : entrée invalide")
}
```

On valide **toutes les entrées**, pas seulement celles du capteur cible. L'exemple de l'énoncé confirme ce choix : `0x00000080` (ID=0, bit 7=1) déclenche une erreur même quand on cherche le capteur 5.

### Shift = division par 2⁸ (Ch. 3 — shifts)

```go
valeur := entry >> 8
```

Le décalage de 8 positions équivaut à diviser par 2⁸ = 256. Les bits 8–31 se retrouvent en positions 0–23, prêts à être analysés.

### Bit identity `x & (x-1)` — détection multi-bits (Ch. 4)

C'est l'optimisation la plus importante. Le chapitre 4 enseigne que `x & (x-1)` efface le bit le plus bas à 1 :

```go
if valeur&(valeur-1) != 0 {
    return counts, errors.New("plus d'un bit à 1 parmi les bits 8 à 31")
}
```

**Comment ça fonctionne :**
- Si `valeur = 0b00000100` (1 seul bit) : `0b00000100 & 0b00000011 = 0` → valide ✓
- Si `valeur = 0b00000110` (2 bits) : `0b00000110 & 0b00000101 = 0b00000100 ≠ 0` → invalide ✗
- Si `valeur = 0` (aucun bit) : `0 & (0-1) = 0 & 0xFFFFFFFF = 0` → valide ✓ (overflow unsigned wraps, Ch. 3)

Cette identité remplace `bits.OnesCount32(valeur) > 1` qui nécessite un comptage complet. Ici, une seule soustraction + un AND suffisent. C'est exactement le type de code « branchless et élégant » recommandé au chapitre 4.

### `bits.TrailingZeros` — localisation O(1) du bit (Ch. 4 — bit counting)

Au lieu d'une boucle de 0 à 23 pour trouver quel bit est à 1, j'utilise :

```go
pos := bits.TrailingZeros32(valeur)
counts[pos]++
```

`TrailingZeros32` retourne le nombre de zéros après le bit le moins significatif. Sur x86-64, cela se traduit par l'instruction CPU `TZCNT` (ou `BSF`) qui s'exécute en 1 cycle. Cela remplace une boucle de jusqu'à 24 itérations par **une seule instruction machine** — un exemple concret d'utilisation des fonctions `math/bits` recommandées au chapitre 4.

## 3. Défis rencontrés

### Overflow awareness (Ch. 3)

Le paramètre `capteur` est `uint8` (0–255, soit 2⁸−1). Les identifiants valides sont 0–127 (2⁷−1). Le chapitre 3 explique que les entiers non signés ne « dépassent » pas de la même façon que les signés — `uint8(200)` est parfaitement valide mais hors de la plage des identifiants. La vérification explicite est donc nécessaire :

```go
if capteur > 127 {
    return counts, errors.New("identifiant de capteur invalide")
}
```

### Distinction « aucune mesure » vs « mesure invalide »

Quand `valeur = 0` (aucun bit à 1), l'identité `0 & (0-1) = 0` est correcte grâce au comportement modulo 2ⁿ des entiers non signés (Ch. 3) : `uint32(0) - 1 = 0xFFFFFFFF` (wrap-around), et `0 & 0xFFFFFFFF = 0`. C'est un cas où la connaissance de l'arithmétique unsigned est essentielle pour éviter un bug.

## 4. Tests unitaires

J'ai écrit **6 tests** couvrant chaque branche de validation :

| Test | Ce qu'il vérifie | Concept Ch. 3/4 testé |
|------|-----------------|----------------------|
| `TestAnalyseDonneesValides` | bit 8 compté 2×, bit 9 compté 1×, bit 23 compté 1×. Ignore ID=3 et entrées sans mesure. | AND masking, TrailingZeros, shift |
| `TestAnalyseBit7Invalide` | `0x00000080` → erreur | Test bit unique `1<<7` |
| `TestAnalysePlusieursBitsValeur` | `0x00000305` → erreur | Bit identity `x & (x-1)` |
| `TestAnalyseCapteurInvalide` | `capteur=200` → erreur | Overflow/plage uint8 (Ch. 3) |
| `TestAnalyseTableauVide` | `[]uint32{}` → zéros, pas d'erreur | Tableau fixe [24]int (Ch. 3) |
| `TestAnalyseExempleEnonce` | Reproduit l'exemple de l'énoncé | Validation complète |

Le test `TestAnalysePlusieursBitsValeur` est crucial pour vérifier l'identité `x & (x-1)` : si la formule était incorrecte, `0x00000305` (bits 8 et 9) ne serait pas détecté comme invalide.

## 5. Récapitulatif : concepts non utilisés et justification

Certains concepts des chapitres 3–4 n'apparaissent pas dans ce travail car ils ne s'appliquent pas au problème :

- **Flottants / IEEE 754** (Ch. 3) : les données capteur sont des entiers.
- **Strings / Unicode / UTF-8** (Ch. 3) : pas de traitement de texte.
- **SWAR** (Ch. 4) : utile pour traiter 8 octets en parallèle, mais nos données sont des mots individuels de 32 bits.
- **Rotations de bits** (Ch. 4) : non nécessaires car la structure des données est fixe (pas de chiffrement ni de hashing).
- **Structures/interfaces** (Ch. 3) : le problème ne nécessite qu'une seule fonction, pas d'abstraction supplémentaire.
