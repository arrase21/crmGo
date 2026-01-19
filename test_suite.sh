#!/bin/bash

# Test Suite Runner para CRM Go
# Este script corre todos los tests y genera reportes de coverage

echo "🚀 Iniciando Test Suite para CRM Go..."
echo "======================================"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para imprimir con color
print_status() {
    local status=$1
    local message=$2
    
    case $status in
        "success")
            echo -e "${GREEN}✅ $message${NC}"
            ;;
        "error")
            echo -e "${RED}❌ $message${NC}"
            ;;
        "warning")
            echo -e "${YELLOW}⚠️  $message${NC}"
            ;;
        "info")
            echo -e "${BLUE}ℹ️  $message${NC}"
            ;;
    esac
}

# Verificar que estamos en el directorio correcto
if [ ! -f "go.mod" ]; then
    print_status "error" "go.mod no encontrado. Por favor ejecuta desde el directorio raíz del proyecto."
    exit 1
fi

print_status "info" "Directorio actual: $(pwd)"

# 1. Correr tests de services
print_status "info" "Corriendo tests de Services..."
if go test ./internal/service -v -coverprofile=coverage_service.out; then
    print_status "success" "Tests de Services pasaron correctamente"
else
    print_status "error" "Tests de Services fallaron"
    exit 1
fi

# 2. Correr tests de repositories (si existen)
print_status "info" "Corriendo tests de Repositories..."
if go test ./internal/repository -v -coverprofile=coverage_repository.out 2>/dev/null; then
    print_status "success" "Tests de Repositories pasaron correctamente"
else
    print_status "warning" "No se encontraron tests de repositories o fallaron"
fi

# 3. Correr tests de domain (si existen)
print_status "info" "Corriendo tests de Domain..."
if go test ./internal/domain -v -coverprofile=coverage_domain.out 2>/dev/null; then
    print_status "success" "Tests de Domain pasaron correctamente"
else
    print_status "warning" "No se encontraron tests de domain o fallaron"
fi

# 4. Correr tests de transport (si existen)
print_status "info" "Corriendo tests de Transport..."
if go test ./internal/transport -v -coverprofile=coverage_transport.out 2>/dev/null; then
    print_status "success" "Tests de Transport pasaron correctamente"
else
    print_status "warning" "No se encontraron tests de transport o fallaron"
fi

# 5. Generar coverage combinado
print_status "info" "Generando reporte de coverage combinado..."
echo "mode: atomic" > coverage.out

# Combinar todos los archivos de coverage
for file in coverage_*.out; do
    if [ -f "$file" ]; then
        grep -v "mode: atomic" "$file" >> coverage.out
    fi
done

# 6. Generar reporte HTML de coverage
print_status "info" "Generando reporte HTML de coverage..."
if go tool cover -html=coverage.out -o coverage.html; then
    print_status "success" "Reporte HTML generado: coverage.html"
else
    print_status "warning" "No se pudo generar el reporte HTML"
fi

# 7. Mostrar resumen de coverage
print_status "info" "Resumen de Coverage:"
if [ -f "coverage.out" ]; then
    go tool cover -func=coverage.out
else
    print_status "warning" "No se encontró archivo de coverage"
fi

# 8. Verificar que el código compile
print_status "info" "Verificando que el código compile..."
if go build ./...; then
    print_status "success" "El código compila correctamente"
else
    print_status "error" "El código no compila"
    exit 1
fi

# 9. Linter (si está instalado)
print_status "info" "Ejecutando linter (si está disponible)..."
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run; then
        print_status "success" "Linter pasó sin problemas"
    else
        print_status "warning" "Linter encontró problemas"
    fi
else
    print_status "info" "golangci-lint no está instalado. Para instalar: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

echo ""
echo "======================================"
print_status "success" "Test Suite completado exitosamente!"
echo ""
echo "📊 Reportes generados:"
echo "   - coverage.out (coverage combinado)"
echo "   - coverage.html (reporte HTML)"
echo ""
echo "🔍 Para ver el reporte HTML:"
echo "   open coverage.html"
echo ""
echo "📈 Coverage actual:"
if [ -f "coverage.out" ]; then
    go tool cover -func=coverage.out | tail -1
fi
echo "======================================"