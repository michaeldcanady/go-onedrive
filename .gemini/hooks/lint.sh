#!/usr/bin/env bash
# .gemini/hooks/lint.sh

# Read tool input from stdin
INPUT=$(cat)

# Extract file path using jq
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path' 2>/dev/null)

# Only lint Go files
if [[ "$FILE_PATH" == *.go ]]; then
    # 1. Skip generated code
    if grep -q "Code generated" "$FILE_PATH" 2>/dev/null; then
        echo '{"decision": "allow"}'
        exit 0
    fi

    # 2. Package-level Linting via Justfile
    # Using 'just' ensures we use the same config as CI
    PKG_DIR=$(dirname "$FILE_PATH")
    LINT_OUTPUT=$(just lint-pkg "$PKG_DIR" 2>&1)
    LINT_EXIT=$?

    if [ $LINT_EXIT -ne 0 ]; then
        # Return 'deny' to feed errors back to the agent
        cat <<EOF
{
  "decision": "deny",
  "reason": "Linting failed for package $PKG_DIR (triggered by $FILE_PATH). Architectural violations or syntax errors found:\n$LINT_OUTPUT",
  "systemMessage": "⚠️ Architectural Guardrails Violated. Please ensure 'Interfaces In, Real Types Out' and check feature boundaries."
}
EOF
        exit 0
    fi
fi

# Default to allow
echo '{"decision": "allow"}'
exit 0
