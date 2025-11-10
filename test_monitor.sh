#!/bin/bash

# Script para monitorear tiempos de ejecución de tests
# Detecta problemas de performance con contenedores

echo "╔══════════════════════════════════════════════════════════════════════╗"
echo "║         🔍 MONITOR DE TESTS - ANÁLISIS DE PERFORMANCE               ║"
echo "╚══════════════════════════════════════════════════════════════════════╝"
echo ""
echo "📅 Fecha: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Función para formatear tiempo
format_time() {
    local seconds=$1
    printf "%02d:%02d" $((seconds/60)) $((seconds%60))
}

# Limpiar contenedores previos
echo "🧹 Limpiando contenedores previos..."
docker ps -a | grep -E "postgres|mongo|rabbitmq" | awk '{print $1}' | xargs -r docker rm -f 2>/dev/null
echo ""

# Fase 1: Tests Unitarios
echo "═══════════════════════════════════════════════════════════════════════"
echo "📦 FASE 1: TESTS UNITARIOS"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""

START_UNIT=$(date +%s)
make test-unit 2>&1 | tee /tmp/unit_tests.log
UNIT_EXIT_CODE=${PIPESTATUS[0]}
END_UNIT=$(date +%s)
UNIT_TIME=$((END_UNIT - START_UNIT))

echo ""
echo "⏱️  Tiempo tests unitarios: $(format_time $UNIT_TIME) (${UNIT_TIME}s)"
echo "📊 Estado: $([ $UNIT_EXIT_CODE -eq 0 ] && echo '✅ PASS' || echo '❌ FAIL')"
echo ""

# Extraer estadísticas de tests unitarios
UNIT_TESTS=$(grep -oP '\d+(?= tests)' /tmp/unit_tests.log | tail -1)
UNIT_PASS=$(grep -oP 'ok\s+' /tmp/unit_tests.log | wc -l)
echo "   Tests ejecutados: ${UNIT_TESTS:-N/A}"
echo "   Paquetes OK: ${UNIT_PASS}"
echo ""

# Fase 2: Tests de Integración (con monitoreo de contenedores)
echo "═══════════════════════════════════════════════════════════════════════"
echo "🐳 FASE 2: TESTS DE INTEGRACIÓN (Monitoreo de Contenedores)"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""

# Iniciar monitoreo de contenedores en background
(
    while true; do
        CONTAINER_COUNT=$(docker ps | grep -cE "postgres|mongo|rabbitmq" || echo "0")
        echo "[$(date '+%H:%M:%S')] Contenedores activos: $CONTAINER_COUNT" >> /tmp/container_monitor.log
        sleep 2
    done
) &
MONITOR_PID=$!

START_INTEGRATION=$(date +%s)
make test-integration 2>&1 | tee /tmp/integration_tests.log
INTEGRATION_EXIT_CODE=${PIPESTATUS[0]}
END_INTEGRATION=$(date +%s)
INTEGRATION_TIME=$((END_INTEGRATION - START_INTEGRATION))

# Detener monitoreo
kill $MONITOR_PID 2>/dev/null

echo ""
echo "⏱️  Tiempo tests integración: $(format_time $INTEGRATION_TIME) (${INTEGRATION_TIME}s)"
echo "📊 Estado: $([ $INTEGRATION_EXIT_CODE -eq 0 ] && echo '✅ PASS' || echo '❌ FAIL')"
echo ""

# Extraer estadísticas de tests de integración
INTEGRATION_TESTS=$(grep -oP '\d+(?= tests)' /tmp/integration_tests.log | tail -1)
INTEGRATION_PASS=$(grep -c "PASS" /tmp/integration_tests.log || echo "0")
INTEGRATION_FAIL=$(grep -c "FAIL" /tmp/integration_tests.log || echo "0")
echo "   Tests ejecutados: ${INTEGRATION_TESTS:-N/A}"
echo "   Tests PASS: ${INTEGRATION_PASS}"
echo "   Tests FAIL: ${INTEGRATION_FAIL}"
echo ""

# Analizar patrón de contenedores
echo "🔍 Análisis de uso de contenedores:"
if [ -f /tmp/container_monitor.log ]; then
    MAX_CONTAINERS=$(sort -t: -k4 -n /tmp/container_monitor.log | tail -1 | grep -oP '\d+$')
    MIN_CONTAINERS=$(sort -t: -k4 -n /tmp/container_monitor.log | head -1 | grep -oP '\d+$')
    AVG_CONTAINERS=$(awk -F': ' '{sum+=$2; count++} END {printf "%.1f", sum/count}' /tmp/container_monitor.log)
    
    echo "   Máximo simultáneo: ${MAX_CONTAINERS:-0}"
    echo "   Mínimo: ${MIN_CONTAINERS:-0}"
    echo "   Promedio: ${AVG_CONTAINERS:-0}"
    
    # Detectar problema de recreación constante
    SAMPLES=$(wc -l < /tmp/container_monitor.log)
    if [ "$MAX_CONTAINERS" -gt 6 ]; then
        echo "   ⚠️  ADVERTENCIA: Más de 6 contenedores simultáneos detectados"
        echo "   ⚠️  Posible problema: contenedores no se están reutilizando"
    fi
    
    if [ "$SAMPLES" -gt 100 ] && [ "$MAX_CONTAINERS" -gt 3 ]; then
        echo "   ⚠️  ADVERTENCIA: Muchas fluctuaciones en contenedores"
        echo "   ⚠️  Posible problema: contenedores se crean/destruyen constantemente"
    fi
