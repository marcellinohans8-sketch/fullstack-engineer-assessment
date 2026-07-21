# Fullstack Engineer Assessment

Task Management application menggunakan backend Go Gin, MySQL, Redis, dan frontend React Native TypeScript.

Project ini mengimplementasikan fitur filtering, update, soft delete, Redis caching, bug fixes, unit/component tests, migration SQL, dan dokumentasi setup.

## Tech Stack

- Backend: Go, Gin, GORM
- Database: MySQL
- Cache: Redis
- Frontend: React Native, Expo, TypeScript
- Testing backend: Go test, SQLite in-memory test database, Miniredis
- Testing frontend: Jest, Jest Expo, React Native Testing Library

## Project Structure

```text
.
├── backend
│   ├── config
│   │   ├── database.go
│   │   └── redis.go
│   ├── controllers
│   │   └── task_controller.go
│   ├── helpers
│   │   ├── cache.go
│   │   └── response.go
│   ├── migrations
│   │   └── 001_create_tasks.sql
│   ├── models
│   │   └── task.go
│   ├── routes
│   │   └── routes.go
│   ├── tests
│   │   ├── setup_test.go
│   │   └── task_test.go
│   ├── go.mod
│   └── main.go
├── frontend
│   ├── src
│   │   ├── api
│   │   │   └── taskApi.ts
│   │   ├── components
│   │   │   ├── __tests__
│   │   │   │   └── SearchBar.test.tsx
│   │   │   ├── EditTaskModal.tsx
│   │   │   ├── Loading.tsx
│   │   │   ├── Pagination.tsx
│   │   │   ├── SearchBar.tsx
│   │   │   ├── StatusFilter.tsx
│   │   │   └── TaskCard.tsx
│   │   ├── constants
│   │   │   └── api.ts
│   │   ├── screens
│   │   │   └── TaskScreen.tsx
│   │   └── types
│   │       └── task.ts
│   ├── package.json
│   └── package-lock.json
└── README.md
```

## Implemented Requirements

| Area | Requirement | Status |
| --- | --- | --- |
| Backend | Filtering by `status`, `keyword`, `assignee`, `page`, `limit`, `sort` | Done |
| Backend | `PUT /api/tasks/:id` | Done |
| Backend | Soft `DELETE /api/tasks/:id` | Done |
| Backend | Consistent error responses | Done |
| Redis | Cache `GET /api/tasks` for 60 seconds | Done |
| Redis | Cache key includes query parameters | Done |
| Redis | Invalidate cache after create/update/delete | Done |
| Frontend | Search input | Done |
| Frontend | Status filter | Done |
| Frontend | Pagination | Done |
| Frontend | Edit modal | Done |
| Frontend | Loading state | Done |
| Bug Fix | Duplicate title returns HTTP `409` | Done |
| Bug Fix | Refresh list after update | Done |
| Bug Fix | Hide soft-deleted tasks | Done |
| Testing | Backend update test | Done |
| Testing | Backend search/filter test | Done |
| Testing | Backend cache invalidation test | Done |
| Testing | Frontend component test | Done |
| Deliverable | README | Done |
| Deliverable | DB migration | Done |
| Deliverable | Unit/component tests | Done |

## Backend Setup

### 1. Prerequisites

Install and run:

- Go `1.26.x`
- MySQL
- Redis

### 2. Environment Variables

Create `backend/.env`:

```env
APP_PORT=8080
DB_USER=root
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=task_management
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
```

Notes:

- If `DB_HOST` is empty, backend defaults to `127.0.0.1`.
- If `DB_PORT` is empty, backend defaults to `3306`.
- If `REDIS_ADDR` is empty, backend defaults to `127.0.0.1:6379`.
- If Redis is unavailable, backend still starts and cache is disabled until Redis is available.

### 3. Install Backend Dependencies

```bash
cd backend
go mod download
```

### 4. Create Database

```sql
CREATE DATABASE task_management;
```

### 5. Run Migration

Migration file:

```text
backend/migrations/001_create_tasks.sql
```

