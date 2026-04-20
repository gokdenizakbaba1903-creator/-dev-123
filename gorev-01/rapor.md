# Gorev 01: OWASP Top 10 Mapping

Bu görevde popüler bir repo üzerinden OWASP Top 10 2021 risklerinin nasıl eşleştiğini analiz edeceğiz.

## Analiz Edilen Senaryo: E-Commerce API

Aşağıdaki liste, modern web uygulamalarında karşılaşılan risklerin OWASP Top 10 ile eşleşmesini göstermektedir:

1.  **A01:2021-Broken Access Control**: Kullanıcının başka bir kullanıcının sepetine (`/api/cart/123`) erişebilmesi.
2.  **A02:2021-Cryptographic Failures**: Şifrelerin veritabanında "plaintext" veya MD5 gibi zayıf algoritmalarla saklanması.
3.  **A03:2021-Injection**: `search?q=' OR 1=1 --` gibi inputlarla SQL Injection yapılması.
4.  **A04:2021-Insecure Design**: Ödeme sayfasında client-side fiyat kontrolü yapılması.
5.  **A05:2021-Security Misconfiguration**: Debug modun production'da açık kalması veya default şifrelerin değiştirilmemesi.
6.  **A06:2021-Vulnerable and Outdated Components**: Log4j gibi bilinen CVE'li kütüphanelerin kullanılması.
7.  **A07:2021-Identification and Authentication Failures**: Brute-force saldırılarına karşı koruma olmaması (rate limiting eksikliği).
8.  **A08:2021-Software and Data Integrity Failures**: Güvenilir olmayan kaynaklardan gelen verilerin (untrusted data) deserialize edilmesi.
9.  **A09:2021-Security Logging and Monitoring Failures**: Login denemelerinin veya kritik hataların loglanmaması.
10. **A10:2021-Server-Side Request Forgery (SSRF)**: Kullanıcıdan alınan URL ile sunucunun iç ağdaki kaynaklara (örn: metadata API) erişmesi.

## Tespit Edilen Zafiyyet Örneği
E-commerce reposunda `/api/user/:id` endpoint'inin JWT içindeki user ID ile eşleşip eşleşmediği kontrol edilmiyorsa, bu **Broken Access Control** zafiyetidir.

## Çözüm
Her request'te kullanıcının yetkisi (`RBAC` veya `ABAC`) middleware seviyesinde kontrol edilmelidir.
