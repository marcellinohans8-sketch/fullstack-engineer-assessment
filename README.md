# Aplikasi Manajemen Tugas

Repositori ini berisi implementasi assessment Fullstack Engineer untuk aplikasi manajemen tugas. Aplikasi menggunakan backend Go dengan Gin, database MySQL, Redis untuk cache, dan frontend React Native dengan TypeScript.

## Ringkasan Fitur

| Area | Kebutuhan | Status |
| --- | --- | --- |
| Backend | Filter daftar tugas berdasarkan `status`, `keyword`, `assignee`, `page`, `limit`, dan `sort` | Selesai |
| Backend | Endpoint `PUT /api/tasks/:id` untuk mengubah tugas | Selesai |
| Backend | Endpoint `DELETE /api/tasks/:id` dengan soft delete | Selesai |
| Backend | Format respons error konsisten | Selesai |
| Redis | Cache `GET /api/tasks` selama 60 detik | Selesai |
| Redis | Key cache menyertakan query parameter | Selesai |
| Redis | Cache dihapus setelah tambah, ubah, dan hapus tugas | Selesai |
| Frontend | Input pencarian | Selesai |
| Frontend | Filter status | Selesai |
| Frontend | Pagination | Selesai |
| Frontend | Modal edit tugas | Selesai |
| Frontend | Indikator memuat | Selesai |
| Bug fix | Judul duplikat mengembalikan HTTP `409` | Selesai |
| Bug fix | Daftar tugas refresh setelah update | Selesai |
| Bug fix | Tugas yang sudah soft delete tidak tampil | Selesai |
| Testing | Test backend untuk update, search, dan cache invalidation | Selesai |
| Testing | Test komponen frontend untuk search | Selesai |
| Deliverable | README, migration database, dan unit/component test | Selesai |

## Teknologi

- Backend: Go, Gin, GORM
- Database: MySQL
- Cache: Redis
- Frontend: React Native, Expo, TypeScript
- Pengujian backend: Go test, SQLite di memori, Miniredis
- Pengujian frontend: Jest, Jest Expo, React Native Testing Library

## Struktur Folder

