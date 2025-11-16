# Dependencias del Sprint 04

## Dependencias Técnicas

- [ ] Sprint-03 completado (repositorios)
- [ ] Gin framework
- [ ] Go validator
- [ ] Swag

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

go get github.com/gin-gonic/gin@v1.10.0
go get github.com/go-playground/validator/v10@v10.16.0
go get github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/gin-swagger@v1.6.0
go get github.com/swaggo/files@v1.0.1

go mod tidy
```

## Verificación

```bash
# Verificar Gin
go list -m github.com/gin-gonic/gin

# Verificar swag instalado
swag --version
```

---

**Sprint:** 04/06