Run:

```bash
mysql -u root -p task_management < backend/migrations/001_create_tasks.sql
```

The migration creates:

- `tasks` table
- unique index on `title`
- index on `deleted_at` for soft delete
- indexes on `status` and `assignee`

### 6. Start Backend

```bash
cd backend
go run .
```

Default base URL:

```text
http://localhost:8080
```

Health check:

```http
GET /
```

Expected response:

```json
{
  "message": "API Running"
}
```

## Backend API

Base API path:

```text
/api
```

### Response Format

Success response:

```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {}
}
```

List response:

```json
{
  "success": true,
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 0
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "validation error detail"
}
```

### Task Object

```json
{
  "id": 1,
  "title": "Build task filter",
  "description": "Add filter query parameters",
  "status": "todo",
  "assignee": "Hans",
  "createdAt": "2026-07-21T10:00:00Z",
  "updatedAt": "2026-07-21T10:00:00Z"
}
```

Allowed status values:

- `todo`
- `in_progress`
- `done`

### Create Task

```http
POST /api/tasks
Content-Type: application/json
```

Request:

```json
{
  "title": "Build task filter",
  "description": "Add keyword and status filters",
  "status": "todo",
  "assignee": "Hans"
}
```

Responses:

- `201 Created`: task created
- `400 Bad Request`: invalid request body
- `409 Conflict`: duplicate title
- `500 Internal Server Error`: unexpected database error

### List Tasks

```http
GET /api/tasks
```

Query parameters:

| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| `status` | string | empty | Filter by `todo`, `in_progress`, or `done` |
| `keyword` | string | empty | Search in `title` or `description` |
| `assignee` | string | empty | Filter by assignee |
| `page` | number | `1` | Page number, must be positive |
| `limit` | number | `10` | Page size, must be positive, capped at `100` |
| `sort` | string | `created_at desc` | Sort field and direction |

Example:

```http
GET /api/tasks?keyword=login&status=todo&assignee=Hans&page=1&limit=10&sort=created_at%20desc
```

Allowed sort values:

- `created_at asc`
- `created_at desc`
- `createdAt asc`
- `createdAt desc`
- `updated_at asc`
- `updated_at desc`
- `updatedAt asc`
- `updatedAt desc`
- `title asc`
- `title desc`
- `status asc`
- `status desc`
- `assignee asc`
- `assignee desc`

Responses:

- `200 OK`: list returned
- `400 Bad Request`: invalid query parameter
- `500 Internal Server Error`: unexpected database error

### Update Task

```http
PUT /api/tasks/:id
Content-Type: application/json
```

Request:

```json
{
  "title": "Build task filter",
  "description": "Add search, filter, pagination, and sorting",
  "status": "in_progress",
  "assignee": "Hans"
}
```

Responses:

- `200 OK`: task updated
- `400 Bad Request`: invalid id or request body
- `404 Not Found`: task not found or already soft-deleted
- `409 Conflict`: duplicate title
- `500 Internal Server Error`: unexpected database error

### Delete Task

```http
DELETE /api/tasks/:id
```

Behavior:

- Deletes task using GORM soft delete.
- Deleted task remains in the database with `deleted_at` filled.
- Deleted task is hidden from `GET /api/tasks`.

Responses:

- `200 OK`: task deleted
- `400 Bad Request`: invalid id
- `404 Not Found`: task not found or already soft-deleted
- `500 Internal Server Error`: unexpected database error

## Redis Cache Behavior

`GET /api/tasks` is cached for 60 seconds.

Cache key format:

```text
tasks:assignee=<value>&keyword=<value>&limit=<value>&page=<value>&sort=<value>&status=<value>
```

Examples:

```text
tasks:assignee=&keyword=&limit=10&page=1&sort=created_at+desc&status=
tasks:assignee=Hans&keyword=login&limit=5&page=1&sort=title+asc&status=todo
```

Invalidation:

- `POST /api/tasks` deletes all `tasks:*` cache entries.
- `PUT /api/tasks/:id` deletes all `tasks:*` cache entries.
- `DELETE /api/tasks/:id` deletes all `tasks:*` cache entries.

