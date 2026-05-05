# 2. Use bbolt for local persistent state

Date: 2025-05-14

## Status

Status: Accepted

## Context

The `odc` CLI tool needs to persist user profiles, session tokens, cached drive data, and configuration state locally. This data must survive application restarts and be available across different CLI invocations. Users need a solution that's lightweight, reliable, and requires no external server installation

## Decision

Users use **bbolt** (`go.etcd.io/bbolt`) as the primary embedded key-value store for local persistence

## Consequences

## Benefits
- **Zero-Dependency:** `bbolt` is a pure Go implementation and requires no CGO or external database server
- **ACID Compliant:** Ensures data integrity with fully ACID transactions
- **Single File:** The entire database is a single file on disk, making it easy to manage and back up
- **Performance:** Highly performant for read-heavy workloads, which is typical for a CLI tool reading configuration and credentials
- **Stability:** `bbolt` (a fork of BoltDB) is widely used and well-tested in the Go ecosystem (for example, in etcd)

## Trade-offs
- **Single Writer:** Only one process can have a write transaction open at a time. This is  acceptable for a CLI tool which is typically run as a single instance per user
- **No Native Migrations:** Managing schema changes in a key-value store requires manual migration logic
- **Limited Querying:** Application code must implement complex queries, unlike a relational database. Developers iterate over keys or use secondary indexes.

## Links

- [bbolt GitHub Repository](https://github.com/etcd-io/bbolt)
