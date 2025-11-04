#!/usr/bin/env python3
"""
Script para generar imágenes PNG de diagramas Mermaid
Autor: Claude Code
Fecha: 2025-11-04
"""

import os
import re
import subprocess
from pathlib import Path

# Configuración
ANALYSIS_DIR = Path("/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/sprint/current/analysis")
TEMP_DIR = ANALYSIS_DIR / "temp_mermaid"

def extract_mermaid_blocks(file_path, output_prefix):
    """Extrae bloques Mermaid de un archivo Markdown"""
    print(f"📄 Procesando: {file_path.name}")

    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # Patrón para encontrar bloques ```mermaid ... ```
    pattern = r'```mermaid\n(.*?)```'
    blocks = re.findall(pattern, content, re.DOTALL)

    # Guardar cada bloque en un archivo temporal
    for i, block in enumerate(blocks, start=1):
        output_file = TEMP_DIR / f"{output_prefix}-{i}.mmd"
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(block.strip())
        print(f"  ✓ Extraído diagrama {i}")

    return len(blocks)

def convert_to_png(mmd_file):
    """Convierte un archivo .mmd a PNG usando mmdc"""
    filename = mmd_file.stem
    output_file = ANALYSIS_DIR / f"{filename}.png"

    print(f"  → Generando: {filename}.png")

    try:
        # Intentar con tema dark y fondo transparente
        subprocess.run(
            ['mmdc', '-i', str(mmd_file), '-o', str(output_file), '-b', 'transparent', '-t', 'dark'],
            check=True,
            capture_output=True,
            text=True
        )
    except subprocess.CalledProcessError:
        try:
            # Si falla, intentar con tema default
            print(f"    ⚠️  Error con tema dark, intentando con tema default...")
            subprocess.run(
                ['mmdc', '-i', str(mmd_file), '-o', str(output_file), '-b', 'transparent'],
                check=True,
                capture_output=True,
                text=True
            )
        except subprocess.CalledProcessError as e:
            print(f"    ❌ Error al convertir: {e.stderr}")
            return False

    return True

def main():
    print("🎨 Generando imágenes de diagramas Mermaid...")
    print()

    # Crear directorio temporal
    TEMP_DIR.mkdir(exist_ok=True)

    # Archivos a procesar
    files_to_process = [
        ('architecture-phase-2.md', 'architecture'),
        ('data-model-phase-2.md', 'data-model'),
        ('process-diagram-phase-2.md', 'process'),
        ('readme-phase-2.md', 'readme'),
    ]

    total_diagrams = 0

    print("═══════════════════════════════════════")

    # Extraer diagramas de cada archivo
    for filename, prefix in files_to_process:
        file_path = ANALYSIS_DIR / filename
        if file_path.exists():
            count = extract_mermaid_blocks(file_path, prefix)
            total_diagrams += count
        else:
            print(f"⚠️  Archivo no encontrado: {filename}")

    print("═══════════════════════════════════════")
    print()
    print(f"📊 Total de diagramas extraídos: {total_diagrams}")

    # Listar archivos temporales
    mmd_files = list(TEMP_DIR.glob('*.mmd'))
    if mmd_files:
        print(f"📁 Archivos temporales creados: {len(mmd_files)}")
        for mmd in sorted(mmd_files):
            print(f"   - {mmd.name}")
    else:
        print("⚠️  No se encontraron diagramas Mermaid")
        return

    print()
    print("🖼️  Convirtiendo diagramas a PNG...")
    print()

    # Convertir cada archivo .mmd a PNG
    success_count = 0
    for mmd_file in sorted(mmd_files):
        if convert_to_png(mmd_file):
            success_count += 1

    # Limpiar archivos temporales
    print()
    print("🧹 Limpiando archivos temporales...")
    for mmd_file in mmd_files:
        mmd_file.unlink()
    TEMP_DIR.rmdir()

    print()
    print("✅ ¡Proceso completado!")
    print()
    print(f"📈 Estadísticas:")
    print(f"   - Diagramas encontrados: {total_diagrams}")
    print(f"   - Imágenes generadas: {success_count}")
    print()
    print("📁 Imágenes generadas en:")
    print(f"   {ANALYSIS_DIR}/")
    print()

    # Listar imágenes generadas
    png_files = sorted(ANALYSIS_DIR.glob('*.png'))
    if png_files:
        print("🖼️  Archivos PNG creados:")
        for png in png_files:
            size = png.stat().st_size / 1024  # KB
            print(f"   - {png.name} ({size:.1f} KB)")
    else:
        print("⚠️  No se generaron imágenes PNG")

    print()

if __name__ == '__main__':
    main()