Cache helper uses Redis `SCAN` instead of `KEYS` for safer pattern deletion.

## Frontend Setup

### 1. Prerequisites

Install:

- Node.js
- npm
- Expo CLI through `npx expo`

### 2. Configure API URL

Update:

```text
frontend/src/constants/api.ts
```

Example for local Android emulator:

```ts
export const BASE_URL = "http://10.0.2.2:8080/api";
```

Example for physical device on same Wi-Fi:

```ts
export const BASE_URL = "http://YOUR_LOCAL_IP:8080/api";
```

Example for web/local browser:

```ts
export const BASE_URL = "http://localhost:8080/api";
```

### 3. Install Frontend Dependencies

```bash
cd frontend
npm install
```

### 4. Start Frontend

```bash
cd frontend
npm run start
```

Other scripts:

```bash
npm run android
npm run ios
npm run web
```

## Frontend Features

- Search input updates `keyword` query parameter.
- Status filter updates `status` query parameter.
- Pagination sends `page` and `limit` to backend.
- Edit modal calls `PUT /api/tasks/:id`.
- After successful edit, modal closes and task list refreshes.
- Loading state is displayed while the list request is running.

## Testing

### Backend Tests

Run:

```bash
cd backend
go test ./...
```

Covered tests:

- `TestUpdateTask`: verifies `PUT /api/tasks/:id`.
- `TestSearchTask`: verifies keyword/status/assignee filtering.
- `TestCacheInvalidation`: verifies Redis cache is created and invalidated after update/delete.
- `TestDuplicateTitleReturnsConflict`: verifies duplicate title returns `409`.
- `TestSoftDeletedTasksAreHidden`: verifies deleted tasks are hidden from list/search.

Backend tests use:

- SQLite in-memory database through `github.com/glebarez/sqlite`.
- Miniredis through `github.com/alicebob/miniredis/v2`.

This keeps tests independent from local MySQL and Redis services.

### Frontend Tests

Run:

```bash
cd frontend
npm run test
```

Covered tests:

- `SearchBar.test.tsx`: verifies search input calls `onChangeText` with typed keyword.

### TypeScript Check

```bash
cd frontend
npx tsc --noEmit
```

### Backend Build Check

```bash
cd backend
go build ./...
```

## Verification Results

The following commands were run successfully:

```bash
cd backend
go test ./...
go build ./...

cd ../frontend
npm run test -- --runInBand
npx tsc --noEmit
```

## Notes for Evaluator

- `PUT` and `DELETE` use `:id` in Gin routes, equivalent to the requested `/api/tasks/{id}` format.
- Soft delete is handled by GORM because the model includes `gorm.DeletedAt`.
- Normal GORM queries exclude soft-deleted rows by default, so `GET /api/tasks` hides deleted tasks.
- Duplicate title is enforced through a unique index and mapped to HTTP `409`.
- Redis is optional at boot to avoid crashing local development if Redis is temporarily unavailable, but caching is active when Redis is connected.
- Cache invalidation deletes all task list cache variants because every query combination can produce a different cached list.

## Troubleshooting

### Backend cannot connect to MySQL

Check:

- MySQL service is running.
- Database `task_management` exists.
- `backend/.env` contains correct `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, and `DB_NAME`.

### Backend starts but cache is disabled

Check:

- Redis service is running.
- `REDIS_ADDR` is correct, for example `127.0.0.1:6379`.
- `REDIS_PASSWORD` matches your Redis configuration.

### Frontend cannot reach backend

Check:

- Backend is running on port `8080`.
- `frontend/src/constants/api.ts` points to the correct host.
- Physical devices usually cannot use `localhost` for your laptop backend. Use your laptop local network IP.

### Duplicate title returns conflict

Expected behavior:

```http
HTTP/1.1 409 Conflict
```

The response body uses the standard error format:

```json
{
  "success": false,
  "message": "Task title already exists",
  "error": "title must be unique"
}
```
