# Calculs détaillés - TN6

Ce fichier documente les calculs numériques cités dans le rapport TN6.
Source des valeurs : benchstat sur 6 exécutions (-count=6), Windows/amd64.

---

## 1) Formule d'accélération (speedup)

Formule :

```text
Accélération = Temps de référence / Temps mesuré
```

Le temps de référence est toujours la configuration 1 goroutine du même scénario.

### 1.1 Serveur unique

Référence : 1.905 ms

- 2 goroutines : 1.905 / 1.176 = 1.6199 -> **1.62x**
- 4 goroutines : 1.905 / 3.094 = 0.6157 -> **0.62x**
- 8 goroutines : 1.905 / 10.80 = 0.1764 -> **0.18x**

### 1.2 Multi-serveurs

Référence : 3.135 ms

- 2 goroutines : 3.135 / 1.719 = 1.8237 -> **1.82x**
- 4 goroutines : 3.135 / 0.9495 = 3.3017 -> **3.30x**
- 8 goroutines : 3.135 / 0.7056 = 4.4422 -> **4.44x**

---

## 2) Comment lire les pourcentages "±"

Les valeurs "± N %" viennent directement de benchstat.
Elles représentent la dispersion relative des exécutions répétées.

Exemples cités dans le rapport :

- Serveur unique, 4 goroutines : **3.094 ms ± 196 %** (très instable)
- Multi-serveurs, 8 goroutines : **0.7056 ms ± 2 %** (stable)
- CountWordsHTML : **345.7 us ± 35 %** (variabilité notable)

---

## 3) Médiane explicite (exemple 4 goroutines, serveur unique)

Le rapport indique que 3 exécutions sont proches de 1.78 ms et 3 autres sont plus lentes (~4.4, ~8.6, ~9.1 ms).

Avec 6 mesures triées :

```text
1.78, 1.78, 1.78, 4.4, 8.6, 9.1
```

Pour un nombre pair de mesures, la médiane est la moyenne des 2 valeurs centrales :

```text
Médiane = (3e + 4e) / 2
Médiane = (1.78 + 4.4) / 2 = 3.09 ms
```

Cela correspond à la valeur benchstat du tableau.

---

## 4) Rendement parallèle (optionnel mais utile)

Formule d'efficacité parallèle :

```text
Efficacité = Accélération / Nombre de goroutines
```

Exemple multi-serveurs à 8 goroutines :

```text
E = 4.44 / 8 = 0.555 -> 55.5%
```

Interprétation : bon gain global, mais avec rendement décroissant (normal en parallèle).

---
