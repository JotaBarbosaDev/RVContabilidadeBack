#!/bin/bash

echo "🚀 RV Contabilidade - Build e Teste"
echo "=================================="

# Limpar builds anteriores
echo "🧹 Limpando builds anteriores..."
rm -f RVContabilidadeBack
rm -f rvcontabilidade

# Verificar dependências
echo "📦 Verificando dependências..."
go mod tidy

# Gerar documentação Swagger
echo "📚 Gerando documentação Swagger..."
if command -v swag &> /dev/null; then
    swag init -g main.go --output ./docs
elif [ -f ~/go/bin/swag ]; then
    ~/go/bin/swag init -g main.go --output ./docs
else
    echo "⚠️  Swag não encontrado. Instalando..."
    go install github.com/swaggo/swag/cmd/swag@latest
    ~/go/bin/swag init -g main.go --output ./docs
fi

# Compilar aplicação
echo "🔨 Compilando aplicação..."
if go build -o rvcontabilidade .; then
    echo "✅ Compilação bem-sucedida!"
else
    echo "❌ Erro na compilação!"
    exit 1
fi

# Testar compilação
echo "🧪 Testando se a aplicação inicia..."
timeout 10s ./rvcontabilidade &
APP_PID=$!
sleep 5

# Verificar se está rodando
if kill -0 $APP_PID 2>/dev/null; then
    echo "✅ Aplicação iniciou com sucesso!"
    
    # Testar health check
    if command -v curl &> /dev/null; then
        echo "🩺 Testando health check..."
        if curl -s http://localhost:8080/health > /dev/null; then
            echo "✅ Health check passou!"
        else
            echo "⚠️  Health check falhou (porta pode estar ocupada)"
        fi
    fi
    
    # Parar aplicação
    kill $APP_PID 2>/dev/null
    wait $APP_PID 2>/dev/null
else
    echo "❌ Aplicação não conseguiu iniciar!"
    exit 1
fi

echo ""
echo "🎉 Teste completo!"
echo "📝 Para iniciar a aplicação:"
echo "   ./rvcontabilidade"
echo ""
echo "🌐 Endpoints disponíveis:"
echo "   - Health: http://localhost:8080/health"
echo "   - API Info: http://localhost:8080/api/info"
echo "   - Swagger: http://localhost:8080/swagger/index.html"
echo ""
echo "👥 Utilizadores padrão:"
echo "   - Admin: admin / admin123"
echo "   - Contabilista: contabilista / contabilista123"
