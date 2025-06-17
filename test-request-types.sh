#!/bin/bash

echo "🧪 Teste do Campo request_type"
echo "=============================="

BASE_URL="http://localhost:8080"

echo ""
echo "1️⃣ Fazendo login como contabilista..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "contabilista",
    "password": "contabilista123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Erro: Não foi possível obter token"
    exit 1
fi

echo "✅ Login bem-sucedido, token obtido"

echo ""
echo "2️⃣ Verificando pedidos pendentes..."
PENDING_REQUESTS=$(curl -s -X GET "$BASE_URL/api/admin/pending-requests" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Pedidos pendentes:"
echo $PENDING_REQUESTS | jq '.'

echo ""
echo "3️⃣ Rejeitando o primeiro utilizador para permitir nova solicitação..."
# Rejeitar o primeiro request (assumindo ID 1)
REJECT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/admin/approve-request" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 1,
    "status": "rejected",
    "review_notes": "Teste: Rejeitado para testar cliente existente"
  }')

echo "Resposta da rejeição:"
echo $REJECT_RESPONSE | jq '.'

echo ""
echo "4️⃣ Aguardando 2 segundos..."
sleep 2

echo ""
echo "5️⃣ Criando nova solicitação com o mesmo NIF (cliente existente)..."
EXISTING_CLIENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "teste.cliente.nova",
    "name": "Cliente Teste Nova Tentativa",
    "email": "testenovatentativa@exemplo.com",
    "phone": "914567890",
    "nif": "987654321",
    "password": "password123",
    "company_name": "Empresa Teste Nova Tentativa Lda",
    "trade_name": "Teste Nova",
    "nipc": "987654324",
    "address": "Rua Nova Tentativa, 789",
    "postal_code": "3000-001",
    "city": "Aveiro",
    "country": "Portugal",
    "cae": "69200",
    "legal_form": "Sociedade por Quotas",
    "share_capital": 8000.00,
    "registration_date": "2024-04-15"
  }')

echo "Resposta do registo de cliente existente:"
echo $EXISTING_CLIENT_RESPONSE | jq '.'

echo ""
echo "6️⃣ Verificando novamente os pedidos pendentes para ver o request_type..."
FINAL_PENDING=$(curl -s -X GET "$BASE_URL/api/admin/pending-requests" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Pedidos pendentes finais:"
echo $FINAL_PENDING | jq '.'

echo ""
echo "✅ Teste concluído!"
