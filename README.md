
# RentTracker — local development README

This repository contains a React Native (Expo) frontend and a Go backend using SQLite for storage. The instructions below explain how to set up the initial database, run the Go server, and run the Expo app locally on Windows (PowerShell).

If you already have a working development environment you can follow just the sections you need (Database → Backend → Frontend).

---

## Requirements

- Node.js (16+ recommended) and npm (or yarn)
- Expo CLI (optional) — you can use `npx expo` without a global install
- Go toolchain (1.18+ recommended)
- sqlite3 CLI (for creating/loading the local DB)

On Windows (PowerShell) you can install:
- Node.js: https://nodejs.org/
- Go: https://go.dev/dl/
- sqlite3: install via Chocolatey (choco) or download prebuilt binaries

Example (Chocolatey):

```powershell
choco install nodejs-lts golang sqlite
```

---

## 1) Setup the initial SQLite database

This project ships with `rt.sql` (top-level). It contains the schema and view definitions used by the backend.

From the project root (PowerShell):

```powershell
# Optional: back up existing DB
if (Test-Path .\rt.db) { copy-item .\rt.db .\rt.db.bak }

# Create/overwrite the DB by executing the schema file
sqlite3 .\rt.db ".read rt.sql"

# Quick verification: list overdue and upcoming views
sqlite3 .\rt.db "SELECT leaseId, firstName, lastName, rentAmount, lastPaymentUnix, paymentStatus FROM overduePayments;"
sqlite3 .\rt.db "SELECT leaseId, firstName, lastName, rentAmount, lastPaymentUnix, paymentStatus FROM upcomingPayments;"
```

Notes
- Running `.read rt.sql` will DROP and CREATE tables and views as defined in the file. Back up `rt.db` first if it contains important data.
- If you modify `rt.sql`, re-run the `.read` command to apply changes to `rt.db`.

---

## 2) Run the Go backend server

The backend lives in the `backend/` directory and exposes API endpoints the app expects (payments, leases, maintenance, dashboard views, etc.).

From PowerShell:

```powershell
cd .\backend

# download modules
go mod download

# build
go build

# run (Windows)
.\main.exe

# or run without building
go run .
```

Important
- The frontend's API client is configured in `apis/client.ts` with a BASE_URL (currently a LAN IP). If you run the backend locally on the same machine, ensure `BASE_URL` points to the backend address and port (for example `http://10.0.0.67:8080` or `http://localhost:8080`). Update `apis/client.ts` if needed.

---

## 3) Run the Expo frontend (React Native)

From the project root:

```powershell
npm install

# Start Expo dev server (tunnel recommended for device testing)
npx expo start --tunnel

# Or if you prefer the default
npx expo start
```

Then open the project in Expo Go or a simulator. If using an emulator/simulator, follow Expo's instructions for Android Studio or Xcode.

Notes
- The frontend periodically fetches data from the backend endpoints. Make sure the backend is running and reachable from your device (emulator or physical device). If using a phone, prefer `--tunnel` or ensure your machine's LAN IP is accessible.

---

## Useful development commands

- Recreate DB: `sqlite3 .\rt.db ".read rt.sql"`
- Build backend: `cd backend; go build`
- Run backend without building: `cd backend; go run .`
- Start Expo: `npx expo start`
- Install JS deps: `npm install`

---

## Troubleshooting

- Backend cannot connect to DB: confirm `rt.db` exists in the repo root and `main.go` opens it. The backend expects a local file named `rt.db`.
- API calls failing from app: check `apis/client.ts` BASE_URL and adjust to your backend address.
- Views not reflecting changes: if you edit `rt.sql`, re-run the `.read rt.sql` command and restart the backend (views are resolved by the DB engine at query time but the backend may cache things).

---

## Dependencies (high level)

- Node / npm: use package.json to install JS dependencies (`npm install`).
- Expo: no global install required; use `npx expo`.
- Go: standard modules are declared in `backend/go.mod`; run `go mod download` to fetch them.
- sqlite3 CLI: used for manually creating the DB from `rt.sql`.

If you want I can add a small script (PowerShell or npm script) to reset/create the DB and start backend + frontend together.

---

If anything above doesn't work in your environment, tell me what OS and exact error messages you see and I'll help you resolve them.

Happy hacking — use the `issues` workflow or ask here if you want automation scripts added.
