# Gorev 10: Security Headers (Helmet ile A+)

## 1. Kritik Güvenlik Başlıkları
Web güvenliğini sunucu seviyesinde artırmak için kullanılan HTTP header'ları:

1.  **Strict-Transport-Security (HSTS):** Sadece HTTPS kullanımını zorunlu kılar.
2.  **X-Content-Type-Options:** `nosniff` değeri ile MIME-type sniffing saldırılarını engeller.
3.  **X-Frame-Options:** `DENY` veya `SAMEORIGIN` ile Clickjacking saldırılarını önler.
4.  **Referrer-Policy:** Hassas URL bilgilerinin diğer sitelerle paylaşımını sınırlar.

## 2. Express Üzerinde Uygulama (Helmet)

```javascript
const express = require('express');
const helmet = require('helmet');

const app = express();

// Tüm standart güvenlik başlıklarını ekler
app.use(helmet());

// Özel yapılandırma
app.use(helmet.hidePoweredBy()); // X-Powered-By headerını gizler
```

## 3. Analiz (SecurityScorecard / SSLLabs)
Başlıklar eklendikten sonra tarayıcı üzerinden kontrol:
*   `curl -I http://localhost:3000` komutu ile dönen header'lar incelenmelidir.
*   `A+` skoru için tüm zorunlu başlıkların (`CSP`, `HSTS`, `XCTO`, `XFO`) varlığı kontrol edilmelidir.

## 4. Rapor
Sunucu yanıtlarında `X-Powered-By: Express` gibi bilgilerin sızdırılması, saldırganın versiyona özel zafiyetleri (exploit) taramasını kolaylaştırır. Helmet bu bilgiyi otomatik gizler.
