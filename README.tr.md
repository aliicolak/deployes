<div align="center">

# 🚀 deployes

### GitHub Deployment Otomasyon Platformu

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Angular](https://img.shields.io/badge/Angular-18+-DD0031?style=for-the-badge&logo=angular&logoColor=white)](https://angular.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/Lisans-MIT-green?style=for-the-badge)](LICENSE)

**deployes, GitHub projelerinizi tek tıklamayla veya webhook'lar aracılığıyla otomatik olarak uzak sunuculara deploy etmenizi sağlayan, kendi sunucunuzda barındırabileceğiniz modern bir deployment otomasyon platformudur.**

[🇬🇧 English Documentation](README.md)

---

<img src="https://raw.githubusercontent.com/aliicolak/deployes/main/docs/screenshot-dashboard.png" alt="deployes Dashboard" width="800"/>

</div>

## ✨ Özellikler

### 🎯 Temel Özellikler
- **Tek Tıkla Deploy** - Herhangi bir projeyi herhangi bir sunucuya tek tıkla deploy edin
- **GitHub Webhook Entegrasyonu** - Push event'lerinde otomatik deployment
- **Branch Tabanlı Tetikleyiciler** - Belirli branch'leri deployment tetiklemesi için yapılandırın
- **Çoklu Sunucu Yönetimi** - Tek bir dashboard'dan sınırsız sunucu yönetin
- **Şifrelenmiş Kimlik Bilgileri** - SSH anahtarları ve secret'lar AES-256 ile şifrelenir

### 🔐 Güvenlik
- **JWT Kimlik Doğrulama** - Güvenli token tabanlı authentication
- **SSH Anahtar Yönetimi** - Private repository'ler için otomatik oluşturulan deploy key'ler
- **Environment Secret'ları** - Deployment'lar için şifrelenmiş ortam değişkenleri
- **Parola Saklanmaz** - Sadece SSH anahtar tabanlı sunucu authentication'ı

### 📊 İzleme ve Loglar
- **Gerçek Zamanlı Deployment Logları** - WebSocket tabanlı canlı streaming
- **Deployment Geçmişi** - Durum takibi ile tam geçmiş
- **Dashboard Analitiği** - Deployment istatistiklerine genel bakış

### 🎨 Modern UI/UX
- **Koyu/Açık Mod** - Yumuşak geçişli güzel temalar
- **Responsive Tasarım** - Masaüstü, tablet ve mobilde çalışır
- **Deploy Script Şablonları** - Node.js, Python, Go, .NET, Docker için hazır şablonlar

---

## 🏗️ Mimari

```
┌─────────────────────────────────────────────────────────────────┐
│                       deployes Mimarisi                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐         ┌──────────────┐                     │
│   │   Angular    │◄───────►│   Go API     │                     │
│   │   Frontend   │  REST   │   Backend    │                     │
│   │   (Port 4200)│  + WS   │  (Port 8080) │                     │
│   └──────────────┘         └──────┬───────┘                     │
│                                   │                              │
│                    ┌──────────────┼──────────────┐              │
│                    ▼              ▼              ▼              │
│           ┌────────────┐  ┌────────────┐  ┌────────────┐       │
│           │ PostgreSQL │  │   GitHub   │  │   Uzak     │       │
│           │  Veritabanı│  │  Webhook   │  │  Sunucular │       │
│           └────────────┘  └────────────┘  └────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Hızlı Başlangıç

### Gereksinimler

- **Go** 1.21+
- **Node.js** 18+
- **Docker** ve Docker Compose
- **Git**

### 1. Repository'yi Klonlayın

```bash
git clone https://github.com/aliicolak/deployes.git
cd deployes
```

### 2. Ortam Değişkenlerini Yapılandırın

```bash
cp .env.example .env
# .env dosyasını editleyip kendi güvenli değerlerinizi girin
```

> ⚠️ **`.env` dosyasını asla git'e commit etmeyin.** Gerekli değişkenler için [`.env.example`](.env.example) dosyasına bakın.

### 3. Veritabanını Başlatın

```bash
docker compose up -d

### 4. Backend'i Çalıştırın

```bash
go run ./cmd/api
```

### 5. Frontend'i Çalıştırın

```bash
cd web
npm install
npm start
```

### 6. Tarayıcıda Açın

[http://localhost:4200](http://localhost:4200) adresine gidin

---

## 📁 Proje Yapısı

```
deployes/
├── cmd/
│   └── api/                 # Uygulama giriş noktası
├── internal/
│   ├── application/         # İş mantığı servisleri
│   │   ├── deployment/      # Deployment servisi
│   │   ├── project/         # Proje yönetimi
│   │   ├── server/          # Sunucu yönetimi
│   │   └── webhook/         # Webhook işleme
│   ├── domain/              # Domain entity'leri ve interface'ler
│   ├── handlers/            # HTTP handler'ları
│   └── infrastructure/      # Veritabanı repository'leri
├── pkg/
│   └── utils/               # Yardımcı fonksiyonlar
├── web/                     # Angular frontend
│   └── src/
│       ├── app/
│       │   ├── core/        # Servisler, guard'lar, interceptor'lar
│       │   ├── features/    # Özellik bileşenleri
│       │   └── shared/      # Paylaşılan bileşenler
│       └── assets/
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## ⚙️ Yapılandırma

### Ortam Değişkenleri

| Değişken | Açıklama | Zorunlu | Örnek |
|----------|----------|---------|-------|
| `DATABASE_URL` | PostgreSQL bağlantı string'i | ✅ | `.env.example` dosyasına bakın |
| `JWT_SECRET` | JWT imzalama anahtarı (min 32 karakter) | ✅ | Rastgele anahtar üretin |
| `ENCRYPTION_KEY` | AES şifreleme anahtarı (tam 32 karakter) | ✅ | Rastgele anahtar üretin |
| `APP_PORT` | API sunucu portu | ❌ | `8080` (varsayılan) |
| `ALLOWED_ORIGINS` | CORS izinli origin'ler | ❌ | `http://localhost:4200` |

---

## 📖 Kullanım Kılavuzu

### Sunucu Ekleme

1. **Sunucular** sayfasına gidin
2. **"+ Yeni Sunucu"** butonuna tıklayın
3. Sunucu bilgilerini girin:
   - **Ad**: Sunucu için görünen ad
   - **Host**: IP adresi veya hostname
   - **Port**: SSH portu (varsayılan: 22)
   - **Kullanıcı Adı**: SSH kullanıcısı
   - **SSH Anahtarı**: Kimlik doğrulama için private key
4. **Kaydet** butonuna tıklayın

### Proje Ekleme

1. **Projeler** sayfasına gidin
2. **"+ Yeni Proje"** butonuna tıklayın
3. Proje bilgilerini girin:
   - **Proje Adı**: Görünen ad
   - **GitHub Repo URL**: Tam repository URL'i
   - **Branch**: Hedef branch (örn: `main`, `master`)
   - **Deploy Script**: Çalıştırılacak shell komutları
4. Private repo'lar için oluşturulan **Deploy Key**'i GitHub'a ekleyin

### Deployment Oluşturma

1. **Deployments** sayfasına gidin
2. **"+ Yeni Deployment"** butonuna tıklayın
3. **Proje** ve **Sunucu** seçin
4. **Deploy** butonuna tıklayın
5. Terminalde gerçek zamanlı logları izleyin

### Webhook Kurulumu

1. **Webhooks** sayfasına gidin
2. Proje/sunucu çifti için yeni webhook oluşturun
3. Oluşturulan **Webhook URL**'sini kopyalayın
4. GitHub'a ekleyin: `Settings → Webhooks → Add webhook`
5. Content type: `application/json`
6. Seçin: **Just the push event**

---

## 🛠️ Deploy Script Şablonları

### Node.js + PM2
```bash
#!/bin/bash
set -e
npm install
npm run build
pm2 reload ecosystem.config.js || pm2 restart all
```

### Docker Compose
```bash
#!/bin/bash
set -e
docker-compose pull
docker-compose up -d --build
docker system prune -f
```

### .NET / ASP.NET Core
```bash
#!/bin/bash
set -e
dotnet restore
dotnet publish -c Release -o ./publish
systemctl restart myapp.service
```

---

## 🔒 Güvenlik En İyi Uygulamaları

1. **Güçlü secret'lar kullanın** - Rastgele, uzun JWT ve şifreleme anahtarları oluşturun
2. **Sadece SSH anahtarları** - Asla parola tabanlı SSH authentication kullanmayın
3. **Firewall kuralları** - 8080 ve 4200 portlarına erişimi kısıtlayın
4. **HTTPS** - SSL sertifikalarıyla reverse proxy (nginx) kullanın
5. **Düzenli güncellemeler** - Bağımlılıkları ve Docker image'larını güncel tutun

---

## 🤝 Katkıda Bulunma

Katkılarınızı bekliyoruz! Lütfen Pull Request göndermekten çekinmeyin.

1. Repository'yi fork edin
2. Feature branch'inizi oluşturun (`git checkout -b feature/harika-ozellik`)
3. Değişikliklerinizi commit edin (`git commit -m 'Harika özellik ekle'`)
4. Branch'e push edin (`git push origin feature/harika-ozellik`)
5. Pull Request açın

---

## 📄 Lisans

Bu proje MIT Lisansı altında lisanslanmıştır - detaylar için [LICENSE](LICENSE) dosyasına bakın.

---

## 👨‍💻 Katkıda Bulunanlar

<a href="https://github.com/aliicolak">
  <img src="https://github.com/aliicolak.png?size=100" width="100" height="100" alt="Ali ÇOLAK" style="border-radius: 50%; margin-right: 20px;" />
</a>
<a href="https://github.com/Alpersahin11">
  <img src="https://github.com/Alpersahin11.png?size=100" width="100" height="100" alt="Alper ŞAHİN" style="border-radius: 50%;" />
</a>

**Ali ÇOLAK** ([@aliicolak](https://github.com/aliicolak))  
**Alper ŞAHİN** ([@Alpersahin11](https://github.com/Alpersahin11))
---

<div align="center">

**⭐ Faydalı bulduysanız bu repository'ye yıldız verin!**

Go ve Angular ile ❤️ yapıldı

</div>
