# Go Package to Documentation Mapping

This reference maps Go packages in `internal2/` to their primary developer documentation in `docs/developer/`.

---

## 1. Core Architecture
- **Package:** `internal2/` (general)
- **Primary Doc:** `docs/developer/explanation/architecture.md`
- **Related Docs:** `docs/developer/explanation/dependency-injection.md`

## 2. Domain Layer (`internal2/domain`)
- **Package:** `internal2/domain/fs`
- **Primary Doc:** `docs/developer/reference/domain-interfaces.md`
- **Related Docs:** `docs/developer/explanation/filesystem-abstraction.md`

- **Package:** `internal2/domain/auth`
- **Primary Doc:** `docs/developer/reference/domain-interfaces.md`, `docs/developer/explanation/authentication.md`

- **Package:** `internal2/domain/cache`
- **Primary Doc:** `docs/developer/reference/domain-interfaces.md`, `docs/developer/explanation/caching-strategy.md`

## 3. Application Layer (`internal2/app`)
- **Package:** `internal2/app/fs`
- **Primary Doc:** `docs/developer/explanation/filesystem-abstraction.md`

- **Package:** `internal2/app/auth`
- **Primary Doc:** `docs/developer/explanation/authentication.md`

- **Package:** `internal2/app/cache`
- **Primary Doc:** `docs/developer/explanation/caching-strategy.md`

- **Package:** `internal2/app/di`
- **Primary Doc:** `docs/developer/explanation/dependency-injection.md`

## 4. Infrastructure Layer (`internal2/infra`)
- **Package:** `internal2/infra/file`
- **Primary Doc:** `docs/developer/reference/metadata-repository.md`, `docs/developer/reference/contents-repository.md`
- **Related Docs:** `docs/developer/explanation/path-normalization.md`, `docs/developer/how-to/conditional-uploads.md`, `docs/developer/how-to/configure-caching.md`

- **Package:** `internal2/infra/auth`
- **Primary Doc:** `docs/developer/explanation/authentication.md`

## 5. Interface Layer (`internal2/interface`)
- **Package:** `internal2/interface/cli`
- **Primary Doc:** `docs/developer/explanation/cli-interface-layer.md`, `docs/developer/reference/cli-command-patterns.md`, `docs/developer/reference/cli-error-handling.md`
- **How-to Docs:** `docs/developer/how-to/add-subcommand.md`
- **Tutorial Docs:** `docs/developer/tutorials/developing-a-command.md`

- **Package:** `internal2/interface/formatting`
- **Primary Doc:** `docs/developer/reference/cli-output-formatting.md`

- **Package:** `internal2/interface/filtering`, `internal2/interface/sorting`
- **Primary Doc:** `docs/developer/reference/cli-filtering-sorting.md`

---

## Guidelines for New Mappings
When a new package is added to `internal2/`, it should be added to this list along with its primary documentation. If no primary documentation exists, the user should be prompted to create it following the Diátaxis framework.
