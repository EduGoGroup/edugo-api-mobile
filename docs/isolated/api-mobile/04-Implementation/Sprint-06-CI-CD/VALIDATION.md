# Validación del Sprint 06

## Validación de CI/CD

### 1. Pipeline Verde en GitHub
```bash
# Ver status de Actions
gh run list --workflow=ci.yml --limit 5

# Ver último run
gh run view

# Logs del último run
gh run view --log
```

**Criterio de éxito:** ✅ Estado: Success

### 2. Tests Automáticos Ejecutándose
```bash
# Verificar que workflow ejecuta tests
gh run view --log | grep "go test"
# Esperado: Encuentra líneas con "go test"
```

### 3. Docker Build Exitoso
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Build local
docker build -t edugo-api-mobile:test .

# Verificar imagen creada
docker images | grep edugo-api-mobile

# Run container
docker run -d -p 8080:8080 --name test-api edugo-api-mobile:test

# Health check
curl http://localhost:8080/health

# Cleanup
docker stop test-api
docker rm test-api
```

---

## Criterios de Éxito

- [ ] GitHub Actions workflow ejecuta en cada push
- [ ] Tests automáticos pasando en CI
- [ ] Linting automático sin errores
- [ ] Docker build exitoso
- [ ] README actualizado

---

**Sprint:** 06/06
