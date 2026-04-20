# Gorev 02: SQL Injection Analizi

## 1. Zafiyetli Kod (Node.js + Express)
Kullanıcı girdisini doğrudan SQL sorgusuna ekleyen kod örneği:

```javascript
// vulnerable.js
app.get('/user', (req, res) => {
    const userId = req.query.id;
    const query = "SELECT * FROM users WHERE id = " + userId; // TEHLİKELİ!
    db.query(query, (err, result) => {
        if (err) throw err;
        res.send(result);
    });
});
```

## 2. Payload Örnekleri
*   **Veri Çekme:** `?id=1 OR 1=1` (Tüm kullanıcıları getirir)
*   **Veritabanı Yapısını Öğrenme:** `?id=1 UNION SELECT 1,2,table_name FROM information_schema.tables`
*   **Dosya Okuma:** `?id=1 UNION SELECT 1,2,load_file('/etc/passwd')`

## 3. Savunma (Fix)
**Parameterized Queries** kullanarak saldırıları engelleme:

```javascript
// secure.js
app.get('/user', (req, res) => {
    const userId = req.query.id;
    const query = "SELECT * FROM users WHERE id = ?"; // Placeholder kullan
    db.query(query, [userId], (err, result) => { // Inputu dizi olarak gönder
        if (err) throw err;
        res.send(result);
    });
});
```

## 4. Rapor
1.  **Adım:** Saldırgan tarayıcı üzerinden `id=1; DROP TABLE users;--` payload'unu gönderir.
2.  **Adım:** Sunucu string birleştirme yaptığı için sorgu `SELECT * FROM users WHERE id = 1; DROP TABLE users;--` haline gelir.
3.  **Adım:** Veritabanı tablosu silinir.
4.  **Çözüm:** Input her zaman veriden (data) ayrılmalı, sorgu motoru inputu komut olarak değil sadece parametre olarak algılamalıdır.