fi
echo ""

# Analizar tiempos individuales de tests de integración
echo "📊 Tiempos de tests de integración individuales:"
grep -E "^--- PASS:|^--- FAIL:" /tmp/integration_tests.log | while read -r line; do
    TEST_NAME=$(echo "$line" | awk '{print $3}')
    TEST_TIME=$(echo "$line" | grep -oP '\(\K[0-9.]+(?=s\))')
    if [ -n "$TEST_TIME" ]; then
        TIME_INT=$(printf "%.0f" "$TEST_TIME")
        if [ "$TIME_INT" -gt 15 ]; then
            echo "   ⚠️  $TEST_NAME: ${TEST_TIME}s (LENTO)"
        else
            echo "   ✅ $TEST_NAME: ${TEST_TIME}s"
        fi
    fi
done
echo ""

# Resumen Final
TOTAL_TIME=$((UNIT_TIME + INTEGRATION_TIME))
echo "═══════════════════════════════════════════════════════════════════════"
echo "📊 RESUMEN FINAL"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""
echo "⏱️  Tiempo Total: $(format_time $TOTAL_TIME) (${TOTAL_TIME}s)"
echo "   • Tests Unitarios: $(format_time $UNIT_TIME) (${UNIT_TIME}s)"
echo "   • Tests Integración: $(format_time $INTEGRATION_TIME) (${INTEGRATION_TIME}s)"
echo ""
echo "📊 Resultados:"
echo "   • Tests Unitarios: $([ $UNIT_EXIT_CODE -eq 0 ] && echo '✅ PASS' || echo '❌ FAIL')"
echo "   • Tests Integración: $([ $INTEGRATION_EXIT_CODE -eq 0 ] && echo '✅ PASS' || echo '❌ FAIL')"
echo ""

# Verificar contenedores residuales
RESIDUAL_CONTAINERS=$(docker ps -a | grep -cE "postgres|mongo|rabbitmq" || echo "0")
if [ "$RESIDUAL_CONTAINERS" -gt 0 ]; then
    echo "⚠️  Contenedores residuales: $RESIDUAL_CONTAINERS"
    echo "   (Estos deberían limpiarse automáticamente)"
    docker ps -a | grep -E "postgres|mongo|rabbitmq" | awk '{print "   - " $2 " (" $1 ")"}'
else
    echo "✅ No hay contenedores residuales"
fi
echo ""

# Recomendaciones de performance
echo "═══════════════════════════════════════════════════════════════════════"
echo "💡 RECOMENDACIONES DE PERFORMANCE"
echo "═══════════════════════════════════════════════════════════════════════"
echo ""

if [ "$INTEGRATION_TIME" -gt 300 ]; then
    echo "⚠️  Tests de integración muy lentos (>5 min)"
    echo "   Recomendaciones:"
    echo "   1. Verificar reutilización de contenedores"
    echo "   2. Considerar paralelización de tests"
    echo "   3. Revisar cleanup entre tests"
fi

if [ "${MAX_CONTAINERS:-0}" -gt 6 ]; then
    echo "⚠️  Demasiados contenedores simultáneos"
    echo "   Recomendaciones:"
    echo "   1. Implementar pool de contenedores compartidos"
    echo "   2. Usar testcontainers.Reuse(true)"
    echo "   3. Revisar que cleanup no destruya contenedores prematuramente"
fi

if [ "$RESIDUAL_CONTAINERS" -gt 0 ]; then
    echo "⚠️  Contenedores no se limpian correctamente"
    echo "   Recomendaciones:"
    echo "   1. Verificar defer statements en tests"
    echo "   2. Usar t.Cleanup() para garantizar limpieza"
    echo "   3. Revisar que testcontainers.Terminate() se llame"
fi

if [ "$INTEGRATION_TIME" -lt 180 ] && [ "${MAX_CONTAINERS:-0}" -le 3 ]; then
    echo "✅ Performance de tests es buena"
    echo "   • Tiempo razonable (<3 min)"
    echo "   • Uso eficiente de contenedores"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════════════"

# Cleanup de archivos temporales
rm -f /tmp/container_monitor.log

# Exit code general
if [ $UNIT_EXIT_CODE -eq 0 ] && [ $INTEGRATION_EXIT_CODE -eq 0 ]; then
    exit 0
else
    exit 1
fi
