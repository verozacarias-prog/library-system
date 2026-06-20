#!/bin/bash

BASE="http://localhost:3000"

# Login como admin
TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com", "password": "secret123"}' | jq -r '.access_token')

# Login como usuario normal
LECTOR_TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "lector@test.com", "password": "secret123"}' | jq -r '.access_token')

echo "Tokens obtenidos"
echo ""

# ==========================================
echo "=== UPDATE BOOK ==="
# ==========================================

echo "--- Caso feliz: actualizar título (admin) ---"
curl -s -X PATCH $BASE/books/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Clean Code - Updated Edition"}' | jq
echo ""

echo "--- Caso feliz: actualizar copias ---"
curl -s -X PATCH $BASE/books/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"available_copies": 10}' | jq
echo ""

echo "--- Error: libro que no existe (404) ---"
curl -s -X PATCH $BASE/books/9999 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Fantasma"}' | jq
echo ""

echo "--- Error: sin token (401) ---"
curl -s -X PATCH $BASE/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Sin auth"}' | jq
echo ""

echo "--- Error: usuario normal (403) ---"
curl -s -X PATCH $BASE/books/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LECTOR_TOKEN" \
  -d '{"title": "Intento fallido"}' | jq
echo ""

echo "--- Error: campo desconocido (400) ---"
curl -s -X PATCH $BASE/books/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Clean Code", "campo_inventado": "xyz"}' | jq
echo ""

echo "--- Error: copias negativas (400) ---"
curl -s -X PATCH $BASE/books/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"available_copies": -1}' | jq
echo ""

echo "--- Error: año futuro (400) ---"
curl -s -X PATCH $BASE/books/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"year": 2099}' | jq
echo ""

# ==========================================
echo "=== DELETE BOOK ==="
# ==========================================

echo "--- Error: sin token (401) ---"
curl -s -X DELETE $BASE/books/4 | jq
echo ""

echo "--- Error: usuario normal (403) ---"
curl -s -X DELETE $BASE/books/4 \
  -H "Authorization: Bearer $LECTOR_TOKEN" | jq
echo ""

echo "--- Error: libro que no existe (404) ---"
curl -s -X DELETE $BASE/books/9999 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Caso feliz: eliminar libro 4 (admin) ---"
curl -s -o /dev/null -w "%{http_code}" -X DELETE $BASE/books/4 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Verificar que ya no existe ---"
curl -s $BASE/books/4 | jq
echo ""

# ==========================================
echo "=== UPDATE USER ==="
# ==========================================

echo "--- Caso feliz: admin actualiza nombre de usuario 3 ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Lector Actualizado"}' | jq
echo ""

echo "--- Caso feliz: admin cambia rol de usuario 3 a admin ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"role": "admin"}' | jq
echo ""

echo "--- Caso feliz: admin revierte rol a user ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"role": "user"}' | jq
echo ""

echo "--- Error: usuario normal intenta cambiar su propio rol (403) ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $LECTOR_TOKEN" \
  -d '{"role": "admin"}' | jq
echo ""

echo "--- Error: rol inválido (400) ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"role": "superadmin"}' | jq
echo ""

echo "--- Error: usuario que no existe (404) ---"
curl -s -X PATCH $BASE/users/9999 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Nadie"}' | jq
echo ""

echo "--- Error: sin token (401) ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -d '{"name": "Sin auth"}' | jq
echo ""

echo "--- Error: campo desconocido (400) ---"
curl -s -X PATCH $BASE/users/3 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Test", "campo_raro": "x"}' | jq
echo ""

# ==========================================
echo "=== DELETE USER ==="
# ==========================================

echo "--- Error: sin token (401) ---"
curl -s -X DELETE $BASE/users/2 | jq
echo ""

echo "--- Error: usuario normal (403) ---"
curl -s -X DELETE $BASE/users/2 \
  -H "Authorization: Bearer $LECTOR_TOKEN" | jq
echo ""

echo "--- Error: usuario que no existe (404) ---"
curl -s -X DELETE $BASE/users/9999 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Caso feliz: admin elimina usuario 2 (carlos) ---"
curl -s -o /dev/null -w "%{http_code}" -X DELETE $BASE/users/2 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Verificar que ya no existe ---"
curl -s $BASE/users/2 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Verificar usuarios restantes ---"
curl -s $BASE/users \
  -H "Authorization: Bearer $TOKEN" | jq