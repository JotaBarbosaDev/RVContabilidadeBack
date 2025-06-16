# 📊 RV Contabilidade - Sistema de Gestão Backend

Sistema completo de gestão contabilística com aprovação de clientes, desenvolvido em Go com Gin Framework.

## 🏗️ Arquitetura do Sistema

### Estrutura de Entidades

1. **User** - Utilizadores (clientes, contabilistas, admins)
2. **Company** - Empresas dos clientes
3. **RegistrationRequest** - Histórico de solicitações de registo

### Status de Utilizadores

- `pending` - Aguarda aprovação da contabilista
- `approved` - Aprovado, pode aceder ao sistema
- `rejected` - Rejeitado, sem acesso (mas dados mantidos)
- `blocked` - Bloqueado (pode ter estado ativo antes)

## 🔐 Roles e Permissões

- **client** - Clientes aprovados (acesso limitado)
- **accountant** - Contabilistas (podem aprovar clientes)
- **admin** - Administradores (acesso total)

## 🚀 Como Funciona

### 1. Processo de Registo
1. Cliente submete formulário com dados pessoais e da empresa
2. Sistema cria utilizador com status `pending`
3. Cria solicitação de registo com todos os dados
4. Contabilista recebe notificação de nova solicitação

### 2. Processo de Aprovação
1. Contabilista revê solicitação pendente
2. Aprova ou rejeita com notas
3. Se aprovado: cliente pode fazer login e aceder ao sistema
4. Se rejeitado: dados mantidos para futuras submissões

### 3. Deteção de Duplicados
- Sistema verifica NIF existente
- Se já existe, mostra status atual da conta
- Permite re-submissão apenas se foi rejeitado

## 📡 Endpoints da API

### Autenticação (Público)
```
POST /api/auth/register          # Registo de novo cliente
POST /api/auth/login             # Login (todos os utilizadores)
POST /api/auth/logout            # Logout
POST /api/auth/register-direct   # Registo direto (interno)
```

### Administração (Contabilistas/Admin)
```
GET  /api/admin/pending-requests     # Solicitações pendentes
POST /api/admin/approve-request      # Aprovar/rejeitar solicitação
GET  /api/admin/requests             # Histórico de solicitações
GET  /api/admin/requests/:id         # Detalhes de solicitação
GET  /api/admin/users                # Listar utilizadores
GET  /api/admin/users/:id            # Detalhes de utilizador
PUT  /api/admin/users/:id/status     # Alterar status de utilizador
```

### Cliente (Área Protegida)
```
GET  /api/client/profile             # Ver perfil
PUT  /api/client/profile             # Atualizar perfil
GET  /api/client/company             # Ver empresa
PUT  /api/client/company             # Atualizar empresa
GET  /api/client/requests            # Histórico de solicitações
```

### Geral (Autenticados)
```
GET  /api/profile                    # Perfil atual
GET  /api/info                       # Informações da API
```

## 🛠️ Instalação e Configuração

### Pré-requisitos
- Go 1.21+
- PostgreSQL
- Git

### 1. Clonar Repositório
```bash
git clone [url-do-repositorio]
cd RVContabilidadeBack
```

### 2. Configurar Base de Dados
```sql
-- PostgreSQL
CREATE DATABASE RVContabilidadeDB;
```

### 3. Configurar Variáveis de Ambiente
```bash
# .env
JWT_SECRET=seu-jwt-secret-aqui
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=RVContabilidadeDB
```

### 4. Instalar Dependências
```bash
go mod tidy
```

### 5. Executar Migrações
```bash
# As migrações são executadas automaticamente no arranque
go run main.go
```

### 6. Iniciar Servidor
```bash
go run main.go
# ou
go build -o main . && ./main
```

O servidor estará disponível em `http://localhost:8080`

## 🧪 Testar o Sistema

### Executar Teste Completo
```bash
# Executar script de teste (requer curl e jq)
./test-api.sh
```

### Utilizadores Pré-criados
- **Admin**: `admin@rvcontabilidade.com` / `admin123`
- **Contabilista**: `contabilista@rvcontabilidade.com` / `contabilista123`

