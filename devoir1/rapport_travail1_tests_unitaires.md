# Rapport - Travail 1 : Tests unitaires

## 1. Contexte et objectif
Ce travail vise a valider le comportement de la fonction DaysUntilDeadline(currentDate, deadline string) dans [deadline.go](deadline.go), qui calcule le nombre de jours entre deux dates au format YYYY-MM-DD.

L objectif principal est de garantir la correction de la fonction avec des tests unitaires couvrant:
- des scenarios valides (tests positifs),
- des scenarios invalides (tests negatifs),
- les branches d erreur importantes.

## 2. Strategie de test
La strategie adoptee combine:
- verification des cas fonctionnels attendus,
- verification de la robustesse face aux entrees invalides,
- verification des conditions limites (exemple: meme date),
- mesure de la couverture pour confirmer que toutes les branches principales sont executees.

Les tests ont ete implementes dans [deadline_test.go](deadline_test.go) avec le package testing de Go.

## 3. Cas de test choisis et justification
### 3.1 Tests positifs
1. Echeance future valide
- Entree: currentDate = 2025-05-26, deadline = 2025-06-01
- Attendu: days = 6, err = nil
- Raison: valide le calcul normal de difference entre deux dates valides.

2. Meme jour
- Entree: currentDate = 2025-05-26, deadline = 2025-05-26
- Attendu: days = 0, err = nil
- Raison: verifie un cas limite important ou aucun jour ne separe les dates.

### 3.2 Tests negatifs
3. Format invalide pour la date courante
- Entree: currentDate = 2025-13-01, deadline = 2025-05-26
- Attendu: err non nil
- Raison: verifie que la fonction detecte une date impossible.

4. Echeance anterieure a la date courante
- Entree: currentDate = 2025-05-26, deadline = 2025-05-25
- Attendu: err non nil
- Raison: verifie la regle metier qui interdit une echeance passee.

5. Mauvais separateur de date
- Entree: currentDate = 2025/05/26, deadline = 2025-06-01
- Attendu: err non nil
- Raison: confirme le respect strict du format YYYY-MM-DD.

6. Format invalide pour la date d echeance
- Entree: currentDate = 2025-05-26, deadline = 2025-13-01
- Attendu: err non nil
- Raison: couvre la branche d erreur du parsing de la date d echeance.

## 4. Comment les tests garantissent la correction
Les tests garantissent la correction de la fonction sur trois axes:
- Precision du resultat numerique: verification de la valeur days dans les cas valides.
- Gestion correcte des erreurs: verification qu une erreur est retournee pour les entrees invalides.
- Respect des regles metier: verification explicite du cas ou deadline < currentDate.

Cette combinaison permet de detecter:
- une erreur de calcul de jours,
- une validation de format insuffisante,
- une absence de controle de coherence temporelle.

## 5. Couverture de tests
La commande go test -cover a ete utilisee. Le projet atteint une couverture de 100 % pour la logique cible de [deadline.go](deadline.go).

Ce resultat depasse l exigence minimale de 90 % demandee par l enonce.

## 6. Defis rencontres
1. Gestion des erreurs
- Defi: identifier toutes les branches d erreur possibles (date courante invalide, date d echeance invalide, echeance anterieure).
- Solution: ajouter des tests negatifs distincts pour chaque cause d erreur.

2. Choix des cas limites
- Defi: couvrir le comportement aux frontieres.
- Solution: ajout du cas meme jour pour valider le retour de 0 jour.

3. Couverture complete
- Defi: atteindre une couverture elevee sans ajouter de tests redondants.
- Solution: ajouter un test cible pour le parsing invalide de la date d echeance, ce qui couvre la branche manquante.

## 7. Conclusion
Le jeu de tests realise dans [deadline_test.go](deadline_test.go) est conforme aux exigences du travail:
- minimum de 5 tests atteint (6 tests implementes),
- cas positifs et negatifs couverts,
- conventions de tests Go respectees,
- couverture superieure a 90 %.

La fonction DaysUntilDeadline est verifiee de facon fiable pour les scenarios fonctionnels principaux et les erreurs courantes de saisie.
