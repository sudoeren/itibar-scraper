# Kurulum

## Gereksinimler

- Docker
- Docker Compose

## Docker Compose (Önerilen)

### 1. Yapılandırma

`docker-compose.yaml` zaten hazır.

### 2. Başlat

```bash
docker-compose up -d --build
```

İlk çalıştırmada build eder. Sonraki çalıştırmalarda sadece:
```bash
docker-compose up -d
```

Scraper `http://localhost:8080` adresinde çalışır.

### 3. Kontrol Et

```bash
curl http://localhost:8080/health
```

## Coolify'da Deploy

1. **Repository:** GitHub'dan repo'yu bağla (`sudoeren/itibar-scraper`)
2. **Build:** Dockerfile kullan (`Dockerfile.rod`)
3. **Port:** `8080`
4. **Volume:** `gmapsdata` → `/gmapsdata`
5. **ENV:** `DISABLE_TELEMETRY=1`
6. **Memory:** 2GB
7. **Shm:** 2GB
8. **Health Check:** `/health`
9. **Deploy:** Auto-deploy veya manual

## Kaynak Koddan Çalıştırma

```bash
# Build
docker build -f Dockerfile.rod -t itibar-scraper .

# Run
docker-compose up -d
```

## Sorun Giderme

### "Page Closed" Hataları

Shm boyutunu artır:
```yaml
shm_size: 4gb
```

### Hafıza Hataları

Memory limit artır:
```yaml
memory: 4G
```

### Rate Limit

Proxy ekle:
```yaml
command: -data-folder /gmapsdata -proxies 'socks5://user:pass@host:1080'
```
