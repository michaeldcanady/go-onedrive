# Conventional Commit Examples

## Basic Feature
```text
feat(mount): add support for drive_id option
```

## Bug Fix
```text
fix(identity): resolve token refresh race condition
```

## Documentation
```text
docs: update README with mount command examples
```

## Breaking Change
```text
feat(api): change user endpoint response format

BREAKING CHANGE: the user endpoint now returns a nested object instead of a flat list.
```

## Multiple Scopes (Avoid if possible, but allowed)
```text
refactor(mount,identity): cleanup cross-cutting dependency
```

## Revert
```text
revert: feat(mount): add support for drive_id option

This reverts commit 12345678.
```
