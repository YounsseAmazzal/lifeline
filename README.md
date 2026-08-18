# Lifeline — Blood Donation Emergency Platform

> **"The Uber of blood donation"** — an intelligent emergency system that connects blood requesters, donors, and corporate sponsors in one place using GPS, push notifications, and real-time matching.

---

## What Is Lifeline?

Lifeline is a Moroccan national blood donation platform that replaces chaotic Facebook groups and outdated directories with a smart, fast, and privacy-first system. It connects three parties:

| Party | Role |
|---|---|
| **Requester** (patient/family) | Submits an urgent blood request with blood type, hospital, and a verified medical document photo |
| **Donor** | Receives a geo-targeted push notification when a match is nearby and accepts or declines |
| **Sponsor** (companies) | Funds the platform via a sponsorship dashboard and gets CSR visibility in return |

---

## Why Lifeline Is Different

| Feature | Others (Facebook/tabarro3.ma) | Lifeline |
|---|---|---|
| Geo-targeting | City-wide noise | Radius-based matching (PostGIS) — only nearby donors are notified |
| Business model | Charity-only, shuts down when funding stops | Sustainable via corporate sponsorships |
| Privacy | Donor's phone number is public | Donor contact is hidden; communication goes through the app |
| Verification | Anyone can post fake requests | Admin/moderator reviews the medical document photo before releasing the alert |
| Performance | PHP / WordPress | Go + Fiber — handles thousands of concurrent requests in milliseconds |

---

## Screenshots

### Home Page (Arabic / RTL)
![Home Arabic](assets/home-ar.png)

### Home Page (English / LTR)
![Home English](assets/home-en.png)

### Sign Up
![Sign Up](assets/sing-up.png)

### Login
![Login](assets/login.png)

### Donate Flow
![Donate](assets/donate.png)

### Donator View
![Donator](assets/donater.png)

### Admin Dashboard
![Admin](assets/admin.png)

### Solution Notification
![Notification](assets/sulotion_nothification.png)

---

## Documentation

Detailed, code-accurate docs live in the `docs/` folder:

- [Backend](docs/backend.md) — architecture, models, seeding, auth, full API reference
- [Frontend](docs/frontend.md) — pages, API layer, i18n, design system, UI features

---

## Tech Stack

