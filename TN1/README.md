# TN1 - Unit Testing and Coverage

## English

### Coverage Setup

1. Push this repository to GitHub.
2. Sign in to Codecov and add this repository.
3. In GitHub repository settings, add secret `CODECOV_TOKEN` (if Codecov requires a token for your repository).
4. Push to `main` or open a pull request to trigger the workflow.

### Run Tests Locally

```bash
cd TN1
go test ./... -cover
```

## Francais

### Badges

![Workflow Go Coverage](https://github.com/moyamelissa/Advanved-Programming/actions/workflows/coverage.yml/badge.svg)
![Codecov](https://codecov.io/gh/moyamelissa/Advanved-Programming/branch/main/graph/badge.svg)

### Configuration de la couverture

1. Poussez ce depot sur GitHub.
2. Connectez-vous a Codecov et ajoutez ce depot.
3. Dans les parametres GitHub du depot, ajoutez le secret `CODECOV_TOKEN` (si Codecov demande un token pour votre depot).
4. Poussez sur `main` ou ouvrez une pull request pour declencher le workflow.

### Lancer les tests en local

```bash
cd TN1
go test ./... -cover
```
