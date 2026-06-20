#!/bin/bash

BASE="http://localhost:3000"

# ==========================================
# Setup tokens
# ==========================================
TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com", "password": "secret123"}' | jq -r '.access_token')

LECTOR_TOKEN=$(curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "lector@test.com", "password": "secret123"}' | jq -r '.access_token')

echo "Tokens obtenidos"
echo ""

# ==========================================
echo "=== AUTH ==="
# ==========================================

echo "--- Error: password incorrecto (401) ---"
curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com", "password": "wrongpassword"}' | jq
echo ""

echo "--- Error: usuario no existe (401) ---"
curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "noexiste@test.com", "password": "secret123"}' | jq
echo ""

echo "--- Error: email faltante (400) ---"
curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password": "secret123"}' | jq
echo ""

echo "--- Error: password faltante (400) ---"
curl -s -X POST $BASE/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "vero@test.com"}' | jq
echo ""

# ==========================================
echo "=== REGISTER USER ==="
# ==========================================

echo "--- Error: email inválido (400) ---"
curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "noesunemail", "password": "secret123"}' | jq
echo ""

echo "--- Error: password muy corto (400) ---"
curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "test2@test.com", "password": "123"}' | jq
echo ""

echo "--- Error: email duplicado (409) ---"
curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Otro", "email": "vero@test.com", "password": "secret123"}' | jq
echo ""

echo "--- Error: rol inválido (400) ---"
curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "test3@test.com", "password": "secret123", "role": "superadmin"}' | jq
echo ""

echo "--- Error: nombre faltante (400) ---"
curl -s -X POST $BASE/users \
  -H "Content-Type: application/json" \
  -d '{"email": "test4@test.com", "password": "secret123"}' | jq
echo ""

# ==========================================
echo "=== GET USERS ==="
# ==========================================

echo "--- Caso feliz: admin lista todos los usuarios ---"
curl -s $BASE/users \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Error: usuario normal lista usuarios (403) ---"
curl -s $BASE/users \
  -H "Authorization: Bearer $LECTOR_TOKEN" | jq
echo ""

echo "--- Error: sin token (401) ---"
curl -s $BASE/users | jq
echo ""

echo "--- Caso feliz: obtener usuario por id ---"
curl -s $BASE/users/1 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Error: usuario que no existe (404) ---"
curl -s $BASE/users/9999 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

# ==========================================
echo "=== GET BOOKS ==="
# ==========================================

echo "--- Caso feliz: listar todos sin filtros ---"
curl -s "$BASE/books" | jq
echo ""

echo "--- Caso feliz: filtrar por autor ---"
curl -s "$BASE/books?author=Martin" | jq '.data[] | {id, title, author}'
echo ""

echo "--- Caso feliz: filtrar por género ---"
curl -s "$BASE/books?genre=tech" | jq '.data[] | {id, title, genre}'
echo ""

echo "--- Caso feliz: filtrar por disponibilidad ---"
curl -s "$BASE/books?available=true" | jq '.data[] | {id, title, available_copies}'
echo ""

echo "--- Caso feliz: paginación página 1 con limit 2 ---"
curl -s "$BASE/books?page=1&limit=2" | jq '{total, count: (.data | length)}'
echo ""

echo "--- Caso feliz: paginación página 2 con limit 2 ---"
curl -s "$BASE/books?page=2&limit=2" | jq '{total, count: (.data | length)}'
echo ""

echo "--- Caso feliz: filtros combinados ---"
curl -s "$BASE/books?author=Martin&genre=tech&available=true" | jq '.data[] | {id, title}'
echo ""

echo "--- Caso feliz: obtener libro por id ---"
curl -s "$BASE/books/1" | jq
echo ""

echo "--- Error: libro que no existe (404) ---"
curl -s "$BASE/books/9999" | jq
echo ""

# ==========================================
echo "=== CREATE BOOK - validaciones ==="
# ==========================================

echo "--- Error: ISBN inválido (400) ---"
curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Test", "author": "Test", "isbn": "123", "year": 2020, "genre": "tech", "available_copies": 1}' | jq
echo ""

echo "--- Error: año futuro (400) ---"
curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Test", "author": "Test", "isbn": "9780132350885", "year": 2099, "genre": "tech", "available_copies": 1}' | jq
echo ""

echo "--- Error: copias 0 en creación (400) ---"
curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Test", "author": "Test", "isbn": "9780132350885", "year": 2020, "genre": "tech", "available_copies": 0}' | jq
echo ""

echo "--- Error: ISBN duplicado (409) ---"
curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Otro titulo", "author": "Otro autor", "isbn": "9780132350884", "year": 2020, "genre": "tech", "available_copies": 1}' | jq
echo ""

echo "--- Error: campos faltantes (400) ---"
curl -s -X POST $BASE/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title": "Solo titulo"}' | jq
echo ""

# ==========================================
echo "=== LOANS ==="
# ==========================================

echo "--- Caso feliz: crear préstamo (libro 2, usuario logueado = vero) ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"book_id": 2}' | jq
echo ""

