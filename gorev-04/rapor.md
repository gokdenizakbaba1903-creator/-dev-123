# Gorev 04: CSRF (Cross-Site Request Forgery)

## 1. Senaryo
Kullanıcının banka hesabından para transferi yapan bir endpoint: `POST /api/transfer`.
Eğer uygulama sadece session cookie'sine güveniyorsa, saldırgan kendi sitesi üzerinden bir form ile bu isteği kullanıcı adına tetikleyebilir.

**Saldırgan Formu (hacker.com):**
```html
<form action="https://banka.com/api/transfer" method="POST" id="csrf-form">
    <input type="hidden" name="to" value="hacker_account">
    <input type="hidden" name="amount" value="10000">
</form>
<script>document.getElementById('csrf-form').submit();</script>
```

## 2. Savunma: CSRF Token (csurf)
Her POST isteği için sunucu tarafından üretilen benzersiz bir token kontrol edilir.

```javascript
const csrf = require('csurf');
const csrfProtection = csrf({ cookie: true });

app.post('/api/transfer', csrfProtection, (req, res) => {
    res.send('Transfer Başarılı!');
});
```

## 3. Frontend Uygulama
Form içine sunucudan gelen token eklenmelidir:
```html
<input type="hidden" name="_csrf" value="<%= csrfToken %>">
```

## 4. Analiz
1.  Saldırgan token'a sahip olmadığı için `POST` isteği `403 Forbidden` alacaktır.
2.  Modern tarayıcılarda `SameSite=Lax` cookie ayarı da ek bir koruma sağlar.
