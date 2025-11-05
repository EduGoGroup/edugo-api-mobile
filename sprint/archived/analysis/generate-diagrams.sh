#!/bin/bash

# Script para generar imágenes PNG de diagramas Mermaid
# Autor: Claude Code
# Fecha: 2025-11-04

set -e

# Obtener directorio del script de forma dinámica
ANALYSIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMP_DIR="${ANALYSIS_DIR}/temp_mermaid"

echo "🎨 Generando imágenes de diagramas Mermaid..."
echo ""

# Crear directorio temporal
mkdir -p "${TEMP_DIR}"

# Función para extraer bloques Mermaid de un archivo
extract_mermaid() {
    local file=$1
    local output_prefix=$2
    local counter=1

    echo "📄 Procesando: $(basename $file)"

    # Usar awk para extraer bloques mermaid
    awk '
        /```mermaid/ { in_block=1; block=""; next }
        in_block {
            if (/```/) {
                in_block=0
                print block > "'${TEMP_DIR}'/'${output_prefix}'-" ++counter ".mmd"
                next
            }
            block = block $0 "\n"
        }
    ' "$file"

    return $counter
}

# Procesar architecture-phase-2.md
echo "═══════════════════════════════════════"
extract_mermaid "${ANALYSIS_DIR}/architecture-phase-2.md" "architecture"
ARCH_COUNT=$?

# Procesar data-model-phase-2.md
extract_mermaid "${ANALYSIS_DIR}/data-model-phase-2.md" "data-model"
DATA_COUNT=$?

# Procesar process-diagram-phase-2.md
extract_mermaid "${ANALYSIS_DIR}/process-diagram-phase-2.md" "process"
PROCESS_COUNT=$?

# Procesar readme-phase-2.md
extract_mermaid "${ANALYSIS_DIR}/readme-phase-2.md" "readme"
README_COUNT=$?

echo "═══════════════════════════════════════"
echo ""
echo "📊 Diagramas extraídos:"
ls -la "${TEMP_DIR}"/*.mmd 2>/dev/null || echo "  (ninguno encontrado)"
echo ""

# Convertir cada archivo .mmd a PNG
echo "🖼️  Convirtiendo diagramas a PNG..."
echo ""

for mmd_file in "${TEMP_DIR}"/*.mmd; do
    if [ -f "$mmd_file" ]; then
        filename=$(basename "$mmd_file" .mmd)
        output_file="${ANALYSIS_DIR}/${filename}.png"

        echo "  → Generando: ${filename}.png"

        # Convertir con mmdc
        mmdc -i "$mmd_file" -o "$output_file" -b transparent -t dark 2>/dev/null || {
            echo "    ⚠️  Error con tema dark, intentando con tema default..."
            mmdc -i "$mmd_file" -o "$output_file" -b transparent
        }
    fi
done

# Limpiar archivos temporales
echo ""
echo "🧹 Limpiando archivos temporales..."
rm -rf "${TEMP_DIR}"

echo ""
echo "✅ ¡Proceso completado!"
echo ""
echo "📁 Imágenes generadas en:"
echo "   ${ANALYSIS_DIR}/"
echo ""
ls -lh "${ANALYSIS_DIR}"/*.png 2>/dev/null || echo "  (no se generaron imágenes)"
echo ""
