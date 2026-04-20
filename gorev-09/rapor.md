# Gorev 09: SBOM ve Trivy (SCA) Analizi

## 1. SBOM (Software Bill of Materials) Nedir?
Uygulamanın kullandığı tüm kütüphanelerin ve bağımlılıkların (transitive dependencies dahil) listesidir.

**Sybth ile SBOM Üretme:**
`syft seri-3-lab-backend -o cyclonedx-json > sbom.json`

## 2. Trivy ile Zafiyet Taraması
Docker imajlarını ve kütüphaneleri bilinen CVE'lere karşı tarar (Software Composition Analysis - SCA).

**Komut:**
`trivy image --severity HIGH,CRITICAL node:14`

## 3. Örnek CVE Bulgu ve Fix
**Bulunan:** `express@4.16.0` -> `CVE-2019-10744` (Prototype Pollution)
**Fix:** `package.json` dosyasında versiyonu `4.17.1` veya üzerine yükselt.

```bash
npm install express@latest
```

## 4. Analiz
Geliştirme sürecine Trivy'yi dahil etmek, zafiyetli kütüphanelerin production ortamına çıkmasını engelleyerek **Supply Chain Attacks** riskini minimize eder.
