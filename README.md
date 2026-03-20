# İtibar Scraper

API-only Google Maps scraper. Sadece REST API sunar.

## Başlat

```bash
# Build et ve çalıştır
docker-compose up -d --build
```

Scraper `http://localhost:8080` adresinde çalışır.

## Durdur

```bash
docker-compose down
```

## Verileri Sil

```bash
docker-compose down -v
```

## Log

```bash
docker-compose logs -f
```

## API

```bash
# Health check
curl http://localhost:8080/health

# İş oluştur
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{"keywords": ["cafe in Istanbul"], "lang": "tr", "depth": 1}'

# İş durumu
curl http://localhost:8080/api/v1/jobs/{id}

# Sonuç indir
curl -O http://localhost:8080/api/v1/jobs/{id}/download
```

## Coolify

1. Repo: `sudoeren/itibar-scraper`
2. Dockerfile: `Dockerfile.rod`
3. CMD: `-data-folder /gmapsdata`