```text
.
|-- backend
|   |-- config
|   |   |-- database.go
|   |   `-- redis.go
|   |-- controllers
|   |   `-- task_controller.go
|   |-- helpers
|   |   |-- cache.go
|   |   `-- response.go
|   |-- migrations
|   |   `-- 001_create_tasks.sql
|   |-- models
|   |   `-- task.go
|   |-- routes
|   |   `-- routes.go
|   |-- tests
|   |   |-- setup_test.go
|   |   `-- task_test.go
|   |-- go.mod
|   `-- main.go
|-- frontend
|   |-- src
|   |   |-- api
|   |   |   `-- taskApi.ts
|   |   |-- components
|   |   |   |-- __tests__
|   |   |   |   `-- SearchBar.test.tsx
|   |   |   |-- EditTaskModal.tsx
|   |   |   |-- Loading.tsx
|   |   |   |-- Pagination.tsx
|   |   |   |-- SearchBar.tsx
|   |   |   |-- StatusFilter.tsx
|   |   |   `-- TaskCard.tsx
|   |   |-- constants
|   |   |   `-- api.ts
|   |   |-- screens
|   |   |   `-- TaskScreen.tsx
|   |   `-- types
|   |       `-- task.ts
|   |-- package.json
|   `-- package-lock.json
`-- README.md
```

## Setup Backend

### 1. Prasyarat

Pastikan sudah terpasang dan berjalan:

- Go
- MySQL
- Redis

### 2. Konfigurasi Environment

Buat file `backend/.env`:

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

Catatan:

- Jika `DB_HOST` kosong, backend memakai default `127.0.0.1`.
- Jika `DB_PORT` kosong, backend memakai default `3306`.
- Jika `REDIS_ADDR` kosong, backend memakai default `127.0.0.1:6379`.
- Jika Redis tidak tersedia saat aplikasi dijalankan, backend tetap hidup dan cache dinonaktifkan sementara.

### 3. Instal Dependensi Backend

```bash
cd backend
go mod download
```

### 4. Buat Database

```sql
CREATE DATABASE task_management;
```

### 5. Jalankan Migrasi

File migrasi tersedia di:

```text
backend/migrations/001_create_tasks.sql
```

Jalankan migrasi:

```bash
mysql -u root -p task_management < backend/migrations/001_create_tasks.sql
```

Migrasi membuat:

- Tabel `tasks`
- Unique index pada kolom `title`
- Index pada kolom `deleted_at` untuk soft delete
- Index pada kolom `status` dan `assignee`

### 6. Jalankan Backend

```bash
cd backend
go run .
```

Backend berjalan di:

```text
http://localhost:8080
```

Endpoint pengecekan:

```http
GET /
```

Contoh respons:

```json
{
  "message": "API Running"
}
```

## API Backend

Path dasar API:

```text
/api
```

### Format Respons

Respons sukses untuk tambah, ubah, dan hapus:

```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {}
}
```

Respons daftar tugas:

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

Respons error:

```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "detail error validasi"
}
```

### Struktur Data Tugas

```json
{
  "id": 1,
  "title": "Membuat filter tugas",
  "description": "Menambahkan filter keyword dan status",
  "status": "todo",
  "assignee": "Hans",
  "createdAt": "2026-07-21T10:00:00Z",
  "updatedAt": "2026-07-21T10:00:00Z"
}
```

Nilai `status` yang diperbolehkan:

- `todo`
- `in_progress`
- `done`

### Tambah Tugas

```http
POST /api/tasks
Content-Type: application/json
```

Isi permintaan:

```json
{
  "title": "Membuat filter tugas",
  "description": "Menambahkan filter keyword dan status",
  "status": "todo",
  "assignee": "Hans"
}
```

Kemungkinan respons:

- `201 Created`: tugas berhasil dibuat
- `400 Bad Request`: body tidak valid
- `409 Conflict`: judul tugas sudah ada
- `500 Internal Server Error`: error database tidak terduga

### Ambil Daftar Tugas

```http
GET /api/tasks
```

Parameter query:

| Parameter | Tipe | Default | Keterangan |
| --- | --- | --- | --- |
| `status` | string | kosong | Filter berdasarkan `todo`, `in_progress`, atau `done` |
| `keyword` | string | kosong | Pencarian pada `title` dan `description` |
| `assignee` | string | kosong | Filter berdasarkan penanggung jawab |
| `page` | number | `1` | Nomor halaman, harus lebih dari `0` |
| `limit` | number | `10` | Jumlah data per halaman, maksimal `100` |
| `sort` | string | `created_at desc` | Field dan arah pengurutan |

Contoh:

```http
GET /api/tasks?keyword=login&status=todo&assignee=Hans&page=1&limit=10&sort=created_at%20desc
```

Nilai `sort` yang diperbolehkan:

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

Kemungkinan respons:

- `200 OK`: daftar tugas berhasil dikembalikan
- `400 Bad Request`: query parameter tidak valid
- `500 Internal Server Error`: error database tidak terduga

### Ubah Tugas

```http
PUT /api/tasks/:id
Content-Type: application/json
```

Isi permintaan:

```json
{
  "title": "Membuat filter tugas",
  "description": "Menambahkan search, filter, pagination, dan sorting",
  "status": "in_progress",
  "assignee": "Hans"
}
```

Kemungkinan respons:

- `200 OK`: tugas berhasil diubah
- `400 Bad Request`: id atau body tidak valid
- `404 Not Found`: tugas tidak ditemukan atau sudah dihapus
- `409 Conflict`: judul tugas sudah ada
- `500 Internal Server Error`: error database tidak terduga

### Hapus Tugas

```http
DELETE /api/tasks/:id
```

Perilaku:

- Endpoint menggunakan hapus lunak (soft delete) dari GORM.
- Data tetap ada di database dengan kolom `deleted_at` terisi.
- Data yang sudah dihapus tidak tampil dari `GET /api/tasks`.

Kemungkinan respons:

- `200 OK`: tugas berhasil dihapus
- `400 Bad Request`: id tidak valid
- `404 Not Found`: tugas tidak ditemukan atau sudah dihapus
- `500 Internal Server Error`: error database tidak terduga

## Perilaku Cache Redis

Endpoint `GET /api/tasks` disimpan di Redis selama 60 detik.

Format key cache:

```text
tasks:assignee=<nilai>&keyword=<nilai>&limit=<nilai>&page=<nilai>&sort=<nilai>&status=<nilai>
```

Contoh:

```text
tasks:assignee=&keyword=&limit=10&page=1&sort=created_at+desc&status=
tasks:assignee=Hans&keyword=login&limit=5&page=1&sort=title+asc&status=todo
```

Invalidasi cache:

- `POST /api/tasks` menghapus semua cache `tasks:*`.
- `PUT /api/tasks/:id` menghapus semua cache `tasks:*`.
- `DELETE /api/tasks/:id` menghapus semua cache `tasks:*`.

Helper cache menggunakan Redis `SCAN`, bukan `KEYS`, agar penghapusan pattern lebih aman untuk jumlah key yang banyak.

## Setup Frontend

### 1. Prasyarat

Pastikan sudah terpasang:

- Node.js
- npm
- Expo melalui `npx expo`

### 2. Konfigurasi URL API

Ubah file:

```text
frontend/src/constants/api.ts
```

Contoh untuk Android emulator:

```ts
export const BASE_URL = "http://10.0.2.2:8080/api";
```

Contoh untuk perangkat fisik di jaringan Wi-Fi yang sama:

```ts
export const BASE_URL = "http://IP_LOKAL_LAPTOP:8080/api";
```

Contoh untuk browser lokal:

```ts
export const BASE_URL = "http://localhost:8080/api";
```

### 3. Instal Dependensi Frontend

```bash
cd frontend
npm install
```

### 4. Jalankan Frontend

```bash
cd frontend
npm run start
```

Skrip lain:

```bash
npm run android
npm run ios
npm run web
```

## Fitur Frontend

- Input pencarian mengirim parameter `keyword`.
- Filter status mengirim parameter `status`.
- Pagination mengirim parameter `page` dan `limit`.
- Modal edit memanggil `PUT /api/tasks/:id`.
- Setelah edit berhasil, modal ditutup dan daftar tugas dimuat ulang.
- Indikator memuat tampil selama permintaan daftar tugas berjalan.

## Pengujian

### Pengujian Backend

Jalankan:

```bash
cd backend
go test ./...
```

Cakupan pengujian backend:

- `TestUpdateTask`: memastikan `PUT /api/tasks/:id` berhasil mengubah data.
- `TestSearchTask`: memastikan filter `keyword`, `status`, dan `assignee` berjalan.
- `TestCacheInvalidation`: memastikan cache dibuat dan dihapus setelah create, update, dan delete.
- `TestDuplicateTitleReturnsConflict`: memastikan judul duplikat mengembalikan `409`.
- `TestSoftDeletedTasksAreHidden`: memastikan tugas soft-deleted tidak muncul pada list/search.

Pengujian backend memakai SQLite di memori dan Miniredis, sehingga tidak bergantung pada layanan MySQL dan Redis lokal.

### Pengujian Frontend

Jalankan:

```bash
cd frontend
npm run test
```

Cakupan pengujian frontend:

- `SearchBar.test.tsx`: memastikan input pencarian memanggil `onChangeText` dengan keyword yang diketik.

### Pengecekan TypeScript

```bash
cd frontend
npx tsc --noEmit
```

### Pengecekan Kompilasi Backend

```bash
cd backend
go build ./...
```

## Hasil Verifikasi

Perintah berikut sudah dijalankan dan berhasil:

```bash
cd backend
go test ./...
go build ./...