### Exemplo de Registo de Cliente
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@exemplo.com",
    "phone": "912345678",
    "nif": "123456789",
    "password": "password123",
    "company_name": "Silva & Associados Lda",
    "trade_name": "Silva Consultoria",
    "nipc": "123456789",
    "address": "Rua das Flores, 123",
    "postal_code": "1000-001",
    "city": "Lisboa",
    "country": "Portugal",
    "cae": "69200",
    "legal_form": "Sociedade por Quotas",
    "share_capital": 5000.00,
    "registration_date": "2024-01-15"
  }'
```

## 📚 Documentação

### Swagger UI
Aceder a `http://localhost:8080/swagger/index.html` para documentação interativa.

### Gerar Documentação Swagger
```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar docs
swag init -g main.go --output ./docs
```

## 🗄️ Estrutura da Base de Dados

### Tabela Users
```sql
id SERIAL PRIMARY KEY
email VARCHAR UNIQUE NOT NULL
password VARCHAR NOT NULL
name VARCHAR NOT NULL
phone VARCHAR
nif VARCHAR UNIQUE NOT NULL
role VARCHAR DEFAULT 'client'
status VARCHAR DEFAULT 'pending'
created_at TIMESTAMP
updated_at TIMESTAMP
```

### Tabela Companies
```sql
id SERIAL PRIMARY KEY
user_id INTEGER REFERENCES users(id)
company_name VARCHAR NOT NULL
trade_name VARCHAR
nipc VARCHAR UNIQUE NOT NULL
address VARCHAR
postal_code VARCHAR
city VARCHAR
country VARCHAR DEFAULT 'Portugal'
cae VARCHAR
legal_form VARCHAR
share_capital DECIMAL
registration_date DATE
status VARCHAR DEFAULT 'active'
created_at TIMESTAMP
updated_at TIMESTAMP
```

### Tabela Registration_Requests
```sql
id SERIAL PRIMARY KEY
user_id INTEGER REFERENCES users(id)
request_data TEXT
status VARCHAR DEFAULT 'pending'
submitted_at TIMESTAMP
reviewed_at TIMESTAMP
reviewed_by INTEGER REFERENCES users(id)
review_notes TEXT
approval_token VARCHAR UNIQUE
created_at TIMESTAMP
updated_at TIMESTAMP
```

## 🔧 Funcionalidades Implementadas

- ✅ Sistema de autenticação JWT
- ✅ Registo de clientes com aprovação
- ✅ Gestão de utilizadores e roles
- ✅ Deteção de duplicados por NIF
- ✅ Histórico completo de solicitações
- ✅ Área administrativa para contabilistas
- ✅ Área do cliente protegida
- ✅ Middleware de autorização
- ✅ Documentação Swagger
- ✅ Migrações automáticas
- ✅ Utilizadores padrão (admin/contabilista)
- ✅ Validação de dados de entrada
- ✅ Gestão de empresas
- ✅ Responses padronizadas
- ✅ Tratamento de erros

## 🛡️ Segurança

- Passwords encriptadas com bcrypt
- Tokens JWT com expiração
- Validação de roles em todas as rotas protegidas
- Verificação de status de utilizador em cada request
- Headers CORS configurados
- Validação de entrada com Gin binding

## 🎯 Próximos Passos

- [ ] Notificações por email
- [ ] Upload de documentos
- [ ] Logs de auditoria
- [ ] API de estatísticas
- [ ] Dashboard com gráficos
- [ ] Exportação de dados
- [ ] Backup automático
- [ ] Rate limiting
- [ ] Monitorização

## 🤝 Contribuição

1. Fork o projeto
2. Criar branch para feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit das alterações (`git commit -am 'Adicionar nova funcionalidade'`)
4. Push para branch (`git push origin feature/nova-funcionalidade`)
5. Abrir Pull Request

## 📄 Licença

Este projeto está sob licença MIT. Ver ficheiro `LICENSE` para mais detalhes.

## 📞 Suporte

Para questões e suporte:
- Email: suporte@rvcontabilidade.com
- Issues: [GitHub Issues](link-para-issues)

---

**RV Contabilidade** - Sistema de Gestão Contabilística
