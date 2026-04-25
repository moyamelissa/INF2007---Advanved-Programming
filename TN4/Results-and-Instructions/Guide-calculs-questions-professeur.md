# Instructions pour les questions spéciales

## Question 1 : Distance parcourue par la lumière pendant un sinus

### Données

- Vitesse de la lumière : 299 792 458 m/s
- Temps pour un sinus : on le lit dans la dernière ligne du benchmark (100pct, soit 1 000 000 éléments)

### Comment trouver le temps par sinus

On prend le `ns/op` du benchmark à 100% et on divise par le nombre d'éléments (1 000 000).

```
Temps par sinus = ns/op du benchmark 100pct / 1 000 000
```

Par exemple, avec nos résultats :
- Int : `36 442 568 ns/op / 1 000 000 = 36.4 ns par sinus`
- Float : `21 053 744 ns/op / 1 000 000 = 21.1 ns par sinus`

### Calcul de la distance

On convertit les nanosecondes en secondes (diviser par 1 000 000 000), puis on multiplie par la vitesse de la lumière.

```
distance = vitesse × temps
distance = 299 792 458 m/s × (temps_par_sinus / 1 000 000 000)
```

Avec nos résultats :
- Int : 299 792 458 × 36.4 / 1 000 000 000 = 10.9 mètres
- Float : 299 792 458 × 21.1 / 1 000 000 000 = 6.3 mètres

La lumière parcourt entre 6 et 11 mètres pendant un seul calcul de sinus, selon le type.

## Question 2 : Combien de sinus dans un tick à 120 fps

### Données

- Fréquence : 120 images par seconde
- Temps par sinus : calculé à la question 1

### Comment trouver la durée d'un tick

```
durée d'un tick = 1 seconde / 120 = 0.008333 seconde = 8 333 333 nanosecondes
```

### Calcul du nombre de sinus par tick

On divise la durée du tick par le temps d'un sinus.

```
nombre de sinus = durée du tick / temps par sinus
nombre de sinus = 8 333 333 ns / temps_par_sinus_en_ns
```

Avec nos résultats :
- Int : 8 333 333 / 36.4 = environ 229 000 sinus par tick
- Float : 8 333 333 / 21.1 = environ 395 000 sinus par tick

On peut donc calculer entre 229 000 et 395 000 sinus par frame à 120 fps sur cette machine.
