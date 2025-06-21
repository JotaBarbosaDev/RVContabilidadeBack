#!/bin/bash

echo "🧪 Testando aprovação de pedido"
echo "==============================="

# Fazer login
echo "1️⃣ Fazendo login como contabilista..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "contabilista",
    "password": "contabilista123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Erro: Token não obtido"
    echo "Resposta: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Token obtido: ${TOKEN:0:20}..."

# Verificar pedidos pendentes
echo ""
echo "2️⃣ Verificando pedidos pendentes..."
PENDING=$(curl -s -X GET http://localhost:8080/api/admin/pending-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json")

echo "Pedidos pendentes:"
echo $PENDING | python3 -m json.tool 2>/dev/null || echo $PENDING

# Testar aprovação
echo ""
echo "3️⃣ Testando aprovação do pedido ID 1..."
APPROVAL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/admin/approve-request \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "request_id": 1,
    "status": "approved",
    "review_notes": "Teste de aprovação via script"
  }')

echo "Resposta da aprovação:"
echo $APPROVAL_RESPONSE | python3 -m json.tool 2>/dev/null || echo $APPROVAL_RESPONSE

echo ""
echo "✅ Teste concluído!"
