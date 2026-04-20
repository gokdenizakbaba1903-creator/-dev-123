# Gorev 06: OAuth 2.0 + PKCE Analizi

## 1. Neden PKCE? (Proof Key for Code Exchange)
Native ve Single Page uygulamalarda (SPA) `Client Secret` güvenli bir şekilde saklanamaz. PKCE, yetkilendirme kodunun (auth code) çalınsa bile saldırgan tarafından kullanılmasını engeller.

## 2. Adım Adım Akış
1.  **Code Verifier Oluştur:** Rastgele bir string (`secret-string-123`).
2.  **Code Challenge Oluştur:** Verifier'in SHA256 hash'i.
3.  **Authorize İsteği:** Challenge'ı Google'a gönder.
    `https://accounts.google.com/o/oauth2/v2/auth?code_challenge=xyz&code_challenge_method=S256...`
4.  **Token İsteği:** Auth code ile birlikte orijinal `Verifier`'ı gönder.
    `POST /token { code: '...', code_verifier: 'secret-string-123' }`

## 3. Demo Kod (Node.js)
```javascript
const crypto = require('crypto');

// 1. Verifier (Rastgele string)
const base64URLEncode = (str) => {
    return str.toString('base64').replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
};
const verifier = base64URLEncode(crypto.randomBytes(32));

// 2. Challenge (SHA256 Hash)
const challenge = base64URLEncode(crypto.createHash('sha256').update(verifier).digest());

console.log("Verifier:", verifier);
console.log("Challenge:", challenge);
```

## 4. Analiz
Google Login implementasyonunda `response_type=code` ve PKCE kullanmak, `implicit flow` (token'ın URL'de dönmesi) kullanımından çok daha güvenlidir.
