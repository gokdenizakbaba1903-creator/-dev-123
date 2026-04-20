# Web Güvenliği Lab: SecScan AI (Neon Shadow Edition) 🚀

Bu proje, modern web güvenliği açıklarını tespit etmek, görselleştirmek ve önlemek için tasarlanmış profesyonel bir siber güvenlik laboratuvarıdır.

## 🛠️ Proje Bileşenleri

1.  **SecScan AI (Fullstack):** Neon Cyber temalı, yüksek performanslı güvenlik tarayıcı uygulaması.
    - **Back-end (Go):** XSS, SQLi, Port, Header, CVE ve **Subdomain Discovery** tarama motoru.
    - **Front-end (Next.js):** Gerçek zamanlı dashboard, tarama geçmişi (zaman damgalı) ve detaylı raporlama.
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
├── secscan/           # Ana Scanner Uygulaması (Port: 3000)
│   ├── backend/       # Go Scanners + SSE + PDF Report
│   └── frontend/      # Next.js Neon Shadow Dashboard
├── gorev-hedef/       # Zafiyetli Test Uygulaması (Port: 9090)
├── gorev-01-10/       # OWASP Top 10 Analiz Raporları
├── cheat-sheet.md     # Saldırı/Savunma Rehberi
└── README.md
```

## 🧪 SecScan Premium Özellikleri
- **Neon Shadow UI:** Cyberpunk estetiğiyle tasarlanmış, kristal netliğinde dashboard.
- **Enhanced Scan History:** Tarama sonuçlarınız zaman damgasıyla `localStorage` (cache) üzerinde saklanır.
- **8 Aktif Modül:** Port, Header, TLS, Fuzzer, XSS, SQLi, CVE ve Subdomain Discovery.
- **Remediation Intelligence:** Her bulgu için teknik çözüm önerileri.
- **Technical PDF Reports:** Tek tıkla indirilebilir profesyonel denetim raporları.

## ⚠️ Uyarı
Bu proje sadece **eğitim amaçlıdır**. Kendi kontrolünüzde olmayan sistemlere karşı tarama yapmayınız.
