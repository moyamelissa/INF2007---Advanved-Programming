package analyse

import (
	"errors"
	"math/bits"
)

// Analyse analyse les données binaires de capteurs IoT pour un capteur donné.
//
// Chaque entrée dans data est un entier 32 bits structuré comme suit :
//   - Bits 0 à 6 : Identifiant du capteur (0–127).
//   - Bit 7 : Bit de validation (doit être 0, sinon erreur).
//   - Bits 8 à 31 : Valeur du capteur. Exactement un de ces bits doit être à 1
//     (mesure spécifique) ou tous à 0 (aucune mesure).
//
// Paramètres :
//   - data : tableau d'entiers 32 bits représentant les données brutes des capteurs.
//   - capteur : identifiant du capteur à analyser (valeur entre 0 et 127).
//
// Retour :
//   - [24]int : tableau de 24 entiers indiquant le nombre de fois où chaque bit
//     (position 8 à 31) est à 1 pour les données valides du capteur spécifié.
//   - error : erreur non nulle si l'identifiant capteur est > 127, si le bit 7
//     est à 1 dans une entrée, ou si plus d'un bit parmi 8 à 31 est à 1.
func Analyse(data []uint32, capteur uint8) ([24]int, error) {
	var counts [24]int

	// Ch3 — Overflow awareness : uint8 va de 0 à 255 (2⁸−1).
	// Les identifiants valides n'utilisent que 7 bits (0 à 127 = 2⁷−1).
	// On doit vérifier explicitement car uint8 ne déborde pas à 127.
	if capteur > 127 {
		return counts, errors.New("identifiant de capteur invalide : doit être entre 0 et 127")
	}

	// Ch4 — Masques : (1<<7)-1 = 0x7F isole les 7 bits de poids faible.
	// C'est la technique "range of bits" du chapitre 4.
	const maskID = (1 << 7) - 1                     // 0x7F : bits 0–6
	const maskBit7 = 1 << 7                         // 0x80 : bit 7 seul
	const maskValeur = ^uint32(0) &^ ((1 << 8) - 1) // bits 8–31 (AND NOT pour effacer bits 0–7)

	for _, entry := range data {
		// Ch4 — AND masking : extraire l'ID capteur (bits 0–6)
		id := entry & maskID

		// Ch4 — Test d'un bit unique : vérifier le bit de validation (bit 7)
		if entry&maskBit7 != 0 {
			return counts, errors.New("bit de validation (bit 7) est à 1 : entrée invalide")
		}

		// Ch4 — Shift right : décaler les bits 8–31 vers les positions 0–23
		// Équivalent à diviser par 2⁸ (Ch3 : shift = division par puissance de 2)
		valeur := entry >> 8

		// Ch4 — Bit identity : x & (x-1) efface le bit le plus bas à 1.
		// Si le résultat est non nul, il y avait au moins 2 bits à 1.
		// Plus efficace que bits.OnesCount car c'est une seule opération.
		if valeur&(valeur-1) != 0 {
			return counts, errors.New("plus d'un bit à 1 parmi les bits 8 à 31 : entrée invalide")
		}

		// Compter uniquement pour le capteur demandé
		if id == uint32(capteur) && valeur != 0 {
			// Ch4 — bits.TrailingZeros : trouver la position du bit à 1
			// en O(1) sans boucle ni branchement (utilise l'instruction CPU BSF/TZCNT).
			// Remplace une boucle de 24 itérations par une seule instruction.
			pos := bits.TrailingZeros32(valeur)
			counts[pos]++
		}
	}

	return counts, nil
}
