# Web Güvenliği Lab: Cheat Sheet 🛡️

Bu rehber, projenizdeki zaafiyetleri anlamanız ve savunma mekanizmalarını öğrenmeniz için hazırlanmıştır.

---

## 1. SQL Injection (SQLi)
**Nedir?** Saldırganın, uygulamanın veritabanı sorgularına müdahale etmesidir.

### 🧪 Test Payloadları
- `' OR 1=1 --`: Tüm kullanıcıları listelemek için.
- `' UNION SELECT NULL, name, secret FROM users --`: Gizli tabloları çekmek için.
- `'; DROP TABLE users --`: (Tehlikeli!) Tabloyu silmek için.

### 🛡️ Savunma
- **Prepared Statements (Parametreli Sorgular):** Veriyi koddan ayırın.
```go
db.Query("SELECT name FROM users WHERE id = ?", id)
```

---

## 2. Cross-Site Scripting (XSS)
**Nedir?** Saldırganın, kurbanın tarayıcısında zararlı JavaScript çalıştırmasıdır.

### 🧪 Test Payloadları
- `<script>alert(1)</script>`
- `<svg/onload=alert('XSS')>`
- `<img src=x onerror=prompt(document.cookie)>`

### 🛡️ Savunma
- **Output Encoding:** HTML çıktılarını encode edin (örn. `<` -> `&lt;`).
- **Content Security Policy (CSP):** Sadece güvenilir kaynaklardan script çalışmasına izin verin.

---

## 3. Server-Side Request Forgery (SSRF)
**Nedir?** Sunucunun, saldırgan tarafından kontrol edilen bir URL'ye (genellikle iç ağdaki) istek atmasıdır.

### 🧪 Test Senaryosu
- `http://169.254.169.254/latest/meta-data/`: AWS metadata'sına erişim.
- `http://localhost:8080/admin`: İç panel erişimi.

### 🛡️ Savunma
- **Deny-list (Private IPs):** İç ağ IP bloklarını engelleyin (Projemizdeki `isPrivateIP` fonksiyonu gibi).

---

## 4. Güvenlik Header'ları
**SecScan** tarafından denetlenen önemli header'lar:
- **X-Frame-Options:** Clickjacking'i önler (DENY/SAMEORIGIN).
- **HSTS:** Bağlantının her zaman HTTPS olmasını zorlar.
- **X-Content-Type-Options:** Tarayıcıların MIME sniffing yapmasını engeller (nosniff).