cd ../frontend
npm run test -- --runInBand
npx tsc --noEmit
```

## Catatan untuk Evaluator

- Route Gin memakai format `:id`, setara dengan kebutuhan `/api/tasks/{id}`.
- Hapus lunak bekerja karena model `Task` memiliki field `gorm.DeletedAt`.
- Query GORM biasa otomatis menyembunyikan data dengan `deleted_at` terisi.
- Judul duplikat dijaga oleh unique index dan dipetakan menjadi HTTP `409 Conflict`.
- Redis bersifat opsional saat aplikasi start agar development lokal tidak gagal ketika Redis belum hidup.
- Saat Redis tersedia, cache list tugas aktif selama 60 detik.
- Semua variasi cache list tugas dihapus setelah create, update, dan delete karena query parameter dapat menghasilkan daftar yang berbeda.

## Pemecahan Masalah

### Backend tidak bisa terhubung ke MySQL

Periksa:

- Layanan MySQL sudah berjalan.
- Database `task_management` sudah dibuat.
- File `backend/.env` berisi `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, dan `DB_NAME` yang benar.

### Backend berjalan tetapi cache tidak aktif

Periksa:

- Layanan Redis sudah berjalan.
- `REDIS_ADDR` benar, misalnya `127.0.0.1:6379`.
- `REDIS_PASSWORD` sesuai konfigurasi Redis.

### Frontend tidak bisa mengakses backend

Periksa:

- Backend berjalan di port `8080`.
- `frontend/src/constants/api.ts` mengarah ke host yang benar.
- Perangkat fisik biasanya tidak bisa memakai `localhost` untuk mengakses backend di laptop. Gunakan IP lokal laptop.

### Judul duplikat mengembalikan konflik

Ini adalah perilaku yang diharapkan:

```http
HTTP/1.1 409 Conflict
```

Contoh respons:

```json
{
  "success": false,
  "message": "Task title already exists",
  "error": "title must be unique"
}
```
