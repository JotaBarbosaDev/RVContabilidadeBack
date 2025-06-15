APP_NAME=RVContabilidadeBack
GO=go

run:
	@echo "🚀 A correr o servidor..."
	$(GO) run main.go

build:
	@echo "🏗️  A compilar $(APP_NAME)..."
	$(GO) build -o $(APP_NAME) main.go

tidy:
	@echo "🧹 A organizar dependências..."
	$(GO) mod tidy

test:
	@echo "🧪 A correr testes (se existirem)..."
	$(GO) test ./...

clean:
	@echo "🧼 A limpar ficheiros binários..."
	rm -f $(APP_NAME)
