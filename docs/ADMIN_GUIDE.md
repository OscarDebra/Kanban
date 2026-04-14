# Admin guide

This guide is for the person responsible for running and maintaining the Kanban board. It assumes basic familiarity with the command line and SSH.

---

## Connecting to the Pi

From any computer on your home network:

```bash
ssh ubuntu@192.168.x.x
```

Replace `192.168.x.x` with your Pi's actual IP address. Check your router's connected devices list if you are unsure what it is.

---

## Starting and stopping the application

### Start everything

```bash
cd ~/kanban
docker compose up -d
```

### Stop everything

```bash
cd ~/kanban
docker compose down
```

### Restart everything (e.g. after a config change)

```bash
cd ~/kanban
docker compose down
docker compose up -d
```

### Restart a single container

```bash
docker restart kanban-backend
docker restart kanban-nginx
docker restart kanban-postgres
```

### Check whether containers are running

```bash
docker ps
```

All three containers (`kanban-nginx`, `kanban-backend`, `kanban-postgres`) should show a status of `Up`.

---

## Deploying updates

When new code has been pushed to the Git repository:

```bash
cd ~/kanban
./deploy.sh
```

This pulls the latest code from Git and rebuilds and restarts any containers whose code has changed. If you prefer to run the steps manually:

```bash
cd ~/kanban
git pull
docker compose up --build -d
```

The `--build` flag only rebuilds containers where something has changed. Postgres and Nginx (which use pre-built images) are not affected unless their configuration files changed.

---

## Viewing logs

### All containers at once

```bash
docker compose logs -f
```

Press `Ctrl+C` to stop following.

### A specific container

```bash
docker logs kanban-backend -f
docker logs kanban-nginx -f
docker logs kanban-postgres -f
```

### Last 100 lines only

```bash
docker logs kanban-backend --tail 100
```

Logs are the first place to look when something is not working correctly.

---

## Managing users

User accounts are managed through the application itself. As an administrator you have access to an admin panel at:

```
https://192.168.x.x/admin
```

From here you can:

- View all registered users
- Deactivate or delete accounts
- Reset a user's password
- View which boards a user has access to

> **Note:** The admin panel is only accessible to accounts with the administrator role. The first account registered automatically receives administrator privileges.

---

## Backups

### What needs to be backed up

Only the USB drive needs to be backed up. Specifically:

| Path | Contents |
|---|---|
| `/mnt/usb/postgres-data` | All database data — boards, tasks, users, everything |
| `/mnt/usb/uploads` | All user file attachments |

The application code lives in Git and does not need separate backup. The SD card (OS) does not need to be backed up — if it fails, you reinstall Ubuntu and run the setup steps again.

### Manual backup

To create a database dump:

```bash
docker exec kanban-postgres pg_dump -U kanban kanban > ~/backup-$(date +%Y%m%d).sql
```

This creates a file like `backup-20250415.sql` in your home directory. Copy it somewhere safe (external drive, cloud storage):

```bash
scp ubuntu@192.168.x.x:~/backup-20250415.sql /path/on/your/mac/
```

To back up file attachments:

```bash
scp -r ubuntu@192.168.x.x:/mnt/usb/uploads /path/on/your/mac/uploads-backup/
```

### Restoring from a backup

Stop the backend first to prevent writes during restore:

```bash
docker stop kanban-backend
```

Restore the database:

```bash
docker exec -i kanban-postgres psql -U kanban kanban < ~/backup-20250415.sql
```

Restart the backend:

```bash
docker start kanban-backend
```

To restore file attachments, copy them back to `/mnt/usb/uploads`.

### Automated backups (recommended)

Set up a weekly automated backup using cron. On the Pi:

```bash
crontab -e
```

Add this line to run a backup every Sunday at 2am:

```
0 2 * * 0 docker exec kanban-postgres pg_dump -U kanban kanban > /mnt/usb/backups/backup-$(date +\%Y\%m\%d).sql
```

Create the backups folder first:

```bash
sudo mkdir -p /mnt/usb/backups
```

---

## Renewing the HTTPS certificate

The self-signed certificate lasts one year. When it expires, browsers will show a more severe warning. To renew it:

```bash
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/certs/kanban.key \
  -out /etc/nginx/certs/kanban.crt \
  -subj "/CN=kanban.local"

docker restart kanban-nginx
```

---

## Rotating secrets

If you suspect a password or secret has been compromised, rotate it immediately.

### Rotating the database password

1. Edit the `.env` file on the Pi:
```bash
nano ~/kanban/.env
```
Update `DB_PASSWORD` to a new value generated with:
```bash
openssl rand -base64 32
```

2. Update the password inside Postgres:
```bash
docker exec -it kanban-postgres psql -U kanban -c "ALTER USER kanban PASSWORD 'your-new-password';"
```

3. Restart the containers:
```bash
cd ~/kanban
docker compose down
docker compose up -d
```

### Rotating the JWT secret

Edit `.env` and update `JWT_SECRET` to a new value. Then restart the backend:

```bash
docker restart kanban-backend
```

> **Note:** Rotating the JWT secret immediately invalidates all existing login sessions. All users will be logged out and will need to log in again.

---

## Checking disk usage

### USB drive usage

```bash
df -h /mnt/usb
```

### Breakdown by folder

```bash
du -sh /mnt/usb/*
```

### Docker image and container sizes

```bash
docker system df
```

If Docker images are taking up too much space, remove unused ones:

```bash
docker image prune -a
```

---

## Updating Docker images

Postgres and Nginx use pre-built images. To update them to the latest patch version:

```bash
cd ~/kanban
docker compose pull
docker compose up -d
```

This pulls newer versions of `nginx:alpine` and `postgres:16-alpine` if available. The `16-alpine` tag tracks the latest Postgres 16 patch release — major version upgrades (e.g. 16 → 17) require a manual migration and should not be done automatically.

---

## Common problems

### A container keeps restarting

Check its logs for the error:

```bash
docker logs kanban-backend --tail 50
```

Common causes:
- The `.env` file is missing or has incorrect values
- The USB drive is not mounted (run `df -h` and check `/mnt/usb` appears)
- Postgres has not finished starting before the backend tries to connect (usually resolves itself within 30 seconds)

### The app is not accessible in the browser

1. Check all containers are running: `docker ps`
2. Check the Pi's IP address has not changed: log into your router
3. Check Nginx logs: `docker logs kanban-nginx --tail 50`
4. Check port 443 is open: `sudo ss -tlnp | grep 443`

### The USB drive is not mounted after a reboot

```bash
sudo mount -a
```

If this fails, check the UUID in `/etc/fstab` still matches the drive:

```bash
sudo blkid /dev/sda1
cat /etc/fstab
```

### Postgres database is corrupt

This is rare but can happen after an unclean shutdown (power loss). Try:

```bash
docker stop kanban-backend kanban-postgres
docker start kanban-postgres
docker logs kanban-postgres --tail 50
```

Postgres will usually attempt automatic recovery. If it cannot recover, restore from your most recent backup.

---

## Security checklist

Run through this periodically:

- [ ] `.env` file exists on Pi and contains strong passwords
- [ ] `.env` is not committed to the Git repository
- [ ] HTTPS certificate is not expired (`openssl x509 -enddate -noout -in /etc/nginx/certs/kanban.crt`)
- [ ] Backups are being created and stored somewhere other than the Pi
- [ ] Docker images are up to date (`docker compose pull`)
- [ ] No unnecessary ports are open (`sudo ss -tlnp`)
- [ ] OS packages are up to date (`sudo apt update && sudo apt upgrade`)
