# Gorev 07: Nmap ve OWASP ZAP Analizi

## 1. Nmap Taraması
Sunucudaki açık portları ve servisleri tespit etmek için kullanılır.

**Komut:**
`nmap -sV -sC -p- 127.0.0.1`

**Sonuç Analizi:**
*   `80/tcp open http`: Apache 2.4.41 (Eski sürüm, CVE-2021-41773 potansiyeli).
*   `22/tcp open ssh`: OpenSSH 8.2 (Brute-force'a açık olabilir).

## 2. OWASP ZAP (DAST)
Web uygulamasındaki çalışma zamanı hatalarını (XSS, SQLi vb.) bulmak için kullanılır.

**Tarama Adımları:**
1.  **Automation Framework:** `zap.sh -cmd -autorun scan.yaml`
2.  **Bulgular:**
    *   **High:** SQL Injection (Search endpoint).
    *   **Medium:** Missing Anti-CSRF Tokens.
    *   **Low:** X-Content-Type-Options Header Missing.

## 3. Raporlama
ZAP taraması sonunda üretilen `report.html` üzerinden kritik bulgular Jira ortamına aktarılmalı ve geliştiriciler tarafından fixed edilmelidir.

## 4. Karşılaştırma
*   **Nmap:** Network katmanındaki zafiyetleri bulur (Açık portlar, OS versiyonları).
*   **ZAP:** Uygulama (Application) katmanındaki mantıksal zafiyetleri bulur.