### Backend
- **Language:** Go 1.24
- **Framework:** [Fiber v2](https://github.com/gofiber/fiber) — express-style, ultra-fast HTTP
- **Database:** SQLite (via GORM) + PostGIS-ready geo models
- **Auth:** JWT (`golang-jwt/jwt`)
- **Media:** Cloudinary (photo upload via `photo_service`)
- **Architecture:** Repository → Service → Handler (clean layered structure)

### Frontend
- **Stack:** Vanilla HTML/CSS/JS — no framework dependency
- **Styling:** Tailwind CSS (CDN)
- **Maps:** Leaflet.js
- **i18n:** Custom auto-translation system (`transletesystem.js`) — Arabic RTL & English LTR
- **Icons:** Font Awesome

---

## Project Structure

```
lifeline/
├── backend/
│   ├── cmd/api/main.go          # Entry point
│   ├── internal/
│   │   ├── handlers/            # HTTP handlers (auth, user, admin, request, bank, geo, notification)
│   │   ├── middleware/          # JWT guard, role guard, activity logger, error handler
│   │   ├── models/              # GORM models (User, Request, Bank, Sponsor, Notification…)
│   │   ├── dto/                 # Request/response shapes
│   │   ├── services/            # Business logic (token, photo upload)
│   │   └── routes/              # Route registration
│   ├── pkg/database/            # DB connection, seeders (cities, roles, banks, fake users)
│   └── uploads/                 # Uploaded medical document photos
├── frontend/
│   ├── index.html               # Landing page
│   ├── auth/                    # login.html, register.html, dashboard.html
│   ├── views/                   # sponsor.html, profile.html, admin.html, bank.html
│   ├── api/                     # JS API layer (auth, admin, dashboard, profile, pagination)
│   └── assets/css|js/           # Tailwind config, global styles, i18n, city/bank data
└── assets/                      # Screenshots used in this README
```

---

## API Overview

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/api/account/register` | — | Register a new user |
| POST | `/api/account/login` | — | Login, returns JWT |
| GET | `/api/account/profile` | JWT | Get current user profile |
| GET | `/api/users/` | JWT | List users |
| GET | `/api/users/:username` | JWT | Get user by username |
| GET | `/api/banks/` | — | List blood banks |
| GET | `/api/geo/cities` | — | List Moroccan cities |
| POST | `/api/requests/` | JWT | Create a blood request |
| GET | `/api/requests/` | JWT | List requests |
| GET | `/api/requests/active` | JWT | Get my active request |
| PUT | `/api/requests/:id/accept` | JWT | Accept a request (donor) |
| GET | `/api/notifications` | JWT | Get my notifications |
| PUT | `/api/notifications/:id/read` | JWT | Mark notification as read |
| GET | `/api/admin/requests` | Admin | View all pending requests |
| PUT | `/api/admin/requests/:id` | Admin | Approve or reject a request |
| GET | `/api/admin/users-with-roles` | Admin | List users with their roles |
| POST | `/api/admin/edit-roles/:username` | Admin | Edit a user's roles |

---

## How to Run

### Prerequisites
- [Go 1.24+](https://go.dev/dl/)
- [Node.js](https://nodejs.org/) (only for `npx serve`)

### 1. Clone the repository

```bash
git clone https://github.com/YounsseAmazzal/lifeline.git
cd lifeline
```

### 2. Configure environment

```bash
cp backend/.env.example backend/.env
# Fill in: JWT_SECRET, CLOUDINARY_URL, PORT (default 8080)
```

### 3. Start the backend

```bash
cd backend
go run cmd/api/main.go
```

The server starts at `http://localhost:8080`. On first run it auto-migrates the database and seeds:
- Moroccan regions & cities
- Default roles (`Admin`, `Donor`, `Moderator`)
- Blood banks
- Demo users

### 4. Start the frontend

Open a new terminal at the project root:

```bash
npx serve frontend
```

Then open the URL shown in your terminal (usually `http://localhost:3000`).

---

## Run with Docker

A `docker-compose.yml` at the repo root builds and runs both services.

```bash
docker compose up --build
```

- **Backend API** → http://localhost:8080
- **Frontend** → http://localhost:3000

Configuration (optional): create a `.env` file at the repo root to override defaults:

```bash
TOKEN_KEY=your_long_random_secret
CLOUDINARY_URL=cloudinary://API_KEY:API_SECRET@CLOUD_NAME
```

Persistent data lives in two named volumes (`backend_data` for SQLite, `backend_uploads` for photos).

> Note: the frontend calls the API at the hardcoded `http://localhost:8080/api` (see `frontend/api/api.js`), which matches the published backend port above.

---

## Contributing

Contributions are welcome. Please follow these steps:

1. **Fork** the repository and create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes** — keep commits small and focused.

3. **Test locally** — run both backend and frontend and verify the affected flow.

4. **Open a Pull Request** against `main` with a clear description of what you changed and why.

### Conventions

- Backend: follow standard Go naming conventions; keep handlers thin, business logic in services.
- Frontend: keep API calls inside the `frontend/api/` layer, not inline in HTML files.
- No external CSS frameworks beyond Tailwind; no frontend build step required.
- For new routes, register them in `backend/internal/routes/routes.go` and document them in the API table above.

### Good First Issues

- Add pagination to the donor list
- Implement the sponsor dashboard UI (`views/sponsor.html`)
- Connect Leaflet map to live request geo-data
- Add email/SMS notification fallback

---

## License

MIT — free to use, modify, and distribute.

---

> Built for Morocco. Powered by Go & community. Every drop counts.
