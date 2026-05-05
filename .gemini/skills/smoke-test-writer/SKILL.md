---
name: smoke-test-writer
description: Generates rapid smoke tests for the CLI. These tests are designed to run in seconds after a build or deployment to ensure that the application's most basic and vital functions are operational.
---

# Smoke Test Writer

## Overview
Smoke tests are a subset of E2E tests focused on speed and high-level health checks. They answer the question: "Does the application start and perform its core function without crashing?"

## Workflow
1. **Identify Vital Signs**: Select 3-5 critical commands (e.g., `version`, `help`, `config list`).
2. **Setup Script**: Create a lightweight shell script or Go test file.
3. **Define Pass/Fail**:
    - Exit code must be 0.
    - Output must contain key phrases (e.g., "Usage:").
    - Execution time must be below a threshold (e.g., < 2 seconds).
4. **Implement**: Keep it extremely simple with no complex mocking.
5. **Verify**: Run `./tests/smoke/run.sh` or `go test ./tests/smoke/...`

## Patterns

### Smoke Test Script (Bash)
```bash
#!/bin/bash
set -e

./odc version
./odc help
./odc config list --session
```

## Guidelines
- **Speed**: Smoke tests should never take more than a few seconds.
- **Reliability**: They should only fail if there is a major problem.
- **Automation**: Designed to be the first step in a CI/CD pipeline.
