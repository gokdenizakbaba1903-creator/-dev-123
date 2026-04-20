# Gorev 08: Semgrep (SAST) CI Pipeline

## 1. Semgrep Nedir?
Static Analysis Security Testing (SAST) aracıdır. Kod yazılırken güvenlik açıklarını (hardcoded secrets, insecure functions) bulur.

## 2. GitHub Actions Konfigürasyonu
Aşağıdaki dosya her PR geldiğinde kodu tarar.

```yaml
# .github/workflows/semgrep.yml
name: Semgrep Analysis

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  semgrep:
    runs-on: ubuntu-latest
    container:
      image: returntocorp/semgrep
    steps:
      - uses: actions/checkout@v3
      - run: semgrep scan --config auto --error
```

## 3. Örnek Bir Semgrep Kuralı
Semgrep aşağıdaki deseni bulduğunda uyarı verir:
```python
# pattern matching
- pattern: "requests.get(..., verify=False)"
  message: "SSL verification is disabled. This is insecure!"
  severity: ERROR
```

## 4. Analiz
SAST araçları kodu çalıştırmadan analiz ettiği için hızlıdır ancak "false positive" (yanlış alarm) oranı daha yüksek olabilir. Bu yüzden kural setleri (rule sets) projeye göre özelleştirilmelidir.
