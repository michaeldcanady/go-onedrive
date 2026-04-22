# Definitive VSA Migration Plan for odc

## 1. Objective
Decompose the project into a proper Vertical Slice Architecture where each feature encapsulates its own domain, infrastructure, and UI (command) logic. 

**STABILITY MANDATE**: NO REVERTS. Fix build errors in-place as they arise.

## 2. Target Structure
For each feature (e.g., identity, fs, config):
```
internal/features/<feature>/
├── domain/         # Interfaces and domain types
├── infrastructure/ # Implementations (e.g., bolt, azure)
└── cmd/            # CLI Commands (Cobra)
```

Cross-cutting concerns:
- `internal/core/<utility>`: logger, env, di, errors, concurrency, shared.
- `internal/core/cli`: Shared CLI logic (e.g., `ProviderPathCompletion`).

## 3. Implementation Roadmap

### Phase 1: Establish Shared CLI Core
1. Create `internal/core/cli`.
2. Move `internal/ui/cli/fs/completion.go` -> `internal/core/cli/path_completion.go`.
3. Update package name to `cli` and update all references.
4. Verify build.

### Phase 2: Migrate Missing Domain Slices
1. Move `internal/drive` -> `internal/features/drive/domain`.
2. Move `internal/editor` -> `internal/features/editor/domain`.
3. Update imports project-wide.
4. Verify build.

### Phase 3: Incremental CLI Migration (One-by-One)
For each feature:
1. Create `internal/features/<feature>/cmd`.
2. Move command files from `internal/ui/cli/<feature>/*` -> `internal/features/<feature>/cmd/`.
3. Fix imports (especially the new slice structure).
4. Resolve circular dependencies (ensure cmd doesn't import something that imports cmd).
5. Verify build.

Order:
- `config`
- `identity`
- `mount`
- `profile`
- `drive`
- `editor`
- `fs` (This is the largest block)

### Phase 4: Final Root Command & Clean up
1. Move `internal/ui/cli/root.go` -> `internal/features/root/root.go` (or leave in `internal/ui/cli` if preferred, but `internal/features/root/cmd` is more VSA).
2. Clean up `internal/ui`.
3. Perform final `go mod tidy` and `go build ./...`.

## 4. Architectural Rules
- Slices must only interact via `domain` interfaces.
- `internal/core/di` wires everything but should not be imported by features.
- CLI handlers must only use the `Service` interface of their own (or other) slices.
