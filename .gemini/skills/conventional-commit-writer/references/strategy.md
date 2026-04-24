# Commit Strategy & Timing

Effective use of conventional commits relies on a solid commit strategy. Follow these guidelines to decide *when* and *how* to commit.

## Atomic Commits

A commit should represent a single logical change. This makes it easier to:
- Review code.
- Revert specific changes.
- Track down bugs (using `git bisect`).
- Maintain a clean history.

**Guideline:** If you find yourself using "and" in your commit message (e.g., "fix login bug and update styles"), you should probably split the commit.

## When to Commit

1.  **When a sub-task is complete**: Don't wait until the entire feature is done. Commit as soon as a meaningful, self-contained part of the work is finished and verified.
2.  **Before and after refactoring**: Never mix refactoring with feature additions or bug fixes. Refactor first, commit, then implement the change.
3.  **Before switching tasks**: If you need to stop working on one thing to fix a high-priority bug, commit your progress (even as a WIP) or use `git stash`.
4.  **When tests pass**: Only commit code that is in a "green" state.

## Handling Large Changes

If you've already made a significant number of changes that span multiple logical areas:

1.  **Analyze the diff**: Run `git diff` and identify the distinct logical changes.
2.  **Use Selective Staging**: Use `git add -p` (patch mode) to interactively choose which hunks of code to include in the next commit.
3.  **Stage by File**: If changes are in separate files, `git add <file>` to commit them one by one.
4.  **Refactor vs. Logic**: If you've mixed formatting or refactoring with logic changes, try to stage the refactors first.

## Squashing vs. Splitting

- **Commit early and often**: It's much easier to squash multiple small commits into one clean conventional commit than it is to split one giant commit into many.
- **Local history**: Your local commit history can be messy. Use `git rebase -i` (interactive rebase) to clean up and organize your commits into well-structured conventional commits before pushing or opening a PR.
