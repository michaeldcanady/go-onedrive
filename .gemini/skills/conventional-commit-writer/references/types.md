# Conventional Commit Types

| Type | Description |
| :--- | :--- |
| **feat** | A new feature |
| **fix** | A bug fix |
| **docs** | Documentation only changes |
| **style** | Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc) |
| **refactor** | A code change that neither fixes a bug nor adds a feature |
| **perf** | A code change that improves performance |
| **test** | Adding missing tests or correcting existing tests |
| **build** | Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm) |
| **ci** | Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs) |
| **chore** | Other changes that don't modify src or test files (e.g., agents, skills, CI) |
| **revert** | Reverts a previous commit |

## Choosing the Right Type

- If the change is visible to the user (e.g., a new flag, a new command), use `feat`.
- If the change fixes a bug that was reported or observed, use `fix`.
- If the change is purely about documentation (comments, README, GEMINI.md), use `docs`.
- If the change is about formatting, use `style`.
- If the change is a refactor that doesn't change behavior, use `refactor`.
- If the change is about tests, use `test`.
- If the change is about build scripts, `justfile`, `go.mod`, use `build`.
- If the change is about CI workflows, use `ci`.
- If the change is something else, use `chore`.
