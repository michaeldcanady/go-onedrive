# Architectural Improvement Plan: OpenTelemetry Tracing

## Goal
Improve system observability by instrumenting inter-service communication to facilitate debugging of latency and authentication failures.

## Plan
1. **Define Trace Points**: Identify critical service boundaries, especially in `internal/features/storage` and `internal/features/identity`.
2. **Setup Exporters**: Configure OTLP exporters (e.g., to stdout for local dev, or external collector) via `internal/core/telemetry`.
3. **Instrument**: Add span creation and propagation across middleware, particularly within the Kiota/Graph SDK stack.
4. **Context Propagation**: Ensure `context.Context` correctly carries trace headers across API boundaries.

## Verification
- Verify that spans are generated and exported for standard CLI operations (e.g., `ls`, `cp`).
- Confirm that traces correctly correlate across different service components.
