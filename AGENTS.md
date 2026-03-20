# İtibar Scraper - Geliştirme

## Proje

[gosom/google-maps-scraper](https://github.com/gosom/google-maps-scraper) üzerine inşa edilmiş, API-only Google Maps scraper.

## Build

```bash
make test
make vet
make format
```

## Proje Yapısı

```
├── api/           # REST API
├── cmd/           # CLI
├── gmaps/         # Veri modelleri
├── runner/        # Çalışma modları
├── web/           # API sunucu
└── main.go
```

## Önemli Dosyalar

| Dosya | Açıklama |
|-------|----------|
| `web/web.go` | API endpoint'leri |
| `web/job.go` | JobData yapısı |
| `runner/webrunner/` | Web runner |
| `gmaps/entry.go` | Entry ve Review modelleri |

## API

| Endpoint | Method | Açıklama |
|----------|--------|----------|
| `/health` | GET | Health check |
| `/api/v1/jobs` | POST | İş oluştur |
| `/api/v1/jobs` | GET | İşleri listele |
| `/api/v1/jobs/{id}` | GET | İş durumu |
| `/api/v1/jobs/{id}` | DELETE | İşi sil |
| `/api/v1/jobs/{id}/download` | GET | CSV indir |
