# Release Process

This document outlines the structured release lifecycle for `odc`, detailing how
the project transitions from initial development to a stable production release.

## Overview

`odc` uses a multi-stage release process to ensure feature stability and quality
before reaching users. This process is managed automatically via
`release-please` and GitHub Actions, using a tiered versioning strategy.

## Release Stages

The release lifecycle consists of four distinct stages. Each stage has specific
entry and exit criteria to maintain project integrity.

| Stage                 | Suffix     | Purpose                          | Transition Criteria                    |
| :-------------------- | :--------- | :------------------------------- | :------------------------------------- |
| **Alpha**             | `-alpha.x` | Early testing of new features.   | All planned features are implemented.  |
| **Beta**              | `-beta.x`  | Bug squashing and stabilization. | Feature freeze; no new features.       |
| **Release Candidate** | `-rc.x`    | Final regression testing.        | No critical bugs or regressions found. |
| **Stable**            | _(None)_   | Production ready.                | Passed RC phase with zero issues.      |

## Managing Release Phases

You control the current release phase by modifying the project configuration.

### 1. Update Prerelease Type

To move to a different stage (e.g., from Alpha to Beta), update the
`prerelease-type` field in `.github/prerelease-config.json`:

```json
{
  "prerelease-type": "beta"
}
```

### 2. Commit and Push

Once you push this change to the `main` branch, `release-please` will detect the
update and adjust the next version bump accordingly.

### 3. Automated PRs

The system maintains two separate release tracks:

- **Prerelease Track:** Manages the `-alpha`, `-beta`, and `-rc` versions.
- **Stable Track:** Prepares the final stable release by stripping the suffix.

## Release Workflow

The automated workflow follows a specific sequence for every release.

1.  **Code Changes:** When you merge features or fixes into `main`,
    `release-please` updates the pending Prerelease Pull Request.
2.  **Tagging:** Merging the Prerelease PR triggers the creation of a GitHub
    Release and a corresponding tag (e.g., `v1.2.3-beta.1`).
3.  **Finalization:** When a version reaches the `rc` stage and is deemed
    stable, the Release PR from the Stable Track is merged to create the final
    production release (e.g., `v1.2.3`).
4.  **Distribution:** Each new tag triggers the build and distribution pipeline
    for binaries and packages.

## Best Practices

- **Conventional Commits:** Always use conventional commit messages to ensure
  accurate changelog generation and version bumping.
- **Feature Freeze:** Strictly enforce a feature freeze when transitioning from
  Alpha to Beta.
- **Testing:** Ensure all automated tests pass in the CI/CD pipeline before
  merging any Release PR.