echo "--- Error: préstamo duplicado - mismo user y libro activo (409) ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"book_id": 1}' | jq
echo ""

echo "--- Error: libro que no existe (404) ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"book_id": 9999}' | jq
echo ""

echo "--- Error: book_id inválido (400) ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"book_id": -1}' | jq
echo ""

echo "--- Error: body inválido (400) ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"invalid": "body"}' | jq
echo ""

echo "--- Error: sin token (401) ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -d '{"book_id": 1}' | jq
echo ""

echo "--- Verificar copias decrementadas libro 2 ---"
curl -s $BASE/books/2 | jq '{id, title, available_copies}'
echo ""

echo "--- Caso feliz: préstamos activos usuario 1 (admin viendo los suyos) ---"
curl -s $BASE/loans/users/1 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Caso feliz: historial usuario 1 ---"
curl -s $BASE/loans/users/1/history \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Error: lector intenta ver préstamos de usuario 1 (403) ---"
curl -s $BASE/loans/users/1 \
  -H "Authorization: Bearer $LECTOR_TOKEN" | jq
echo ""

echo "--- Error: lector intenta ver historial de usuario 1 (403) ---"
curl -s $BASE/loans/users/1/history \
  -H "Authorization: Bearer $LECTOR_TOKEN" | jq
echo ""

echo "--- Caso feliz: lector ve sus propios préstamos (userId = 3) ---"
curl -s $BASE/loans/users/3 \
  -H "Authorization: Bearer $LECTOR_TOKEN" | jq
echo ""

echo "--- Error: sin token en loans by user (401) ---"
curl -s $BASE/loans/users/1 | jq
echo ""

echo "--- Caso feliz: devolver libro (loan 1) ---"
curl -s -X PATCH $BASE/loans/1 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Error: devolver préstamo ya devuelto (409) ---"
curl -s -X PATCH $BASE/loans/1 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Error: préstamo que no existe (404) ---"
curl -s -X PATCH $BASE/loans/9999 \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Verificar copias restauradas libro 1 ---"
curl -s $BASE/books/1 | jq '{id, title, available_copies}'
echo ""

echo "--- Historial usuario 1 con préstamo devuelto ---"
curl -s $BASE/loans/users/1/history \
  -H "Authorization: Bearer $TOKEN" | jq
echo ""

echo "--- Caso feliz: volver a pedir el mismo libro después de devolverlo ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"book_id": 1}' | jq
echo ""

echo "=== VALIDACIÓN DE COPIAS ==="

echo "--- Agotar copias libro 3 (4 copias, lector ya tiene 1 activo desde seed, crear 3 usuarios temporales para los slots restantes) ---"

echo "--- Verificar copias iniciales de libro 3 ---"
curl -s http://localhost:3000/books/3 | jq '{id, title, available_copies}'
echo ""

# Crear 3 usuarios temporales con sus tokens (autosuficiente, sin depender de seed)
curl -s -X POST $BASE/users -H "Content-Type: application/json" \
  -d '{"name": "Temp1", "email": "temp1@test.com", "password": "secret123", "role": "user"}' > /dev/null
TEMP1_TOKEN=$(curl -s -X POST $BASE/auth/login -H "Content-Type: application/json" \
  -d '{"email": "temp1@test.com", "password": "secret123"}' | jq -r '.access_token')

curl -s -X POST $BASE/users -H "Content-Type: application/json" \
  -d '{"name": "Temp2", "email": "temp2@test.com", "password": "secret123", "role": "user"}' > /dev/null
TEMP2_TOKEN=$(curl -s -X POST $BASE/auth/login -H "Content-Type: application/json" \
  -d '{"email": "temp2@test.com", "password": "secret123"}' | jq -r '.access_token')

curl -s -X POST $BASE/users -H "Content-Type: application/json" \
  -d '{"name": "Temp3", "email": "temp3@test.com", "password": "secret123", "role": "user"}' > /dev/null
TEMP3_TOKEN=$(curl -s -X POST $BASE/auth/login -H "Content-Type: application/json" \
  -d '{"email": "temp3@test.com", "password": "secret123"}' | jq -r '.access_token')

# slots 2, 3 y 4 de 4 (slot 1 = lector desde seed)
curl -s -X POST $BASE/loans -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TEMP1_TOKEN" -d '{"book_id": 3}' | jq '{id, status}'

curl -s -X POST $BASE/loans -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TEMP2_TOKEN" -d '{"book_id": 3}' | jq '{id, status}'

curl -s -X POST $BASE/loans -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" -d '{"book_id": 3}' | jq '{id, status}'
echo ""

echo "--- Verificar copias en 0 ---"
curl -s http://localhost:3000/books/3 | jq '{id, title, available_copies}'
echo ""

echo "--- Error: sin copias disponibles (409) — temp3 no tiene préstamo activo del libro, error real es de stock ---"
curl -s -X POST $BASE/loans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TEMP3_TOKEN" \
  -d '{"book_id": 3}' | jq
echo ""

echo "--- Health check loans-service ---"
curl -s http://localhost:8081/health | jq
echo ""
