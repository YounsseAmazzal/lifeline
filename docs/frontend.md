# Lifeline — Frontend Documentation

The frontend is a **vanilla HTML/CSS/JS** single-page-style app (no build step, no framework) that talks to the Go backend over `fetch`. It includes a tri-lingual interface (Arabic RTL, English LTR, French LTR), an interactive Leaflet map, live notifications, and a full auth flow.

---

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Getting Started](#getting-started)
3. [Project Structure](#project-structure)
4. [Pages](#pages)
5. [API Layer](#api-layer)
6. [Authentication Flow](#authentication-flow)
7. [Internationalization (i18n)](#internationalization-i18n)
8. [Design System](#design-system)
9. [Key UI Features](#key-ui-features)
10. [Known Issues & TODOs](#known-issues--todos)

---

## Tech Stack

| Concern | Choice |
|---|---|
| Markup | HTML5 |
| Styling | [Tailwind CSS](https://tailwindcss.com) (CDN, no build step) + one custom `main.css` |
| Scripting | Vanilla JavaScript (ES modules + global `window` functions) |
| Maps | [Leaflet 1.9.4](https://leafletjs.com) + Carto Voyager tiles |
| i18n | Custom translation system (`transletesystem.js` + `lang.js`) |
| Icons | [Font Awesome 6](https://fontawesome.com) (CDN) |
| Animations | [AOS](https://michalsnik.github.io/aos/) scroll animations |
| Charts | [Chart.js](https://www.chartjs.org) (sponsor page only) |
| Fonts | Cairo (Arabic) + Plus Jakarta Sans (Latin) |

> There is no `package.json` dependency graph for the frontend — everything is loaded from CDNs. The root `package.json` is empty.

---

## Getting Started

### Prerequisites
- [Node.js](https://nodejs.org) (only for the static server; `npx` ships with npm)

### Run

Serve the `frontend/` directory as static files (from the project root):

```bash
npx serve frontend
```

Then open the printed URL (usually `http://localhost:3000`).

> The backend must be running separately at `http://localhost:8080` (see [`backend.md`](backend.md)). The API base URL is hardcoded in `frontend/api/api.js`.

---

## Project Structure

```
frontend/
├── index.html                  # Landing page (public marketing page)
├── auth/
│   ├── login.html              # Login page
│   ├── register.html           # Registration page
│   └── dashboard.html          # Donor dashboard (map + banks + requests)
├── views/
│   ├── admin.html              # Admin panel (approve/reject requests, roles)
│   ├── profile.html            # User profile (edit + availability toggle)
│   ├── bank.html               # Blood-bank management (partial)
│   └── sponsor.html            # Sponsor dashboard (partial, Chart.js)
├── api/                        # API access layer (thin fetch wrappers)
│   ├── api.js                  #   apiRequest / apiRequestForm + auth/account/geo objects
│   ├── auth/login.js           #   login form handler
│   ├── auth/register.js        #   register form handler (+ city loading)
│   ├── dashbord/main.js        #   dashboard: map, banks, requests, notifications
│   ├── admin/adminrequest.js   #   admin: load + approve/reject requests
│   ├── profile/main.js         #   profile: load/update profile
│   ├── pasination.js           #   bank pagination logic
│   ├── globalalert.js          #   showAutoAlert toast helper
│   └── confirm.js              #   showConfirm modal (ES module)
└── assets/
    ├── css/
    │   ├── main.css            #   custom design tokens + confirm-modal styles
    │   ├── login.css           #   (2 lines, effectively unused)
    │   ├── register.css        #   (2 lines, effectively unused)
    │   └── tailwindconfig.js   #   Tailwind theme (brand colors, fonts, animations)
    ├── js/
    │   ├── lang.js             #   language preference + RTL/LTR layout helper
    │   ├── transletesystem.js  #   i18n dictionary + changeLanguage()
    │   ├── CitiesAndBanks.js   #   landing-page city/bank loading
    │   └── main.js             #   (empty placeholder)
    └── img/sponsors/           #   sponsor logos (CDG, Maroc Telecom, OCP, RAM, CIH, inwi, Attijariwafa)
```

---

## Pages

Each HTML page is a self-contained entry point that pulls in the shared `api.js`, `lang.js`, `globalalert.js`, and a page-specific script.

| Page | Purpose | Page-specific script |
|---|---|---|
| `index.html` | Public landing page: hero, sponsors marquee, impact stats, "how it works", bank stock grid | `CitiesAndBanks.js`, `transletesystem.js` |
| `auth/login.html` | Email + password login | `auth/login.js` |
| `auth/register.html` | Registration with city dropdown (loaded from backend) + photo upload | `auth/register.js` |
| `auth/dashboard.html` | Donor dashboard: Leaflet map of banks, bank list, blood requests feed, notifications panel, request modal | `dashbord/main.js` |
| `views/profile.html` | View/edit profile, availability toggle, logout | `profile/main.js` |
| `views/admin.html` | List all blood requests; approve/reject pending ones | `admin/adminrequest.js` |
| `views/bank.html` | Blood-bank management UI | *(no page script wired yet)* |
| `views/sponsor.html` | Sponsor analytics dashboard (Chart.js) | *(inline script only)* |

### Routing / redirects

There is no client-side router. Navigation is done with `window.location.href`. After login, users are redirected based on their role:

| Role | Redirect target |
|---|---|
| `Admin` | `views/admin.html` |
| `Sponsor` | `sponsor.html` (note: relative path bug, see known issues) |
| `Donor` (default) | `dashboard.html` (relative to `auth/`) |

---

## API Layer

All network access is centralized in `api/api.js`.

### `apiRequest(endpoint, method, data)`

- Base URL: `const API_URL = "http://localhost:8080/api"`.
- Sends `Content-Type: application/json`.
- Attaches `Authorization: Bearer <token>` if `localStorage.lifeline_token` exists.
- On `401`, clears the stored token and throws `"Unauthorized"`.
- Parses JSON when the response content-type is JSON, and throws the server's `error` field on non-OK responses.

### `apiRequestForm(endpoint, method, formData)`

Same as above but for `multipart/form-data` (no content-type header; body is a `FormData` object).

### Pre-built API objects

```js
auth.login(userName, password)             // POST /account/login
auth.register(userData)                    // POST /account/register (FormData or JSON)

accountApi.profile()                       // GET  /account/profile
accountApi.updateProfile(formData)         // PUT  /account/profile/photo (currently 404s)
accountApi.uploadPhoto(formData)           // POST /account/profile/photo (currently 404s)

geoApi.cities()                            // GET  /geo/cities
```

Other endpoints are called directly with `apiRequest(...)` inside the page scripts:

| Call site | Endpoint |
|---|---|
| `dashbord/main.js` | `/banks?pageSize=50`, `/requests`, `/requests/active`, `/requests/:id/accept`, `/notifications`, `/notifications/:id/read` |
| `admin/adminrequest.js` | `/admin/requests`, `/admin/requests/:id` |
| `pasination.js` | `/banks?pageNumber=..&pageSize=..` (raw `fetch`) |
| `auth/register.js` | `/geo/cities` (via `geoApi`) |

> ⚠️ `CitiesAndBanks.js` calls `banksApi.getAll(...)`, but **`banksApi` is not defined** anywhere in `api.js`. The landing-page bank grid therefore falls back to its error branch.

---

## Authentication Flow

1. **Login** (`auth/login.js`) submits `{ userName, password }`. On success it stores the JWT in `localStorage.lifeline_token` and redirects by role.
2. **Registration** (`auth/register.js`) builds a `FormData` with `name`, `userName` (= email), `email`, `phoneNumber`, `password`, `bloodGroup`, `city`, `country`, optional `photo`. On success it stores the token and redirects to `login.html`.
3. **Guarded pages** (`dashboard`, `profile`) check for a token on `DOMContentLoaded`; if absent they redirect to `login.html`. `apiRequest` also auto-clears the token on `401` and the page-level `catch` redirects to login.
4. **Logout** (`profile/main.js`) simply removes `localStorage.lifeline_token` and redirects.

The token is stored in `localStorage` and read by `api.js` on every request.

---

## Internationalization (i18n)

Two files cooperate:

### `lang.js` — language state & layout

- Key: `localStorage.lifeline_lang` (`ar` | `en` | `fr`, default `ar`).
- `getPreferredLanguage()` — reads stored value, falls back to browser language.
- `persistLanguage(lang)` — saves the choice.
- `applyLanguageLayout(lang)` — sets `<html lang>` and `dir` (`rtl` for Arabic, `ltr` otherwise) and swaps the body font class (`font-arabic` vs `font-sans`).

### `transletesystem.js` — dictionary & switching

- Defines an `i18n` object with `ar`, `en`, and `fr` dictionaries keyed by a translation ID (e.g. `hero_title_1`, `btn_donate`).
- `changeLanguage(lang)` updates every element with a `data-i18n="<key>"` attribute and flips direction/font.
- Runs on `DOMContentLoaded` using the preferred language.

Usage in HTML:

```html
<span data-i18n="btn_donate">تبرع الآن</span>
```

> Note: only the landing page is fully wired for translation. The auth/views pages are largely hardcoded in Arabic.

---

## Design System

Brand colors and tokens are defined in `assets/css/tailwindconfig.js` and `assets/css/main.css`.

### Colors

| Token | Hex | Usage |
|---|---|---|
| `morocco-red` | `#C1272D` | Primary CTA, accents, the Moroccan flag red |
| `morocco-green` | `#006233` | Secondary accent, success, the flag green |
| `morocco-gold` | `#D4AF37` | Decorative accents, crown icon |
| `morocco-dark` | `#0f172a` | Dark sections |

### Fonts

- **Arabic:** Cairo (`.font-arabic`, applied when `dir="rtl"`)
- **Latin:** Plus Jakarta Sans (`.font-sans`)

### Custom CSS classes (`main.css`)

- `.glass-royal` — frosted-glass sticky navbar.
- `.pattern-overlay` — subtle zellige/tile SVG background.
- `.text-gradient-flag` — red→green gradient text.
- `.confirm-overlay` / `.confirm-box` / `.btn` — the reusable confirmation modal styling.

### Tailwind animation keyframes

`scroll` (sponsor logo marquee), `float`, `pulse-slow`.

---

## Key UI Features

### Leaflet map (`dashbord/main.js`)

- Initializes a map centered on Morocco (`[31.7917, -7.0926]`).
- Fetches `accountApi.profile()` and `/banks?pageSize=50` in parallel.
- Uses the browser's `navigator.geolocation` to center on the user and drop a pulsing marker.
- Renders a custom hospital icon per bank, with a popup showing stock status and a Google Maps directions link.

### Blood request feed

- `loadRequests()` fetches `/requests` (approved only) and renders cards with an `Accept & Go` button.
- `acceptRequest(id, lat, lng)` calls `PUT /requests/:id/accept` and opens Google Maps directions to the hospital.
- A "tracking card" is shown if the user has an in-progress request (`GET /requests/active`).

### Request modal

- Loads cities from `geoApi.cities()`, then loads hospitals for the selected city via `/banks?city=..`.
- On submit, POSTs `{ bloodType, hospitalName, isUrgent, latitude, longitude }` to `/requests`.

### Notifications

- `checkNotifications()` polls `GET /notifications` **every 30 seconds** (`setInterval`).
- Unread count is shown as a badge; clicking a notification calls `PUT /notifications/:id/read`.

### Pagination (`pasination.js`)

- Paginates banks at 6 per page using the `X-Pagination` response header.
- Renders prev/next + numbered buttons with ellipsis gaps.

### Toast alerts (`globalalert.js`)

- `showAutoAlert(message, type)` shows a fixed-position toast (green for `success`, red for `error`) that fades out after 3 seconds.
- Exposed globally as `window.showAutoAlert`.

### Confirmation modal (`confirm.js` / inline in dashboard & admin)

- `showConfirm(message)` returns a `Promise<boolean>` and drives a shared `#globalConfirmModal`.
- Used before destructive/outward actions (approving/rejecting requests, accepting a donation).

---

## Known Issues & TODOs

- [ ] **`banksApi` is undefined** in `CitiesAndBanks.js` — the landing-page bank grid always falls back to the error state.
- [ ] **`accountApi.updateProfile` / `uploadPhoto`** point to endpoints (`PUT /account/profile`, `POST /account/profile/photo`) that are **not registered** in the backend `routes.go`.
- [ ] **Sponsor redirect bug** — `login.js` redirects `Sponsor` to `sponsor.html` (relative to `auth/`), but the page lives in `views/sponsor.html`.
- [ ] **`login.html`** includes a stray `<script src="js/api.js">` (wrong path) plus a large inline script before the real scripts load.
- [ ] **Profile city loading** (`profile/main.js`) reads `city.name` from `/geo/cities`, but the API returns `ville` — the select is populated with `undefined`.
- [ ] **Profile form** references fields (`name`, `email`, etc.) via `form.elements`; if the HTML field names drift, these silently no-op.
- [ ] **Bank management** (`views/bank.html`) and **sponsor dashboard** (`views/sponsor.html`) have no functional page script wired up yet.
- [ ] **Duplicated `showConfirm` / `loadRequests`** — defined both inline and in `confirm.js`; several pages re-define the same helpers.
- [ ] **Hardcoded `API_URL`** — `http://localhost:8080/api` is baked in; no environment-based config for the frontend.
- [ ] **Translation coverage** — only the landing page uses `data-i18n`; auth/views pages are hardcoded Arabic.
- [ ] **No frontend tests, linter, or bundler** — everything is loaded ad-hoc from CDNs.
