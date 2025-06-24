# Makefile para RV Contabilidade Backend

.PHONY: help run build test clean docs docker dev all start stop status logs

# Configurações
BINARY_NAME=rvcontabilidade
GO=go
DOCKER_IMAGE=rv-contabilidade-backend
PORT=8080

help: ## Mostrar ajuda
	@echo "RV Contabilidade Backend - Comandos Disponíveis:"
	@echo "================================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Instalar dependências Go
	@echo "📦 Instalando dependências..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependências instaladas!"

install-tools: ## Instalar ferramentas de desenvolvimento
	@echo "🛠️ Instalando ferramentas..."
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Ferramentas instaladas!"

run: ## Executar aplicação em modo desenvolvimento
	@echo "🚀 Iniciando servidor de desenvolvimento..."
	$(GO) run main.go

build: ## Compilar aplicação
	@echo "🔨 Compilando aplicação..."
	mkdir -p bin
	$(GO) build -o bin/$(BINARY_NAME) .
	@echo "✅ Binário criado: bin/$(BINARY_NAME)"

build-release: ## Compilar para produção
	@echo "🏗️ Compilando para produção..."
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/$(BINARY_NAME) .
	@echo "✅ Build de produção criado!"

test: ## Executar testes
	@echo "🧪 Executando testes..."
	$(GO) test -v ./...

test-coverage: ## Executar testes com coverage
	@echo "📊 Executando testes com coverage..."
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

docs: ## Gerar documentação Swagger
	@echo "📚 Gerando documentação Swagger..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g main.go --output ./docs; \
	elif [ -f ~/go/bin/swag ]; then \
		~/go/bin/swag init -g main.go --output ./docs; \
	else \
		echo "Installing swag..."; \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
		~/go/bin/swag init -g main.go --output ./docs; \
	fi
	@echo "✅ Documentação gerada! Acesse: http://localhost:$(PORT)/swagger/index.html"

clean: ## Limpar arquivos gerados
	@echo "🧹 Limpando arquivos gerados..."
	rm -rf bin/
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "✅ Limpeza concluída!"

lint: ## Executar linting
	@echo "🔍 Executando linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint não encontrado. Executando go vet..."; \
		$(GO) vet ./...; \
	fi

format: ## Formatar código
	@echo "💅 Formatando código..."
	$(GO) fmt ./...
	@echo "✅ Código formatado!"

start: build ## Compilar e iniciar aplicação
	@echo "🚀 Iniciando aplicação..."
	./bin/$(BINARY_NAME) &
	@echo "✅ Aplicação iniciada! PID: $$!"

stop: ## Parar aplicação
	@echo "⏹️ Parando aplicação..."
	@pkill -f $(BINARY_NAME) || echo "Aplicação não estava rodando"

status: ## Verificar status da aplicação
	@echo "📊 Status da aplicação:"
	@if pgrep -f $(BINARY_NAME) > /dev/null; then \
		echo "✅ Aplicação está rodando"; \
		echo "🌐 Endpoints:"; \
		echo "   - Health: http://localhost:$(PORT)/health"; \
		echo "   - API Info: http://localhost:$(PORT)/api/info"; \
		echo "   - Swagger: http://localhost:$(PORT)/swagger/index.html"; \
	else \
		echo "❌ Aplicação não está rodando"; \
	fi

quick-test: ## Teste rápido da API
	@echo "⚡ Teste rápido da API..."
	@if curl -s http://localhost:$(PORT)/health > /dev/null; then \
		echo "✅ Health check passou!"; \
		curl -s http://localhost:$(PORT)/api/info | head -n 5; \
	else \
		echo "❌ API não está respondendo em localhost:$(PORT)"; \
	fi

dev-setup: install install-tools docs ## Configurar ambiente completo de desenvolvimento
	@echo "🎉 Ambiente de desenvolvimento configurado!"
	@echo "Execute 'make run' para iniciar o servidor"

dev: dev-setup run ## Configurar e executar em modo desenvolvimento

all: clean install docs build ## Pipeline completa

# Comandos Docker (se necessário no futuro)
docker-build: ## Construir imagem Docker
	@echo "🐳 Construindo imagem Docker..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Executar com Docker
	@echo "🐳 Executando com Docker..."
	docker run -p $(PORT):$(PORT) $(DOCKER_IMAGE)

# Comando para CI/CD
ci: install lint test build ## Pipeline de CI/CD
