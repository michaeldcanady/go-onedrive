---
name: performance-test-writer
description: Generates performance and benchmark tests for the CLI and its underlying services. These tests measure execution time, memory usage, and throughput to ensure the application remains efficient under load.
---

# Performance Test Writer

## Overview
Performance tests (Benchmarks) identify bottlenecks and quantify the efficiency of critical code paths. In Go, these are implemented using the `testing.B` framework.

## Workflow
1. **Identify Critical Path**: Determine which operations are likely to be slow (e.g., large file transfers, complex URI parsing).
2. **Setup Benchmark**:
    - Initialize required state outside the benchmark loop.
    - Use `b.ResetTimer()` to exclude setup time.
3. **Implement Loop**: Execute the operation `b.N` times.
4. **Configure Flags**: Use `-benchmem` to track memory allocations.
5. **Analyze Results**: Compare results against baselines or previous versions.
6. **Verify**: Run `go test -v -bench=. ./tests/performance/...`

## Patterns

### Benchmark Template
```go
func BenchmarkOperation(b *testing.B) {
    // Setup
    data := prepareLargeDataset()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = CriticalOperation(data)
    }
}
```

## Guidelines
- **Real-World Scenarios**: Use data sizes that reflect actual usage (e.g., 1GB file copy, 10,000 item listing).
- **Consistency**: Run benchmarks in a controlled environment to minimize noise.
- **Optimization**: Focus performance efforts on code paths that the benchmarks prove are slow.
