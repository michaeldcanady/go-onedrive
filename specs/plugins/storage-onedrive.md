# Plugin Specification: `storage-onedrive`

## Overview
The `storage-onedrive` plugin provides storage operations for Microsoft OneDrive and SharePoint. Requests require the drive id and a valid auth token.

## Capabilities
- **File Operations:** CRUD operations for files and folders (ls, cp, mv, rm, mkdir, stat).
- **Streaming:** Supports chunked uploads and downloads for large files using the Microsoft Graph API's upload session.
- **Drive Discovery:** Lists available drives associated with the account (Personal, OneDrive for Business, SharePoint Sites).
- **Metadata:** Retrieves item metadata including size, hashes (SHA1, QuickXorHash), and timestamps.

## Configuration Options
The following options can be set via `odc config set storage.onedrive.<key> <value>`:
- `chunk_size`: Size of upload chunks in bytes (default 10MB).
- `concurrent_transfers`: Number of concurrent file transfers (default 3).

## Interface
Implements the `StorageService` gRPC interface as defined in `specs/proto/storage.proto`.

## Behavior
- **Authentication**: Uses the `token` string provided in the `options` map of every gRPC request. 
- **Targeting**: Uses the `drive_id` string from the `options` map to target specific drives (defaults to `root` for personal drives).
- **Path Mapping**: Maps VFS paths to Graph API endpoints using the `root:/path` or `drives/{id}/items/root:/path` addressing schemes.
- **I/O Handling**: Supports chunked transfers. For simple Read/Write, it handles full `[]byte` payloads for small files.
- **Throttling**: Handles API rate limiting with exponential backoff.
