# Lifeline — Backend Documentation

The backend is a **Go + Fiber** REST API that powers the Lifeline blood-donation platform. It handles authentication, user management, blood-bank stock, geo data (Moroccan cities), blood requests, notifications, and admin moderation.

---

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Getting Started](#getting-started)
3. [Environment Variables](#environment-variables)
4. [Project Structure](#project-structure)
5. [Architecture](#architecture)
6. [Database Models](#database-models)
7. [Seeding](#seeding)
8. [Authentication & Authorization](#authentication--authorization)
9. [API Reference](#api-reference)
10. [Middleware](#middleware)
11. [Services](#services)
12. [Known Issues & TODOs](#known-issues--todos)

---

## Tech Stack

| Concern | Choice |
|---|---|
| Language | Go 1.24 |
| HTTP framework | [Fiber v2](https://github.com/gofiber/fiber) (`gofiber/fiber/v2`) |
| ORM | [GORM](https://gorm.io) (`gorm.io/gorm`) |
| Database | **SQLite** via `gorm.io/driver/sqlite` |
| Auth | JWT — `golang-jwt/jwt/v5` + `gofiber/contrib/jwt` |
| Password hashing | `golang.org/x/crypto/bcrypt` |
| Media upload | Cloudinary (`cloudinary-go/v2`) — *partially wired* |
| UUIDs | `github.com/google/uuid` |

> **Note on PostGIS:** the README and models mention PostGIS/geo features, but the current implementation is **SQLite** (single file `lifeline.db`). Latitude/longitude are stored as plain `float64` columns. There is no PostGIS dependency in `go.mod` yet.

---

## Getting Started

### Prerequisites
- Go 1.24+

### Run

```bash
cd backend
go run cmd/api/main.go
```

The server listens on `http://localhost:8080` (override with the `PORT` env var).

On startup, `main.go` runs, in order:

1. `database.Connect()` — opens SQLite, auto-migrates all models, enables foreign keys.
2. `database.SeedMoroccanCities()` — imports regions & cities from JSON.
3. `database.SeedRolesAndAdmin()` — creates roles + a default admin.
4. `database.SeedBanks()` — imports blood banks/hospitals.
5. `database.SeedFakeUsers()` — creates demo donors.
6. Registers middleware (logger, CORS, activity logger) and routes.
7. Serves `/uploads` as static files from `./backend/uploads`.

> ⚠️ The SQLite path in `pkg/database/database.go` is hardcoded to `../../lifeline.db`, which resolves correctly **only when run from the `backend/` directory** (i.e. `go run cmd/api/main.go`).

---

## Environment Variables

| Variable | Used by | Default / Fallback |
|---|---|---|
| `PORT` | `main.go` (server port) | `8080` |
| `TOKEN_KEY` | JWT signing (`token_service.go`, `middleware/activity.go`) | `super_secret_default_key_change_me` |
| `CLOUDINARY_URL` | Cloudinary upload (`photo_service.go`) | *(not configured → uploads no-op)* |

> `.env.example` also lists `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT` (PostgreSQL-style). These are **not read by the current code** — the app uses SQLite. They are aspirational for a future PostGIS migration.

---

## Project Structure

```
backend/
├── cmd/api/main.go                # Entry point: DB connect, seed, middleware, routes
├── internal/
│   ├── dto/                       # Request/response shapes (Data Transfer Objects)
│   │   ├── auth_dto.go            #   RegisterInput, LoginInput, UserResponse
│   │   ├── request_dto.go         #   CreateRequestInput
│   │   ├── bank_dto.go            #   BankRegisterInput
│   │   ├── bank_response.go       #   BankResponse, BloodGroupDto
│   │   ├── profile_dto.go         #   UserProfile
│   │   ├── user_dto.go            #   MemberResponse
│   │   ├── params.go              #   UserParams, PaginationHeader
│   │   ├── moderator_dto.go
│   │   ├── sponsor_dto.go         #   SponsorStatsResponse
│   │   └── stock_dto.go           #   BloodGroupUpdateInput
│   ├── handlers/                  # HTTP handlers (thin, orchestrate services/repos)
│   │   ├── auth_handler.go        #   Register, Login, GetUserProfile
│   │   ├── user_handler.go        #   GetUsers, GetUser
│   │   ├── request_handler.go     #   Create/Get/Accept blood requests + admin status
│   │   ├── bank_handler.go        #   GetBanks (filters + pagination)
│   │   ├── geo_handler.go         #   GetCities
│   │   ├── notification_handler.go#   GetMyNotifications, MarkAsRead
│   │   ├── admin_handler.go       #   GetUsersWithRoles, EditRoles
│   │   ├── moderate_handler.go    #   (stub) bank moderators
│   │   └── sponsor.go             #   (stub) GetImpactReport
│   ├── middleware/
│   │   ├── activity.go            #   Protected (JWT), RequireRole, LogUserActivity
│   │   └── error_handler.go       #   Global error handler
│   ├── models/                    # GORM entities
│   ├── repository/
│   │   └── user_repository.go     #   DB access for users (pagination/filtering)
│   ├── routes/routes.go           #   Single place where routes are registered
│   └── services/
│       ├── token_service.go       #   JWT creation
│       └── photo_service.go       #   Cloudinary upload/delete
├── pkg/database/
│   ├── database.go                #   Connect() + AutoMigrate()
│   └── seed.go                    #   All seeders
├── assets/                        #   Seed data (banks.json, cities JSON)
├── uploads/                       #   Locally-stored user photos
├── go.mod / go.sum
└── .env / .env.example
```

---

## Architecture

The project follows a **layered architecture** (Repository → Service → Handler):

```
HTTP Request
    │
    ▼
Middleware (CORS → Logger → Activity → Protected/RequireRole)
    │
    ▼
Handler  ──► Service ──► Repository ──► GORM ──► SQLite
    │            │             │
    └── DTO mapping (JSON responses)
```

- **Handlers** parse input, call services/repositories, and return JSON. Some handlers (e.g. `request_handler`, `admin_handler`) bypass the service/repository layer and query `database.DB` directly.
- **Services** hold business logic (token signing, photo upload). Defined as Go **interfaces** for testability.
- **Repository** (`user_repository.go`) encapsulates user queries (filtering, pagination).
- **DTOs** decouple the wire format from the database models.
- **Models** are GORM structs that map to SQLite tables.

Dependencies are wired manually in `routes.go` (no DI framework).

---

## Database Models

GORM `AutoMigrate` creates tables for all of these on startup.

| Model | Table | Key fields / notes |
|---|---|---|
| `User` | `users` | `username` (unique), `email` (unique), `password_hash`, `phone_number`, `name`, `date_of_birth`, `gender`, `last_active`, `blood_group`, `available`, `profile` (local photo path) |
| `Role` | `roles` | `name` (unique): `Admin`, `Sponsor`, `Donor`, `Moderator` |
| `UserRole` | `user_roles` | composite PK (`user_id`, `role_id`) — many-to-many join |
| `Address` | `addresses` | `area`, `city`, `state`, `country`, `postal_code`, `latitude`, `longitude`; nullable `user_id` / `bank_id` |
| `Photo` | `photos` | `public_id`, `url`; nullable `user_id` / `bank_id` |
| `Bank` | `banks` | `name`, `city`, `phone_number`, `email`, `website`; has-many `BloodGroup`, one `Address`, one `Photo`, many `Moderator` |
| `BloodGroup` | `blood_groups` | `group` (e.g. `O+`), `quantity`, `bank_id` |
| `Moderator` | `moderators` | composite PK (`user_id`, `bank_id`), `type` |
| `Sponsor` | `sponsors` | `name` (unique), `logo_url`, `website`, `total_paid`, `campaign_budget`, `is_active`, `views_count`, `clicks_count` |
| `Region` | `regions` | `region`; has-many `City` |
| `City` | `cities` | `ville`, `region_id` |
| `BloodRequest` | `blood_requests` | `user_id`, `blood_type`, `is_urgent`, `hospital_name`, `latitude`, `longitude`, `prescription_photo`, `status` |
| `Notification` | `notifications` | `user_id`, `title`, `message`, `type`, `is_read` |

### Blood request statuses

Declared in `models/request.go`:

```go
StatusPending   // "Pending"
StatusApproved  // "Approved"
StatusFulfilled // "Fulfilled"
StatusCancelled // "Cancelled"
```

> ⚠️ Two extra status strings are used in code but **not** part of the enum:
> - `"In_Progress"` — set by `AcceptRequest` when a donor accepts a request.
> - `"Rejected"` — expected by the frontend admin panel (but the backend never sets it).

### Blood groups

Canonical set (used by the bank seeder and UI): `A+`, `A-`, `B+`, `B-`, `AB+`, `AB-`, `O+`, `O-`.

---

## Seeding

All seeders are idempotent (skip if data already exists).

| Seeder | Source | What it creates |
|---|---|---|
| `SeedMoroccanCities` | `assets/sql-moroccan-cities/json/region.json` + `ville.json` | 12 Moroccan regions + all cities |
| `SeedRolesAndAdmin` | hardcoded | Roles `Admin`, `Sponsor`, `Donor`, `Moderator`; admin user `admin` / `younsse1234` |
| `SeedBanks` | `assets/banks.json` | Blood banks with address + 8 blood-group stock rows each |
| `SeedFakeUsers` | hardcoded | 4 demo donors (password `123456`): `karim`, `fatima`, `said`, `anass` |

Default admin credentials:

| Username | Password |
|---|---|
| `admin` | `younsse1234` |

---

## Authentication & Authorization

### Register

`POST /api/account/register` accepts **`multipart/form-data`** (not JSON). Fields read via `c.FormValue`:

| Field | Required | Notes |
|---|---|---|
| `userName` | ✓ | lowercased, must be unique |
| `password` | ✓ | bcrypt-hashed |
| `name` | ✓ | |
| `email` | | |
| `phoneNumber` | | |
| `bloodGroup` | | |
| `city` | | stored on `Address.City` |
| `country` | | stored on `Address.Country` |
| `photo` | | optional file, saved to `./uploads/<uuid>.<ext>` |

Returns `{ userName, name, token }`.

### Login

`POST /api/account/login` accepts JSON `{ userName, password }`. The `userName` field is matched against **either** `user_name` **or** `email` (case-insensitive).

Returns `{ userName, name, gender, photoUrl, token, role }`, where `role` is `Admin`, `Sponsor`, or `Donor` (resolved from the user's roles).

### JWT

- Signed with **HS512**.
- Claims: `nameid` (user ID), `unique_name` (username), `role` (array, if present), `exp` (24h expiry).
- Secret from `TOKEN_KEY` env var (fallback constant).
- Tokens are stored by the frontend under `localStorage.lifeline_token` and sent as `Authorization: Bearer <token>`.

### Protecting routes

- `middleware.Protected()` — validates the JWT; on failure returns `401`.
- `middleware.RequireRole("Admin")` — additionally loads the user and checks they have the required role; returns `403` on failure. The `/admin/*` routes are wrapped in `Protected()` **and** `RequireRole("Admin")`.

---

## API Reference

Base URL: `http://localhost:8080/api`

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/account/register` | — | Register (multipart form) |
| POST | `/account/login` | — | Login, returns JWT |
| GET | `/account/profile` | JWT | Current user's profile |
| GET | `/users/` | JWT | List users (filtered + paginated) |
| GET | `/users/:username` | JWT | Get one user by username |
| GET | `/banks/` | — | List blood banks (filters + pagination) |
| GET | `/geo/cities` | — | List Moroccan cities (sorted A→Z) |
| POST | `/requests/` | JWT | Create a blood request |
| GET | `/requests/` | JWT | List **Approved** requests only |
| GET | `/requests/active` | JWT | My in-progress request (or `null`) |
| PUT | `/requests/:id/accept` | JWT | Donor accepts a request |
| GET | `/notifications` | JWT | My last 10 notifications |
| PUT | `/notifications/:id/read` | JWT | Mark a notification as read |
| GET | `/admin/requests` | Admin | All requests (newest first) |
| PUT | `/admin/requests/:id` | Admin | Set status (`Approved`, etc.) |
| GET | `/admin/users-with-roles` | Admin | Users with their roles |
| POST | `/admin/edit-roles/:username` | Admin | Replace a user's roles |

### Query parameters

**`GET /users/`** (via `dto.UserParams`):

| Param | Default | Notes |
|---|---|---|
| `pageNumber` | `1` | |
| `pageSize` | `10` (max 50) | |
| `gender` | — | |
| `bloodGroup` | — | |
| `minAge` | `18` | |
| `maxAge` | `65` | |
| `orderBy` | `last_active` | `created` also supported |

**`GET /banks/`**:

| Param | Default | Notes |
|---|---|---|
| `pageNumber` | `1` | |
| `pageSize` | `10` (max 100) | |
| `bloodGroup` | — | filters to banks with `quantity > 0` for that group |
| `city` | — | case-insensitive match on `banks.city` |

### Pagination headers

Responses include pagination metadata in a custom header:

- `GET /users/` → `Pagination` header: `{ pageNumber, pageSize, totalPages, totalCount }`
- `GET /banks/` → `X-Pagination` header: `{ currentPage, itemsPerPage, totalItems, totalPages }`

> ⚠️ Inconsistency: the two endpoints use **different header names and field names**. The CORS config in `main.go` only exposes the `Pagination` header explicitly.

### Blood request lifecycle

```
Request created (status: Pending)
        │
        ▼  Admin approves via PUT /admin/requests/:id { status: "Approved" }
        │   └── background goroutine creates a Notification for every donor
        │       with the matching blood_group (skipping the requester)
        ▼
Approved  (visible in GET /requests/)
        │
        ▼  Donor accepts via PUT /requests/:id/accept
        │   └── status set to "In_Progress"
        │   └── Notification sent to the requester
        ▼
In_Progress  (tracked via GET /requests/active)
```

---

## Middleware

| Middleware | Purpose |
|---|---|
| `logger.New()` | Request logging (Fiber built-in) |
| `cors.New(...)` | Allow all origins; expose `Pagination` header |
| `middleware.LogUserActivity()` | After each authenticated request, updates the user's `last_active` (async) |
| `middleware.Protected()` | JWT verification (`gofiber/contrib/jwt`) |
| `middleware.RequireRole(role)` | Role-based access control |
| `middleware.ErrorHandler` | Global error handler (returns `{ statusCode, message }`) |

---

## Services

### `TokenService` (`token_service.go`)

Interface:

```go
type TokenService interface {
    CreateToken(user *models.User, existingToken string) (string, error)
}
```

Creates an HS512-signed JWT with 24h expiry. The `existingToken` parameter is currently unused (refresh-token logic was deferred for the MVP).

### `PhotoService` (`photo_service.go`)

Interface:

```go
type PhotoService interface {
    UploadImage(file *multipart.FileHeader) (url, publicID string, err error)
    DeleteImage(publicID string) error
}
```

Wraps Cloudinary. Uploads to the `lifeline-app` folder with a `w_500,h_500,c_fill,g_face` transformation. **Currently not used by any handler** — registration still saves photos to local disk (`./uploads`) instead.

---

## Known Issues & TODOs

Collected from code comments and the project `todo` file:

- [ ] **Database is SQLite, not PostGIS** — geo proximity matching is not implemented. `go.mod` has no PostGIS driver.
- [ ] **PostGIS radius matching** — `request_handler.go` notes "add geo-location here" when approving a request; currently it notifies *all* donors of the same blood type, not nearby ones.
- [ ] **Photo upload** — `PhotoService` (Cloudinary) is not wired into registration/request creation; files are saved locally.
- [ ] **Status enum drift** — `"In_Progress"` and `"Rejected"` are used as strings but not defined in `RequestStatus`.
- [ ] **Pagination header inconsistency** — `users` uses `Pagination`, `banks` uses `X-Pagination`, with different field names.
- [ ] **Profile update** — `GetUserProfile` has a commented-out update implementation; `PUT /account/profile` is not registered (frontend calls it but it 404s).
- [ ] **Sponsor & moderator endpoints** — `sponsor.go` and `moderate_handler.go` are stubs.
- [ ] **`.env.example`** lists PostgreSQL vars that the code doesn't use.
- [ ] **Hardcoded SQLite path** (`../../lifeline.db`) — brittle if the binary is run from a different cwd.
- [ ] **`PRAGMA foreign_keys = ON`** is enabled, but several handlers rely on manual association handling rather than GORM constraints.
