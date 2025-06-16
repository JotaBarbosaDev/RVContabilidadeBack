# Makefile para RV Contabilidade Backend

.PHONY: help run build test clean docs api-test install dev all

# Configurações
BINARY_NAME=rvcontabilidade
GO=go

help: ## Mostrar ajuda
	@echo "RV Contabilidade Backend - Comandos Disponíveis:"
	@echo "================================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Instalar dependências
	@echo "� Instalando dependências..."
	$(GO) mod download
	$(GO) mod tidy

run: ## Executar aplicação em modo desenvolvimento
	@echo "🚀 Iniciando servidor de desenvolvimento..."
	$(GO) run main.go

build: ## Compilar aplicação
	@echo "🔨 Compilando aplicação..."
	$(GO) build -o bin/$(BINARY_NAME) main.go
	@echo "✅ Binário criado: bin/$(BINARY_NAME)"

test: ## Executar testes
	@echo "🧪 Executando testes..."
	$(GO) test -v ./...

docs: ## Gerar documentação Swagger
	@echo "📚 Gerando documentação Swagger..."
	swag init -g main.go --output ./docs
	@echo "✅ Documentação gerada em: http://localhost:8080/swagger/index.html"

api-test: ## Testar API completa
	@echo "🧪 Testando API..."
	chmod +x test-api.sh
	./test-api.sh

clean: ## Limpar arquivos gerados
	@echo "� Limpando arquivos gerados..."
	rm -rf bin/
	rm -f $(BINARY_NAME)

dev: install docs run ## Configurar ambiente de desenvolvimento

all: clean build docs ## Executar pipeline completa
