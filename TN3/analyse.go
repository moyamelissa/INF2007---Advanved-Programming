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

	// uint8 accepte 0–255, mais les identifiants valides n'utilisent que 7 bits (0–127).
	// On vérifie explicitement car le type ne contraint pas cette plage. (cf. Ch. 3)
	if capteur > 127 {
		return counts, errors.New("identifiant de capteur invalide : doit être entre 0 et 127")
	}

	// Masques construits par décalage plutôt qu'en dur pour lisibilité. (cf. Ch. 4)
	const maskID = (1 << 7) - 1 // 0x7F : bits 0–6
	const maskBit7 = 1 << 7     // 0x80 : bit 7 seul

	for _, entry := range data {
		// Extraire l'ID capteur (bits 0–6) par AND masking. (cf. Ch. 4)
		id := entry & maskID

		// Vérifier le bit de validation (bit 7). (cf. Ch. 4)
		if entry&maskBit7 != 0 {
			return counts, errors.New("bit de validation (bit 7) est à 1 : entrée invalide")
		}

		// Décaler les bits 8–31 vers les positions 0–23 (>> 8 = ÷ 2⁸). (cf. Ch. 3)
		valeur := entry >> 8

		// x & (x-1) efface le bit le plus bas à 1.
		// Si le résultat est non nul, il y a plus d'un bit actif. (cf. Ch. 4)
		if valeur&(valeur-1) != 0 {
			return counts, errors.New("plus d'un bit à 1 parmi les bits 8 à 31 : entrée invalide")
		}

		// Compter uniquement pour le capteur demandé
		if id == uint32(capteur) && valeur != 0 {
			// Trouver la position du bit à 1 en O(1) via l'instruction CPU TZCNT. (cf. Ch. 4)
			pos := bits.TrailingZeros32(valeur)
			counts[pos]++
		}
	}

	return counts, nil
}
