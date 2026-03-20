# İtibar Scraper

API-only Google Maps scraper. Sadece REST API sunar.

## Başlat

```bash
# İlk çalıştırma (build gerekir)
docker-compose up -d --build

# Sonraki çalıştırmalar
docker-compose up -d
```

Bu kadar. Scraper `http://localhost:8080` adresinde çalışır.

## Durdur

```bash
docker-compose down
```

## Verileri Sil

```bash
docker-compose down -v
```

## Logları İncele

```bash
docker-compose logs -f
```

## API Kullanımı

### Health Check

```bash
curl http://localhost:8080/health
```

### İş Oluştur

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["cafe in Istanbul"],
    "lang": "tr",
    "depth": 1,
    "max_time": 600
  }'
```

### İş Durumu

```bash
curl http://localhost:8080/api/v1/jobs/{id}
```

### Sonuç İndir

```bash
curl -O http://localhost:8080/api/v1/jobs/{id}/download
```

## Geliştirme

```bash
# Build et ve çalıştır
docker-compose up -d --build

# Sadece build
docker build -f Dockerfile.rod -t itibar-scraper .

# Sadece run
docker-compose up -d
```

## Coolify Deploy

1. GitHub repo'yu bağla: `sudoeren/itibar-scraper`
2. Dockerfile: `Dockerfile.rod`
3. CMD: `-data-folder /gmapsdata`
