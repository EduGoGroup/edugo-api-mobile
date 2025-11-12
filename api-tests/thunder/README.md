# Thunder Client Collection

Este directorio contiene una colección exportada para Thunder Client (extensión de VSCode).

## ⚠️ Nota Importante

Este proyecto usa **httpyac** como herramienta principal de testing. Los archivos `.http` en `api-tests/requests/` son compatibles tanto con httpyac como con Thunder Client.

## 🚀 Uso con Thunder Client

### Opción 1: Importar colección existente

1. Instala Thunder Client en VSCode
2. Abre Thunder Client desde la barra lateral
3. Click en "Import" → Selecciona `edugo-postman-collection.json`

### Opción 2: Usar archivos .http directamente

Thunder Client puede ejecutar archivos `.http` directamente:

1. Abre cualquier archivo `.http` en `api-tests/requests/`
2. Thunder Client detectará las peticiones
3. Ejecuta desde Thunder Client UI

## 📝 Mantenimiento

Si necesitas regenerar la colección:

1. Exporta desde Thunder Client: Collections → Export → Postman v2.1
2. Guarda como `edugo-postman-collection.json` en este directorio

## 💡 Recomendación

Para testing local, usa **httpyac** desde la línea de comandos (ver `api-tests/README.md`).

Thunder Client es útil para:
- Testing visual e interactivo
- Usuarios que prefieren UI sobre CLI
- Compartir colecciones con el equipo

---

**Última actualización**: 11 de noviembre de 2025
