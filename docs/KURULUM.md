# Kurulum

## Docker Compose

```bash
docker-compose up -d --build
```

## Coolify

1. Repo: `sudoeren/itibar-scraper`
2. Dockerfile: `Dockerfile.rod`
3. Port: 8080
4. Volume: `gmapsdata` → `/gmapsdata`
5. ENV: `DISABLE_TELEMETRY=1`
6. Memory: 2GB, Shm: 2GB
7. Health: `/health`
