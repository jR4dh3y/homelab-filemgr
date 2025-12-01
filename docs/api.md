# API Reference

The Homelab File Manager provides a REST API for file operations and a WebSocket endpoint for real-time updates.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All endpoints except `/auth/*` and `/health` require JWT authentication.

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin"
}
```

**Response:**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 3600
}
```

### Using Tokens

Include the access token in the Authorization header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### Refresh Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

## File Operations

### List Mount Points

```http
GET /api/v1/files
```

**Response:**
```json
{
  "roots": [
    { "name": "media", "readOnly": false },
    { "name": "documents", "readOnly": false },
    { "name": "backups", "readOnly": true }
  ]
}
```

### List Directory

```http
GET /api/v1/files/{path}?page=1&pageSize=50&sortBy=name&sortDir=asc
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Page number |
| pageSize | int | 50 | Items per page |
| sortBy | string | name | Sort field: name, size, modTime, type |
| sortDir | string | asc | Sort direction: asc, desc |
| filter | string | | Filter by name (contains) |

**Response:**
```json
{
  "path": "media/movies",
  "items": [
    {
      "name": "movie.mkv",
      "path": "media/movies/movie.mkv",
      "size": 4294967296,
      "isDir": false,
      "modTime": "2024-01-15T10:30:00Z",
      "permissions": "-rw-r--r--",
      "mimeType": "video/x-matroska"
    }
  ],
  "totalCount": 150,
  "page": 1,
  "pageSize": 50
}
```

### Get File Info

```http
GET /api/v1/files/{path}
```

Returns file metadata for a single file.

### Create Directory

```http
POST /api/v1/files/{path}
Content-Type: application/json

{
  "name": "new-folder"
}
```

**Response:**
```json
{
  "name": "new-folder",
  "path": "media/new-folder",
  "size": 0,
  "isDir": true,
  "modTime": "2024-01-15T10:30:00Z",
  "permissions": "drwxr-xr-x"
}
```

### Rename/Move

```http
PUT /api/v1/files/{path}
Content-Type: application/json

{
  "newPath": "media/renamed-file.txt"
}
```

### Delete

```http
DELETE /api/v1/files/{path}?confirm=true
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| confirm | bool | Required for directories |

---

## Streaming

### Download File

```http
GET /api/v1/stream/download/{path}
```

Supports HTTP Range requests for resumable downloads:

```http
GET /api/v1/stream/download/{path}
Range: bytes=0-1023
```

**Response Headers:**
```
Content-Type: application/octet-stream
Content-Disposition: attachment; filename="file.txt"
Accept-Ranges: bytes
Content-Length: 1024
```

### Chunked Upload

Upload files in chunks for large file support and resumability.

```http
POST /api/v1/stream/upload/{path}
Content-Type: application/octet-stream
X-Upload-ID: unique-upload-id
X-Chunk-Index: 0
X-Total-Chunks: 10
X-Checksum: sha256:abc123...  (optional, on last chunk)

[binary chunk data]
```

**Headers:**
| Header | Required | Description |
|--------|----------|-------------|
| X-Upload-ID | Yes | Unique identifier for this upload |
| X-Chunk-Index | Yes | Zero-based chunk index |
| X-Total-Chunks | Yes | Total number of chunks |
| X-Checksum | No | SHA256 checksum (final chunk only) |

**Response:**
```json
{
  "uploadId": "unique-upload-id",
  "chunkIndex": 0,
  "received": true
}
```

**Final chunk response:**
```json
{
  "uploadId": "unique-upload-id",
  "complete": true,
  "path": "media/uploaded-file.zip",
  "size": 104857600,
  "checksum": "sha256:abc123..."
}
```

---

## Search

### Search Files

```http
GET /api/v1/search?path=media&q=movie
```

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| path | string | Yes | Directory to search in |
| q | string | Yes | Search query (case-insensitive) |

**Response:**
```json
{
  "path": "media",
  "query": "movie",
  "results": [
    {
      "name": "movie.mkv",
      "path": "media/movies/movie.mkv",
      "size": 4294967296,
      "isDir": false,
      "modTime": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

## Background Jobs

### List Jobs

```http
GET /api/v1/jobs
```

**Response:**
```json
{
  "jobs": [
    {
      "id": "job_abc123",
      "type": "copy",
      "state": "running",
      "progress": 45,
      "sourcePath": "media/movie.mkv",
      "destPath": "backups/movie.mkv",
      "createdAt": "2024-01-15T10:30:00Z",
      "startedAt": "2024-01-15T10:30:01Z"
    }
  ]
}
```

### Get Job Status

```http
GET /api/v1/jobs/{id}
```

### Create Job

```http
POST /api/v1/jobs
Content-Type: application/json

{
  "type": "copy",
  "sourcePath": "media/movie.mkv",
  "destPath": "backups/movie.mkv"
}
```

**Job Types:**
| Type | Description |
|------|-------------|
| copy | Copy file/directory |
| move | Move file/directory |
| delete | Delete file/directory |

**Response:**
```json
{
  "id": "job_abc123",
  "type": "copy",
  "state": "pending",
  "progress": 0,
  "sourcePath": "media/movie.mkv",
  "destPath": "backups/movie.mkv",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

### Cancel Job

```http
DELETE /api/v1/jobs/{id}
```

---

## WebSocket

### Connect

```
ws://localhost:8080/api/v1/ws?token=eyJhbGciOiJIUzI1NiIs...
```

Or with header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### Client Messages

**Subscribe to job updates:**
```json
{
  "type": "subscribe",
  "jobId": "job_abc123"
}
```

**Unsubscribe:**
```json
{
  "type": "unsubscribe",
  "jobId": "job_abc123"
}
```

**Ping:**
```json
{
  "type": "ping"
}
```

### Server Messages

**Job update:**
```json
{
  "type": "job_update",
  "payload": {
    "jobId": "job_abc123",
    "state": "running",
    "progress": 45
  }
}
```

**Job complete:**
```json
{
  "type": "job_complete",
  "payload": {
    "jobId": "job_abc123",
    "state": "completed",
    "progress": 100
  }
}
```

**Error:**
```json
{
  "type": "error",
  "payload": {
    "message": "Job not found"
  }
}
```

**Pong:**
```json
{
  "type": "pong"
}
```

---

## Health Check

```http
GET /api/v1/health
```

**Response:**
```json
{
  "status": "ok"
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Access denied (mount point, read-only) |
| 404 | Not Found - Path does not exist |
| 409 | Conflict - File already exists |
| 500 | Internal Server Error |

### Error Codes

| Code | Description |
|------|-------------|
| PATH_NOT_FOUND | Requested path does not exist |
| PATH_TRAVERSAL | Path traversal attempt detected |
| MOUNT_DENIED | Access outside mount points |
| READ_ONLY | Write operation on read-only mount |
| INVALID_TOKEN | JWT token is invalid |
| TOKEN_EXPIRED | JWT token has expired |
