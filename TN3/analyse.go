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

	// Valider l'identifiant du capteur (doit être ≤ 127, c'est-à-dire 7 bits)
	if capteur > 127 {
		return counts, errors.New("identifiant de capteur invalide : doit être entre 0 et 127")
	}

	for _, entry := range data {
		// Extraire l'identifiant du capteur (bits 0–6) avec un masque 0x7F
		id := entry & 0x7F

		// Vérifier le bit de validation (bit 7) avec un masque 0x80
		if entry&0x80 != 0 {
			return counts, errors.New("bit de validation (bit 7) est à 1 : entrée invalide")
		}

		// Extraire les bits 8 à 31 (valeur du capteur)
		valeur := entry >> 8

		// Vérifier qu'au plus un bit parmi les bits 8–31 est à 1
		if bits.OnesCount32(valeur) > 1 {
			return counts, errors.New("plus d'un bit à 1 parmi les bits 8 à 31 : entrée invalide")
		}

		// Compter uniquement pour le capteur demandé
		if id == uint32(capteur) && valeur != 0 {
			// Trouver quel bit est à 1 et incrémenter le compteur correspondant
			for i := 0; i < 24; i++ {
				if valeur&(1<<i) != 0 {
					counts[i]++
					break // Un seul bit est à 1, on peut arrêter
				}
			}
		}
	}

	return counts, nil
}
