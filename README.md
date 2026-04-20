# Web Güvenliği Lab Projesi (seri-3-lab) 🚀

Bu proje, modern web güvenliği açıklarını tespit etmek, anlamak ve önlemek için tasarlanmış kapsamlı bir laboratuvardır.

## 🛠️ Proje Bileşenleri

1.  **SecScan AI (Fullstack):** Modern, premium arayüzlü bir güvenlik tarayıcı uygulaması.
    - **Back-end (Go):** XSS, SQLi, Port, Header ve CVE tarama motoru.
    - **Front-end (Next.js):** Gerçek zamanlı dashboard, tarama geçmişi ve detaylı raporlama.
2.  **Vulnerable Target App (Go):** Bilerek zafiyet eklenmiş (SQLi, XSS, File Exposure) bir test uygulaması.
3.  **Güvenlik Laboratuvarı (10 Görev):** OWASP Top 10 temelli raporlar ve dokümantasyonlar.
4.  **Security Cheat Sheet:** Yaygın saldırı payloadları ve savunma rehberi.

## 🚀 Hızlı Başlangıç

### 1. Sistem Gereksinimleri
- Docker & Docker Compose
- Node.js (Opsiyonel, yerel geliştirme için)
- Go (Opsiyonel, yerel geliştirme için)

### 2. Uygulamaları Başlat

```bash
# Proje dizinine girin
cd seri-3-lab

# 1. SecScan Uygulamasını Başlat (Scanner + UI)
cd secscan
docker-compose up --build

# 2. (Opsiyonel) Hedef Uygulamayı Başlat (Ayrı bir terminalde)
cd ../gorev-hedef
go run main.go
```

## 📁 Klasör Yapısı

```text
seri-3-lab/
├── secscan/           # Ana Scanner Uygulaması
│   ├── backend/       # Go Scanners + SSE + PDF Report
│   └── frontend/      # Next.js Premium Dashboard
├── gorev-hedef/       # Zafiyetli Test Uygulaması (9090 portu)
├── gorev-01-10/       # OWASP Top 10 Analiz Raporları
├── cheat-sheet.md     # Saldırı/Savunma Rehberi
└── README.md
```

## 🧪 SecScan Premium Özellikleri
- **Interactive Dashboard:** Radar grafiği ile risk dağılımı görselleştirmesi.
- **Scan History:** Tarama sonuçlarınız `localStorage` üzerinde kaydedilir.
- **Deep Scanning:** Gerçek payloadlar ile aktif XSS ve SQLi tespiti.
- **Professional Reports:** Tek tıkla indirilebilir teknik PDF raporlar.
- **SSE Streaming:** Tarama ilerlemesini canlı olarak takip edin.

## ⚠️ Uyarı
Bu proje sadece **eğitim amaçlıdır**. Kendi kontrolünüzde olmayan sistemlere karşı tarama yapmayınız.
