# Gorev 05: JWT (JSON Web Token) Güvenlik Denetimi

## 1. Audit Checklist
JWT kullanırken şu kontroller yapılmalıdır:

*   [ ] **Algorithm Check:** `alg`: "none" saldırılarına karşı kütüphanede algoritma kısıtlanmış mı?
*   [ ] **Expiration:** `exp` (son kullanma tarihi) makul bir sürede mi? (15-60 dk)
*   [ ] **Signature:** Token imzası (`HMAC` veya `RSA`) sunucu tarafında doğrulanıyor mu?
*   [ ] **Secret Strength:** HMAC secret en az 256-bit karmaşıklığında mı?
*   [ ] **Payload:** Token içinde hassas veri (şifre, TC no) bulunuyor mu? (JWT'ler base64'tür, herkes okuyabilir).

## 2. Hatalı Kullanım (None Attack)
```javascript
const jwt = require('jsonwebtoken');
// Tehlikeli: Algoritma kontrolü yok
const decoded = jwt.decode(token); // Sadece decode eder, verify etmez!
```

## 3. Güvenli Kullanım (Verify)
```javascript
const jwt = require('jsonwebtoken');

try {
    const verified = jwt.verify(token, process.env.JWT_SECRET, {
        algorithms: ['HS256'] // Sadece bu algoritmayı kabul et
    });
} catch (err) {
    // Hatalı token
}
```

## 4. Analiz
JWT bir "session storage" değildir. Eğer bir kullanıcıyı banlamak gerekiyorsa JWT'yi geçersiz kılmak için **Blacklist** (Redis vb.) veya **Refresh Token** mekanizması kullanılmalıdır.
