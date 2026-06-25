# Contributing to ripchords

This project uses a few lightweight conventions so that automated tooling
(linting, versioning, releases) can do its job. They're summarized here.

## One-time setup

```bash
make tools   # install golangci-lint (pinned version) into $GOPATH/bin
make hooks   # enable the pre-commit hook (runs checks before each commit)
```

## Branch names

Branches follow `<type>/<short-kebab-description>`, where `<type>` is one of the
Conventional Commit types below. Including the issue number is encouraged.

```
fix/18-rr-confirm
feat/capo-support
ci/add-linting
```

This is enforced by the `branch-name` job in `.github/workflows/pr-checks.yml`.

## Commits & PR titles (Conventional Commits)

Versioning and releases are automated by release-please, which reads commit
messages, so they must follow [Conventional Commits](https://www.conventionalcommits.org/).
Because PRs are **squash-merged**, the **PR title** becomes the commit on `main`
— that's what matters, and it's validated by the `conventional-title` CI job.

| Type        | Effect (pre-1.0)          | Example                              |
|-------------|---------------------------|--------------------------------------|
| `fix:`      | patch bump                | `fix: correct mini-barre render`     |
| `feat:`     | minor bump                | `feat: add capo support`             |
| `feat!:`    | minor bump (breaking <1.0)| `feat!: change settings format`      |
| `chore:`/`docs:`/`refactor:`/`test:`/`ci:`/`build:`/`perf:` | no release | `ci: add linting workflow` |

## Checks

Run the full check suite locally before pushing (CI runs the same thing):

```bash
make check   # go vet + golangci-lint + go test
```

Individual targets: `make build`, `make test`, `make vet`, `make lint`.

CI is defined in:
- `.github/workflows/ci.yml` — build, vet, test, golangci-lint
- `.github/workflows/pr-checks.yml` — PR-title (Conventional Commits) + branch name
- `.github/workflows/release-please.yml` — automated versioning/releases
