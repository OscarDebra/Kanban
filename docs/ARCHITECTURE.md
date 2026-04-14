# Architecture

## Overview

This is a self-hosted, full-stack Kanban board application running on a Raspberry Pi 3B+. The system is composed of three Docker containers managed by Docker Compose, sitting behind an Nginx reverse proxy, with all persistent data stored on a USB drive rather than the Pi's SD card.

---

## Physical hardware

| Component | Role |
|---|---|
| Raspberry Pi 3B+ | Application host — runs all containers |
| 32GB SD card | Operating system (Ubuntu lite 64-bit) and applications (git and Docker f.ex) |
| 64GB SanDisk USB drive | All the rest: code, database data, Docker images |

The SD card is intentionally kept light. SD cards degrade quickly under the sustained random write patterns that a database produces. Keeping the OS on the SD card and all write-heavy data on the USB drive extends the life of both storage devices.

Docker's data root is pointed at `/mnt/usb/docker` via `/etc/docker/daemon.json`. This means pulling a new image (e.g. `postgres:16-alpine`) writes to the USB drive, not the SD card.

---

## Network architecture

```
Browser (any device on Bakka-IM network)
        |
        | HTTPS (port 443) / HTTP redirect (port 80)
        |
Raspberry Pi — port 443 exposed on home network
        |
      Nginx container
      /         \
     /           \
React app      Reverse proxy
               /api/ → Go backend :8080
               /ws   → Go backend :8080 (WebSocket upgrade)
                |
           Go backend container
                |
           Postgres container :5432
                |
           data (/mnt/usb/postgres-data)
```

No ports other than 80 and 443 are exposed outside Docker. The Go backend and Postgres are only reachable from within Docker's internal network (`kanban-network`). This means even if someone on your home network scanned for open ports, they would only see Nginx.

---

## Container architecture

Three containers run via Docker Compose:

### Nginx
- Image: `nginx:alpine`
- Exposed ports: 80 (redirects to HTTPS), 443 (serves app)
- Responsibilities:
  - Serves the built React app as static files
  - Proxies `/api/` requests to the Go backend
  - Upgrades `/ws` connections to WebSocket
  - Enforces HTTPS, security headers, and rate limiting
  - Acts as the only entry point into the system

### Go backend
- Image: built from `./backend/Dockerfile`
- Not exposed outside Docker
- Responsibilities:
  - REST API for all CRUD operations
  - WebSocket hub for real-time collaboration
  - JWT authentication and authorisation
  - File upload handling
  - Scheduled jobs (due date notifications)
  - Input validation and sanitisation

### Postgres
- Image: `postgres:16-alpine`
- Not exposed outside Docker
- Responsibilities:
  - Stores all application data
  - Memory-tuned for Pi 3B+ (32MB shared_buffers, 2MB work_mem, max 20 connections)
  - Data directory mounted from USB drive at `/mnt/usb/postgres-data`

---

## Security architecture

Security is applied in layers — each layer provides independent protection so that if one is bypassed, others remain in place.

### Layer 1 — Transport (Nginx)
- All traffic forced to HTTPS via HTTP → HTTPS redirect on port 80
- TLS 1.2 and 1.3 only — older broken versions blocked
- Weak cipher suites excluded (`!aNULL:!MD5`)
- Self-signed certificate (can be replaced with Let's Encrypt for public exposure)

### Layer 2 — Rate limiting (Nginx)
- Login endpoint: 5 requests per minute per IP (blocks brute force attacks)
- General API: 10 requests per second per IP with burst of 20 (blocks flood attacks)
- Applied at the proxy layer before requests reach the Go backend

### Layer 3 — Browser security (Nginx headers)
- `X-Frame-Options: DENY` — prevents clickjacking via iframes
- `X-Content-Type-Options: nosniff` — prevents MIME-type sniffing attacks
- `Content-Security-Policy` — restricts script and style sources to own domain only
- `Referrer-Policy` — limits referrer information sent to third parties

### Layer 4 — Application (Go backend)
- JWT authentication on all protected routes
- Password hashing with bcrypt
- Input validation and sanitisation on all endpoints
- File type validation on uploads

### Layer 5 — Secrets management
- No secrets stored in code or version control
- All passwords and keys stored in `.env` on the Pi only
- `.env` is gitignored and never committed
- Docker Compose reads variables from `.env` at runtime

### Layer 6 — Network isolation (Docker)
- Only Nginx is exposed to the network
- Go backend and Postgres only reachable within Docker's internal bridge network
- Postgres has no direct external access

---

## Protocols used

| Protocol | Where | Purpose |
|---|---|---|
| HTTPS (HTTP over TLS) | Browser → Nginx | Encrypted web traffic |
| HTTP | Browser → Nginx | Redirected immediately to HTTPS |
| WebSocket Secure (WSS) | Browser → Nginx → Go | Real-time collaboration |
| HTTP/1.1 (internal) | Nginx → Go | API proxying inside Docker |
| TCP | Go → Postgres | Database connections |
| DNS | Pi resolves hostnames | Docker container name resolution |

Docker provides internal DNS — containers can reach each other by name (e.g. the Go backend connects to Postgres at hostname `postgres`) rather than IP address. Docker resolves those names automatically within the `kanban-network` bridge.

---

## Deployment workflow

Development happens on a Mac using VS Code. Changes are pushed to a Git repository, pulled on the Pi, and Docker rebuilds affected containers.

```
Mac (VS Code)
    ↓ git push
GitHub repository
    ↓ git pull (manual or via deploy.sh)
Raspberry Pi
    ↓ docker compose up --build -d
Running containers
```

---

## Current status

- [x] OS installed and updated (Ubuntu 64-bit)
- [x] USB drive mounted permanently at `/mnt/usb`
- [x] Docker and Docker Compose installed
- [x] Docker storage pointed at USB drive
- [x] `docker-compose.yml` configured with all three services
- [x] Nginx configured with HTTPS, rate limiting, security headers
- [x] Secrets moved to `.env` file (not committed to git)
- [x] Git workflow established (Mac → GitHub → Pi)
- [ ] Go backend (in progress)
- [ ] Database schema and migrations
- [ ] React frontend
- [ ] JWT authentication
- [ ] WebSocket collaboration
- [ ] File uploads
- [ ] Notifications
